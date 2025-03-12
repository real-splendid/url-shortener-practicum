// Пакет handlers предоставляет обработчики HTTP-запросов для сервиса сокращения URL.
package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/real-splendid/url-shortener-practicum/internal"
	"github.com/real-splendid/url-shortener-practicum/internal/middleware"
)

// MakeDeleteUserURLsHandler создает обработчик для удаления URL пользователя.
// Принимает хранилище и логгер.
// Возвращает функцию обработчика.
func MakeDeleteUserURLsHandler(storage internal.Storage, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		body, readErr := io.ReadAll(r.Body)
		if readErr != nil {
			logger.Error(readErr)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		defer func() {
			if closeErr := r.Body.Close(); closeErr != nil {
				logger.Error("failed to close request body", closeErr)
			}
		}()

		var shortURLs []string
		if unmarshalErr := json.Unmarshal(body, &shortURLs); unmarshalErr != nil {
			logger.Error(unmarshalErr)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		go func() {
			if deleteErr := storage.DeleteUserURLs(userID, shortURLs); deleteErr != nil {
				logger.Error(deleteErr)
			}
		}()

		w.WriteHeader(http.StatusAccepted)
	}
}
