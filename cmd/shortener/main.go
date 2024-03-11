package main

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

const HOST = "http://localhost"
const PORT = "8080"

var shortURLStorage map[string]string

func readRequestBody(r *http.Request) string {
	body, err := io.ReadAll(r.Body)
	// FIXME: better error handle
	if err != nil {
		return ""
	}
	defer r.Body.Close()
	return string(body)
}

// FIXME: better key
func makeKey() string {
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	return timestamp[len(timestamp)-8:]
}

func handleShortLinkCreation(w http.ResponseWriter, r *http.Request) {
	key := makeKey()
	shortURLStorage[key] = readRequestBody(r)
	w.WriteHeader(http.StatusCreated)
	shortURL := HOST + ":" + PORT + "/" + key
	w.Write([]byte(shortURL))
}

func handleRedirection(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	originalURL, ok := shortURLStorage[key]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	shortURLStorage = make(map[string]string)
	r := chi.NewRouter()
	r.Post("/", handleShortLinkCreation)
	r.Get("/{key}", handleRedirection)
	http.ListenAndServe(":"+PORT, r)
}
