package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/real-splendid/url-shortener-practicum/internal/middleware"
	"github.com/real-splendid/url-shortener-practicum/internal/storage"
	"github.com/real-splendid/url-shortener-practicum/mocks"
)

func TestHandleDeleteUserURLs(t *testing.T) {
	t.Run("successful-delete", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockStorage(ctrl)
		mockStorage.EXPECT().
			DeleteUserURLs(gomock.Any(), gomock.Any()).
			Return(nil)

		handler := MakeDeleteUserURLsHandler(mockStorage, zap.NewNop().Sugar())

		urlsToDelete := []string{"abc123", "def456"}
		body, err := json.Marshal(urlsToDelete)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewReader(body))
		mockUserID := "test-user-id"
		mockSignedCookie := middleware.SignCookie(mockUserID)
		req.AddCookie(&http.Cookie{Name: "user_id", Value: mockSignedCookie})

		rec := httptest.NewRecorder()
		handler(rec, req)

		assert.Equal(t, http.StatusAccepted, rec.Code)
		time.Sleep(100 * time.Millisecond)
	})

	t.Run("invalid-json", func(t *testing.T) {
		handler := MakeDeleteUserURLsHandler(storage.NewMemoryStorage(), zap.NewNop().Sugar())

		req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewReader([]byte("invalid json")))
		mockUserID := "test-user-id"
		mockSignedCookie := middleware.SignCookie(mockUserID)
		req.AddCookie(&http.Cookie{Name: "user_id", Value: mockSignedCookie})

		rec := httptest.NewRecorder()
		handler(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
