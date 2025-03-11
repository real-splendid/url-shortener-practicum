package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestLogMiddleware(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()

	t.Run("logging response writer", func(t *testing.T) {
		lw := &loggingResponseWriter{
			ResponseWriter: httptest.NewRecorder(),
			status:         0,
			size:           0,
		}

		lw.WriteHeader(http.StatusCreated)
		assert.Equal(t, http.StatusCreated, lw.status)

		n, err := lw.Write([]byte("test data"))
		require.NoError(t, err)
		assert.Equal(t, 9, n)
		assert.Equal(t, 9, lw.size)

		n, err = lw.Write([]byte(" more"))
		require.NoError(t, err)
		assert.Equal(t, 5, n)
		assert.Equal(t, 14, lw.size)
	})

	t.Run("log middleware", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test-path", nil)

		rr := httptest.NewRecorder()

		handlerCalled := false
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("response data"))
			require.NoError(t, err)
		})

		handler := MakeLogMiddleware(logger)(testHandler)
		handler.ServeHTTP(rr, req)

		assert.True(t, handlerCalled)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "response data", rr.Body.String())
	})

	t.Run("log middleware with error status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/error-path", nil)

		rr := httptest.NewRecorder()

		handlerCalled := false
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("error response"))
			require.NoError(t, err)
		})

		handler := MakeLogMiddleware(logger)(testHandler)
		handler.ServeHTTP(rr, req)

		assert.True(t, handlerCalled)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "error response", rr.Body.String())
	})
}
