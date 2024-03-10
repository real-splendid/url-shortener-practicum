package main

import (
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"
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
	originalURL, ok := shortURLStorage[r.URL.Path[1:]]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	re := regexp.MustCompile(`^/\w{8}$`)
	if r.Method == http.MethodGet && re.MatchString(r.URL.Path) {
		handleRedirection(w, r)
		return
	}
	if r.Method == http.MethodPost && r.URL.Path == "/" {
		handleShortLinkCreation(w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func main() {
	shortURLStorage = make(map[string]string)
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRequest)
	http.ListenAndServe(":"+PORT, mux)
}
