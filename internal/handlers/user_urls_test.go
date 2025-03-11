package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/real-splendid/url-shortener-practicum/internal"
	"github.com/real-splendid/url-shortener-practicum/internal/middleware"
	"github.com/real-splendid/url-shortener-practicum/mocks"
)

func TestHandleUserURLs(t *testing.T) {
	baseURL := "http://localhost:8080"

	t.Run("storage-error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockStorage(ctrl)
		mockStorage.EXPECT().
			GetUserURLs(gomock.Any()).
			Return(nil, sql.ErrConnDone)

		handler := MakeUserURLsHandler(mockStorage, zap.NewNop().Sugar(), baseURL)
		req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)

		mockUserID := "test-user-id"
		mockSignedCookie := middleware.SignCookie(mockUserID)
		req.AddCookie(&http.Cookie{Name: "user_id", Value: mockSignedCookie})

		rec := httptest.NewRecorder()
		handler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockStorage(ctrl)
		mockStorage.EXPECT().
			GetUserURLs(gomock.Any()).
			Return([]internal.URLPair{{
				ShortURL:    "abc123",
				OriginalURL: "https://example.com",
			}}, nil)

		handler := MakeUserURLsHandler(mockStorage, zap.NewNop().Sugar(), baseURL)
		req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)

		mockUserID := "test-user-id"
		mockSignedCookie := middleware.SignCookie(mockUserID)
		req.AddCookie(&http.Cookie{Name: "user_id", Value: mockSignedCookie})

		rec := httptest.NewRecorder()
		handler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `[{"short_url":"http://localhost:8080/abc123","original_url":"https://example.com"}]`, rec.Body.String())
	})

	t.Run("no-cookie", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockStorage(ctrl)

		handler := MakeUserURLsHandler(mockStorage, zap.NewNop().Sugar(), baseURL)
		req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)

		rec := httptest.NewRecorder()
		handler(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("invalid-cookie", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockStorage(ctrl)

		handler := MakeUserURLsHandler(mockStorage, zap.NewNop().Sugar(), baseURL)
		req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
		req.AddCookie(&http.Cookie{Name: "user_id", Value: "invalid"})

		rec := httptest.NewRecorder()
		handler(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("empty-list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockStorage(ctrl)
		mockStorage.EXPECT().
			GetUserURLs(gomock.Any()).
			Return([]internal.URLPair{}, nil)

		handler := MakeUserURLsHandler(mockStorage, zap.NewNop().Sugar(), baseURL)
		req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)

		mockUserID := "test-user-id"
		mockSignedCookie := middleware.SignCookie(mockUserID)
		req.AddCookie(&http.Cookie{Name: "user_id", Value: mockSignedCookie})

		rec := httptest.NewRecorder()
		handler(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})
}
