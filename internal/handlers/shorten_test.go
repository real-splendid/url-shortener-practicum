package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/real-splendid/url-shortener-practicum/internal"
	"github.com/real-splendid/url-shortener-practicum/internal/middleware"
	"github.com/real-splendid/url-shortener-practicum/internal/storage"
	"github.com/real-splendid/url-shortener-practicum/mocks"
)

func TestHandleShorten(t *testing.T) {
	t.Run("create-short-link", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://ya.ru"))

		mockUserID := "test-user-id"
		mockSignedCookie := middleware.SignCookie(mockUserID)
		request.AddCookie(&http.Cookie{Name: "user_id", Value: mockSignedCookie})

		handler := MakeShortenHandler(storage.NewMemoryStorage(), zap.NewNop().Sugar(), "http://localhost")
		handler(recorder, request)
		result := recorder.Result()
		body, err := io.ReadAll(result.Body)
		defer result.Body.Close()

		assert.Equal(t, http.StatusCreated, result.StatusCode)
		assert.NoError(t, err)
		assert.Regexp(t, "http://[^/]+/\\w{12}$", string(body))
	})

	t.Run("unique-violation", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://ya.ru"))

		mockUserID := "test-user-id"
		mockSignedCookie := middleware.SignCookie(mockUserID)
		request.AddCookie(&http.Cookie{Name: "user_id", Value: mockSignedCookie})

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		s := mocks.NewMockStorage(ctrl)
		s.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return("test", internal.ErrDuplicateKey).Times(1)
		handler := MakeShortenHandler(s, zap.NewNop().Sugar(), "http://localhost")

		handler(recorder, request)
		result := recorder.Result()
		body, err := io.ReadAll(result.Body)
		defer result.Body.Close()

		assert.Equal(t, http.StatusConflict, result.StatusCode)
		assert.NoError(t, err)
		assert.Regexp(t, "http://[^/]+/\\w+$", string(body))
	})
}
