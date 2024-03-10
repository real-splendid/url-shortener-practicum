package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleShortLinkCreation(t *testing.T) {
	shortURLStorage = make(map[string]string)
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
	shortURLStorage = make(map[string]string)
	t.Run("redirect", func(t *testing.T) {
		key := "testtest"
		originalURL := "https://ya.ru"
		shortURLStorage[key] = originalURL
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/"+key, nil)

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
