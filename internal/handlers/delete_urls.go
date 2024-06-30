package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/real-splendid/url-shortener-practicum/internal"
	"github.com/real-splendid/url-shortener-practicum/internal/middleware"
)

func MakeDeleteUserURLsHandler(storage internal.Storage, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error(err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var shortURLs []string
		if err := json.Unmarshal(body, &shortURLs); err != nil {
			logger.Error(err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		go func() {
			if err := storage.DeleteUserURLs(userID, shortURLs); err != nil {
				logger.Error(err)
			}
		}()

		w.WriteHeader(http.StatusAccepted)
	}
}
