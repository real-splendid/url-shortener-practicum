package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/real-splendid/url-shortener-practicum/internal/middleware"
	"github.com/real-splendid/url-shortener-practicum/internal/storage"
)

func TestHandleAPIShortenBatch(t *testing.T) {
	t.Run("create-short-links-batch", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		reqBody, err := json.Marshal([]ShortenBatchReq{
			{CorrelationID: "1", OriginalURL: "https://ya.ru"},
			{CorrelationID: "2", OriginalURL: "https://google.com"},
		})
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(reqBody))

		mockSignedCookie := middleware.SignCookie("test-user-id")
		req.AddCookie(&http.Cookie{Name: "user_id", Value: mockSignedCookie})

		handler := MakeAPIShortenBatchHandler(storage.NewMemoryStorage(), zap.NewNop().Sugar(), "http://localhost")
		handler(recorder, req)
		result := recorder.Result()
		body, err := io.ReadAll(result.Body)
		assert.NoError(t, err)
		var resp []ShortenBatchResp
		err = json.Unmarshal(body, &resp)
		assert.NoError(t, err)
		defer result.Body.Close()

		assert.Equal(t, http.StatusCreated, result.StatusCode)
		assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
		assert.NoError(t, err)
		assert.Len(t, resp, 2)
		assert.Equal(t, "1", resp[0].CorrelationID)
		assert.Regexp(t, "http://[^/]+/\\w{12}$", resp[0].ShortURL)
		assert.Equal(t, "2", resp[1].CorrelationID)
		assert.Regexp(t, "http://[^/]+/\\w{12}$", resp[1].ShortURL)
	})

	t.Run("unauthorized-request", func(t *testing.T) {
		reqBody, err := json.Marshal([]ShortenBatchReq{
			{CorrelationID: "1", OriginalURL: "https://ya.ru"},
		})
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(reqBody))
		recorder := httptest.NewRecorder()
		handler := MakeAPIShortenBatchHandler(storage.NewMemoryStorage(), zap.NewNop().Sugar(), "http://localhost")

		handler(recorder, req)

		result := recorder.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, result.StatusCode)
	})
}
