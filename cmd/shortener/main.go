package main

import (
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var (
	storage = make(map[string]string)
	address = flag.String("a", ":8080", "server address")
	baseURL = flag.String("b", "http://localhost:8080", "base url")
	logger  *zap.SugaredLogger
)

type (
	loggingResponseWriter struct {
		http.ResponseWriter
		status int
		size   int
	}

	ShortenReq struct {
		URL string `json:"url"`
	}

	ShortenResp struct {
		Result string `json:"result"`
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.status = statusCode
}

func init() {
	rawLogger, _ := zap.NewDevelopment()
	logger = rawLogger.Sugar()
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

func makeKey() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
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

func handleAPIShorten(w http.ResponseWriter, r *http.Request) {
	key := makeKey()
	URL, err := readURLFromAPIRequestBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	storage[key] = URL
	shortURL := *baseURL + "/" + key
	jsonResp, _ := json.Marshal(ShortenResp{Result: shortURL})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResp)
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

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := &loggingResponseWriter{ResponseWriter: w}
		next.ServeHTTP(lw, r)
		logger.Infow("request processed",
			"uri", r.RequestURI,
			"method", r.Method,
			"duration", time.Since(start),
			"status", lw.status,
			"size", lw.size,
		)
	})
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
	r.Use(loggingMiddleware)
	r.Post("/", handleShorten)
	r.Post("/api/shorten", handleAPIShorten)
	r.Get("/{key}", handleRedirection)
	if err := http.ListenAndServe(*address, r); err != nil {
		panic(err)
	}
	logger.Sync()
}
