// Пакет handlers предоставляет обработчики HTTP-запросов для сервиса сокращения URL.
package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"go.uber.org/zap"

	"github.com/real-splendid/url-shortener-practicum/internal"
	"github.com/real-splendid/url-shortener-practicum/internal/middleware"
)

// ShortenReq структура запроса для сокращения URL.
type ShortenReq struct {
	// URL URL для сокращения.
	URL string `json:"url"`
}

// ShortenResp структура ответа для сокращения URL.
type ShortenResp struct {
	// Result сокращенный URL.
	Result string `json:"result"`
}

// MakeAPIShortenHandler создает обработчик для сокращения URL через API.
// Принимает хранилище, логгер и базовый URL.
// Возвращает функцию обработчика.
func MakeAPIShortenHandler(storage internal.Storage, logger *zap.SugaredLogger, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		key := internal.MakeKey()
		URL, err := readURLFromAPIRequestBody(r)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		duplicateKey, setErr := storage.Set(key, URL, userID)
		if errors.Is(setErr, internal.ErrDuplicateKey) {
			logger.Info("key already exists", duplicateKey)
			shortURL := baseURL + "/" + duplicateKey
			jsonResp, marshalErr := json.Marshal(ShortenResp{Result: shortURL})
			if marshalErr != nil {
				logger.Error(marshalErr)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			_, writeErr := w.Write(jsonResp)
			if writeErr != nil {
				logger.Error("failed to write response", writeErr)
			}
			return
		}
		shortURL := baseURL + "/" + key
		jsonResp, err := json.Marshal(ShortenResp{Result: shortURL})
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, writeErr := w.Write(jsonResp)
		if writeErr != nil {
			logger.Error("failed to write response", writeErr)
		}
	}
}

// readURLFromAPIRequestBody читает URL из тела запроса API.
// Принимает запрос HTTP.
// Возвращает URL и ошибку, если произошла ошибка при чтении или парсинге URL.
func readURLFromAPIRequestBody(r *http.Request) (string, error) {
	limitedReader := io.LimitReader(r.Body, maxBodySize)
	body, err := io.ReadAll(limitedReader)

	if err != nil {
		return "", err
	}

	defer func() {
		_ = r.Body.Close()
	}()

	var req ShortenReq
	if err = json.Unmarshal(body, &req); err != nil {
		return "", err
	}

	if _, err = url.Parse(req.URL); err != nil {
		return "", err
	}
	return req.URL, nil
}
