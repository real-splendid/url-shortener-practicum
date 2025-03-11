package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/real-splendid/url-shortener-practicum/internal"
	"github.com/real-splendid/url-shortener-practicum/internal/storage"
	"github.com/real-splendid/url-shortener-practicum/mocks"
)

func TestHandleRedirection(t *testing.T) {
	originalURL := "https://ya.ru"
	key := "testtest"
	userID := "test-user-id"
	s := storage.NewMemoryStorage()
	_, err := s.Set(key, originalURL, userID)
	assert.NoError(t, err)
	handler := MakeRedirectionHandler(s, zap.NewNop().Sugar())
	t.Run("redirect", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/{key}", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("key", key)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

		handler(recorder, request)
		result := recorder.Result()
		defer func() {
			if closeErr := result.Body.Close(); closeErr != nil {
				t.Errorf("failed to close response body: %v", closeErr)
			}
		}()

		assert.Equal(t, http.StatusTemporaryRedirect, result.StatusCode)
		assert.Equal(t, originalURL, result.Header.Get("Location"))
	})
	t.Run("not-found", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/notfound", nil)

		handler(recorder, request)
		result := recorder.Result()
		defer func() {
			if closeErr := result.Body.Close(); closeErr != nil {
				t.Errorf("failed to close response body: %v", closeErr)
			}
		}()

		assert.Equal(t, http.StatusNotFound, result.StatusCode)
	})

	t.Run("deleted-url", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockStorage(ctrl)
		mockStorage.EXPECT().
			Get(gomock.Any()).
			Return("", internal.ErrURLDeleted)

		handler := MakeRedirectionHandler(mockStorage, zap.NewNop().Sugar())

		req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("key", "abc123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()
		handler(rec, req)

		assert.Equal(t, http.StatusGone, rec.Code)
	})
}
