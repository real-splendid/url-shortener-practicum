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

type (
	ShortenBatchReq struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	ShortenBatchResp struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}
)

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
		defer r.Body.Close()
		err = json.Unmarshal(body, &req)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var resp []ShortenBatchResp
		for _, item := range req {
			if _, err := url.Parse(item.OriginalURL); err != nil {
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			key := internal.MakeKey()
			if _, err := storage.Set(key, item.OriginalURL, userID); err != nil {
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
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
		w.Write(jsonResp)
	}
}
