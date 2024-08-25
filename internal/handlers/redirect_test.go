package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/real-splendid/url-shortener-practicum/internal/storage"
)

func TestHandleRedirection(t *testing.T) {
	originalURL := "https://ya.ru"
	key := "testtest"
	userID := "test-user-id"
	s := storage.NewMemoryStorage()
	s.Set(key, originalURL, userID)
	handler := MakeRedirectionHandler(s, zap.NewNop().Sugar())
	t.Run("redirect", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/{key}", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("key", key)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

		handler(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusTemporaryRedirect, result.StatusCode)
		assert.Equal(t, originalURL, result.Header.Get("Location"))
	})
	t.Run("not-found", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/notfound", nil)

		handler(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusNotFound, result.StatusCode)
	})
}
