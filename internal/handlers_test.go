package internal

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	SetStorage(NewMemoryStorage())
	m.Run()
}

func TestHandleShorten(t *testing.T) {
	t.Run("create-short-link", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://ya.ru"))

		HandleShorten(recorder, request)
		result := recorder.Result()
		body, err := io.ReadAll(result.Body)
		defer result.Body.Close()

		assert.Equal(t, http.StatusCreated, result.StatusCode)
		assert.NoError(t, err)
		assert.Regexp(t, "http://[^/]+/\\w{12}$", string(body))
	})
}

func TestHandleAPIShorten(t *testing.T) {
	t.Run("create-short-link", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		reqBody, err := json.Marshal(ShortenReq{URL: "https://ya.ru"})
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(reqBody)))

		HandleAPIShorten(recorder, req)
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

func TestHandleRedirection(t *testing.T) {
	t.Run("redirect", func(t *testing.T) {
		key := "testtest"
		originalURL := "https://ya.ru"
		internalStorage.Set(key, originalURL)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/{key}", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("key", key)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

		HandleRedirection(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusTemporaryRedirect, result.StatusCode)
		assert.Equal(t, originalURL, result.Header.Get("Location"))
	})
	t.Run("not-found", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/notfound", nil)

		HandleRedirection(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusNotFound, result.StatusCode)
	})
}
