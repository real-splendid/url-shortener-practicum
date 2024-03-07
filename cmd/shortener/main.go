package main

import (
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var shortURLStorage map[string]string

func main() {
	shortURLStorage = make(map[string]string)
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handleRequest)

	if err := http.ListenAndServe(`:8080`, mux); err != nil {
		panic(err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.URL.Path == "/" {
		key := makeKey()
		shortURLStorage[key] = readRequestBody(r)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://localhost:8080/" + key))
		return
	}

	re := regexp.MustCompile(`^/\w{8}$`)
	if r.Method == http.MethodGet && re.MatchString(r.URL.Path) {
		w.Header().Set("Location", shortURLStorage[r.URL.Path[1:]])
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

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
