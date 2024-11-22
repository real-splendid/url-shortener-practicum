package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHandlePing(t *testing.T) {
	t.Run("ping-on-bad-dsn", func(t *testing.T) {
		testDSN := "postgres://test:test@localhost:5432/test?sslmode=disable"

		handler := MakePingHandler(testDSN, zap.NewNop().Sugar())
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		rec := httptest.NewRecorder()

		handler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
