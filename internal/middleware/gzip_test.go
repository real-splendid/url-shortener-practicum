package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGzipMiddleware(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	t.Run("request with gzip content-encoding", func(t *testing.T) {
		var buf bytes.Buffer
		gzipWriter := gzip.NewWriter(&buf)
		_, err := gzipWriter.Write([]byte("test data"))
		require.NoError(t, err)
		closeErr := gzipWriter.Close()
		require.NoError(t, closeErr)

		req := httptest.NewRequest(http.MethodPost, "/", &buf)
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Accept-Encoding", "gzip")

		rr := httptest.NewRecorder()

		handlerCalled := false
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			body, errRead := io.ReadAll(r.Body)
			require.NoError(t, errRead)
			assert.Equal(t, "test data", string(body))
			w.WriteHeader(http.StatusOK)
			_, errWrite := w.Write([]byte("response data"))
			require.NoError(t, errWrite)
		})

		handler := MakeGzipMiddleware(sugar)(testHandler)
		handler.ServeHTTP(rr, req)

		assert.True(t, handlerCalled)

		assert.Equal(t, "gzip", rr.Header().Get("Content-Encoding"))

		gzipReader, err := gzip.NewReader(rr.Body)
		require.NoError(t, err)
		defer func() {
			closeErr := gzipReader.Close()
			require.NoError(t, closeErr)
		}()

		decompressed, err := io.ReadAll(gzipReader)
		require.NoError(t, err)
		assert.Equal(t, "response data", string(decompressed))
	})

	t.Run("request without gzip support", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader("test data"))

		rr := httptest.NewRecorder()

		handlerCalled := false
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			assert.Equal(t, "test data", string(body))
			w.WriteHeader(http.StatusOK)
			_, err = w.Write([]byte("uncompressed response"))
			require.NoError(t, err)
		})

		handler := MakeGzipMiddleware(sugar)(testHandler)
		handler.ServeHTTP(rr, req)

		assert.True(t, handlerCalled)

		assert.Empty(t, rr.Header().Get("Content-Encoding"))
		assert.Equal(t, "uncompressed response", rr.Body.String())
	})

	t.Run("invalid gzip content", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not gzipped"))
		req.Header.Set("Content-Encoding", "gzip")

		rr := httptest.NewRecorder()

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Handler should not be called with invalid gzip content")
		})

		handler := MakeGzipMiddleware(sugar)(testHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("gzip writer test", func(t *testing.T) {
		mockResponseWriter := httptest.NewRecorder()
		var buf bytes.Buffer
		gzWriter := gzip.NewWriter(&buf)

		gzipW := gzipWriter{
			ResponseWriter: mockResponseWriter,
			Writer:         gzWriter,
		}

		n, err := gzipW.Write([]byte("test data"))
		require.NoError(t, err)
		assert.Equal(t, 9, n)

		closeErr := gzWriter.Close()
		require.NoError(t, closeErr)

		gzipReader, err := gzip.NewReader(&buf)
		require.NoError(t, err)
		defer func() {
			closeErr := gzipReader.Close()
			require.NoError(t, closeErr)
		}()

		decompressed, err := io.ReadAll(gzipReader)
		require.NoError(t, err)
		assert.Equal(t, "test data", string(decompressed))
	})
}
