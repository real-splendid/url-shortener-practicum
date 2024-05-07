package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/real-splendid/url-shortener-practicum/internal/middleware"
	"github.com/real-splendid/url-shortener-practicum/internal/storage"
)

func TestHandleAPIShorten(t *testing.T) {
	t.Run("create-short-link", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		reqBody, err := json.Marshal(ShortenReq{URL: "https://ya.ru"})
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(reqBody))

		mockSignedCookie := middleware.SignCookie("test-user-id")
		req.AddCookie(&http.Cookie{Name: "user_id", Value: mockSignedCookie})

		handler := MakeAPIShortenHandler(storage.NewMemoryStorage(), zap.NewNop().Sugar(), "http://localhost")
		handler(recorder, req)

		var resp ShortenResp
		err = json.Unmarshal(recorder.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		assert.Regexp(t, "http://[^/]+/\\w{12}$", resp.Result)
	})

	t.Run("unauthorized-request", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		reqBody, err := json.Marshal(ShortenReq{URL: "https://ya.ru"})
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(reqBody))
		handler := MakeAPIShortenHandler(storage.NewMemoryStorage(), zap.NewNop().Sugar(), "http://localhost")

		handler(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})
}
