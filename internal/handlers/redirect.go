package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/real-splendid/url-shortener-practicum/internal"
	"go.uber.org/zap"
)

func MakeRedirectionHandler(storage internal.Storage, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		originalURL, err := storage.Get(key)
		if err != nil {
			logger.Infof("key %s not found", key)
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
