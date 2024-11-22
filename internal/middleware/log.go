// Package middleware предоставляет middleware для HTTP-сервера.
package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

// Write записывает данные в ответ и обновляет размер ответа.
// Возвращает количество записанных байт и ошибку, если таковая имеется.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.size += size
	return size, err
}

// WriteHeader записывает код состояния HTTP и сохраняет его.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.status = statusCode
}

// MakeLogMiddleware создает middleware для логирования запросов.
// Принимает логгер для записи информации о запросе.
// Возвращает функцию, которая принимает http.Handler и возвращает http.Handler, обернутый middleware.
// Middleware записывает URI запроса, метод, время обработки, статус ответа и размер ответа.
func MakeLogMiddleware(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			lw := &loggingResponseWriter{ResponseWriter: w}
			next.ServeHTTP(lw, r)
			logger.Infow("request processed",
				"uri", r.RequestURI,
				"method", r.Method,
				"duration", time.Since(start),
				"status", lw.status,
				"size", lw.size,
			)
		})
	}
}
