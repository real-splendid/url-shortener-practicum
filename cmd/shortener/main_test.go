package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestHandleShortLinkCreation(t *testing.T) {
	t.Run("create-short-link", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://ya.ru"))

		handleShortLinkCreation(recorder, request)
		result := recorder.Result()
		body, err := io.ReadAll(result.Body)
		defer result.Body.Close()

		assert.Equal(t, http.StatusCreated, result.StatusCode)
		assert.NoError(t, err)
		assert.Regexp(t, "http://[^/]+/\\w{8}$", string(body))
	})
}

func TestHandleRedirection(t *testing.T) {
	t.Run("redirect", func(t *testing.T) {
		key := "testtest"
		originalURL := "https://ya.ru"
		storage[key] = originalURL
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/{key}", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("key", key)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

		handleRedirection(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusTemporaryRedirect, result.StatusCode)
		assert.Equal(t, originalURL, result.Header.Get("Location"))
	})
	t.Run("not-found", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/notfound", nil)

		handleRedirection(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusNotFound, result.StatusCode)
	})
}
