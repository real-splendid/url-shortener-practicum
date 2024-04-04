package internal

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type (
	ShortenReq struct {
		URL string `json:"url"`
	}

	ShortenResp struct {
		Result string `json:"result"`
	}
)

func HandleShorten(w http.ResponseWriter, r *http.Request) {
	key := makeKey()
	reqBody, err := readRequestBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	internalStorage.Set(key, reqBody)
	shortURL := *BaseURL + "/" + key
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func HandleAPIShorten(w http.ResponseWriter, r *http.Request) {
	key := makeKey()
	URL, err := readURLFromAPIRequestBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	internalStorage.Set(key, URL)
	shortURL := *BaseURL + "/" + key
	jsonResp, _ := json.Marshal(ShortenResp{Result: shortURL})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResp)
}

func HandleRedirection(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	originalURL, err := internalStorage.Get(key)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
