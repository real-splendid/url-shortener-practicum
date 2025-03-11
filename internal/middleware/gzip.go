// Package middleware предоставляет middleware для HTTP-сервера.
package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Write записывает данные в gzip writer.
// Возвращает количество записанных байт и ошибку, если таковая имеется.
func (w gzipWriter) Write(b []byte) (int, error) {
	n, err := w.Writer.Write(b)
	if err != nil {
		return 0, err
	}

	return n, nil
}

// MakeGzipMiddleware создает middleware для gzip сжатия.
// Принимает логгер для записи ошибок.
// Возвращает функцию, которая принимает http.Handler и возвращает http.Handler, обернутый middleware.
// Middleware проверяет заголовок "Content-Encoding" на наличие gzip и распаковывает тело запроса, если необходимо.
// Если клиент поддерживает gzip (заголовок "Accept-Encoding"), то ответ сжимается с использованием gzip.
func MakeGzipMiddleware(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
				gz, err := gzip.NewReader(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				defer func() {
					if closeErr := gz.Close(); closeErr != nil {
						logger.Error("failed to close gzip writer", closeErr)
					}
				}()
				r.Body = gz
			}

			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				_, writeErr := io.WriteString(w, err.Error())
				if writeErr != nil {
					logger.Error("failed to write error message", writeErr)
				}
				return
			}
			defer func() {
				if closeErr := gz.Close(); closeErr != nil {
					logger.Error("failed to close gzip writer", closeErr)
				}
			}()

			w.Header().Set("Content-Encoding", "gzip")
			next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
		})
	}
}
