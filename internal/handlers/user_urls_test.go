package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

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
}
