package handlers

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/real-splendid/url-shortener-practicum/internal"
	"go.uber.org/zap"
)

func MakeShortenHandler(storage internal.Storage, logger *zap.SugaredLogger, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := internal.MakeKey()
		reqBody, err := readRequestBody(r)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		storage.Set(key, reqBody)
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
