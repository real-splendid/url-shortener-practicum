package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/real-splendid/url-shortener-practicum/internal"
	"github.com/real-splendid/url-shortener-practicum/internal/middleware"
)

func MakeShortenHandler(storage internal.Storage, logger *zap.SugaredLogger, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		key := internal.MakeKey()
		reqBody, err := readRequestBody(r)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		duplicateKey, err := storage.Set(key, reqBody, userID)
		if errors.Is(err, internal.ErrDuplicateKey) {
			logger.Info("key already exists", duplicateKey)
			shortURL := baseURL + "/" + duplicateKey
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(shortURL))
			return
		}
		shortURL := baseURL + "/" + key
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortURL))
	}
}

func makeKey() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}

func readRequestBody(r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	if _, err = url.Parse(string(body)); err != nil {
		return "", err
	}
	return string(body), nil
}
