package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"encoding/json"

	"github.com/real-splendid/url-shortener-practicum/internal/storage"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHandleAPIShorten(t *testing.T) {
	t.Run("create-short-link", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		reqBody, err := json.Marshal(ShortenReq{URL: "https://ya.ru"})
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(reqBody)))

		handler := MakeAPIShortenHandler(storage.NewMemoryStorage(), zap.NewNop().Sugar(), "http://localhost")
		handler(recorder, req)
		result := recorder.Result()
		body, err := io.ReadAll(result.Body)
		assert.NoError(t, err)
		resp := &ShortenResp{}
		err = json.Unmarshal(body, resp)
		assert.NoError(t, err)
		defer result.Body.Close()

		assert.Equal(t, http.StatusCreated, result.StatusCode)
		assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
		assert.NoError(t, err)
		assert.Regexp(t, "http://[^/]+/\\w{12}$", string(resp.Result))
	})
}
