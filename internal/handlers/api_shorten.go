package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/real-splendid/url-shortener-practicum/internal"
	"go.uber.org/zap"
)

type (
	ShortenReq struct {
		URL string `json:"url"`
	}

	ShortenResp struct {
		Result string `json:"result"`
	}
)

func MakeAPIShortenHandler(storage internal.Storage, logger *zap.SugaredLogger, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := internal.MakeKey()
		URL, err := readURLFromAPIRequestBody(r)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		storage.Set(key, URL)
		shortURL := baseURL + "/" + key
		jsonResp, err := json.Marshal(ShortenResp{Result: shortURL})
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResp)
	}
}

func readURLFromAPIRequestBody(r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	var req ShortenReq
	if err = json.Unmarshal(body, &req); err != nil {
		return "", err
	}

	if _, err = url.Parse(req.URL); err != nil {
		return "", err
	}
	return req.URL, nil
}
