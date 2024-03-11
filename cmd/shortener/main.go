package main

import (
	"flag"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

var (
	storage = make(map[string]string)
	address = flag.String("a", ":8080", "server address")
	baseURL = flag.String("b", "http://localhost:8080", "base url")
)

func readRequestBody(r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	if _, err = url.Parse(string(body)); err != nil {
		panic("invalid url")
	}
	return string(body), nil
}

func makeKey() string {
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	return timestamp[len(timestamp)-8:]
}

func handleShortLinkCreation(w http.ResponseWriter, r *http.Request) {
	key := makeKey()
	reqBody, err := readRequestBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	storage[key] = reqBody
	shortURL := *baseURL + "/" + key
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func handleRedirection(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	originalURL, ok := storage[key]
	if !ok {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	flag.Parse()
	if envAddress, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		*address = envAddress
	}
	if envBaseURL, ok := os.LookupEnv("BASE_URL"); ok {
		*baseURL = envBaseURL
	}
	r := chi.NewRouter()
	r.Post("/", handleShortLinkCreation)
	r.Get("/{key}", handleRedirection)
	err := http.ListenAndServe(*address, r)
	if err != nil {
		panic(err)
	}
}
