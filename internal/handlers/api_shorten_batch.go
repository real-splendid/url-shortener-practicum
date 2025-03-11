// Пакет handlers предоставляет обработчики HTTP-запросов для сервиса сокращения URL.
package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"go.uber.org/zap"

	"github.com/real-splendid/url-shortener-practicum/internal"
	"github.com/real-splendid/url-shortener-practicum/internal/middleware"
)

// ShortenBatchReq структура запроса для пакетного сокращения URL.
type ShortenBatchReq struct {
	// CorrelationID идентификатор корреляции.
	CorrelationID string `json:"correlation_id"`
	// OriginalURL оригинальный URL.
	OriginalURL string `json:"original_url"`
}

// ShortenBatchResp структура ответа для пакетного сокращения URL.
type ShortenBatchResp struct {
	// CorrelationID идентификатор корреляции.
	CorrelationID string `json:"correlation_id"`
	// ShortURL сокращенный URL.
	ShortURL string `json:"short_url"`
}

// MakeAPIShortenBatchHandler создает обработчик для пакетного сокращения URL через API.
// Принимает хранилище, логгер и базовый URL.
// Возвращает функцию обработчика.
func MakeAPIShortenBatchHandler(storage internal.Storage, logger *zap.SugaredLogger, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req []ShortenBatchReq
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer func() {
			if closeErr := r.Body.Close(); closeErr != nil {
				logger.Error("failed to close request body", closeErr)
			}
		}()
		err = json.Unmarshal(body, &req)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var resp []ShortenBatchResp
		for _, item := range req {
			if _, parseErr := url.Parse(item.OriginalURL); parseErr != nil {
				logger.Error(parseErr)
				http.Error(w, parseErr.Error(), http.StatusBadRequest)
				return
			}

			key := internal.MakeKey()
			if _, setErr := storage.Set(key, item.OriginalURL, userID); setErr != nil {
				logger.Error(setErr)
				http.Error(w, setErr.Error(), http.StatusInternalServerError)
				return
			}

			resp = append(resp, ShortenBatchResp{
				CorrelationID: item.CorrelationID,
				ShortURL:      baseURL + "/" + key,
			})
		}

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
