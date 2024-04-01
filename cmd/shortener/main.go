package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/real-splendid/url-shortener-practicum/internal"
	"go.uber.org/zap"
)

// TODO: добавить логирование ошибокок

var (
	logger *zap.SugaredLogger
)

func init() {
	rawLogger, _ := zap.NewDevelopment()
	logger = rawLogger.Sugar()
}

func main() {
	address := flag.String("a", ":8080", "server address")
	baseURL := flag.String("b", internal.DefaultBaseURL, "base url")
	flag.Parse()
	if envAddress, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		*address = envAddress
	}
	if envBaseURL, ok := os.LookupEnv("BASE_URL"); ok {
		*baseURL = envBaseURL
	}
	internal.BaseURL = baseURL
	internal.Logger = logger
	r := chi.NewRouter()
	r.Use(internal.LoggingMiddleware)
	r.Use(internal.GzipMiddleware)
	r.Post("/", internal.HandleShorten)
	r.Post("/api/shorten", internal.HandleAPIShorten)
	r.Get("/{key}", internal.HandleRedirection)
	if err := http.ListenAndServe(*address, r); err != nil {
		panic(err)
	}
	logger.Sync()
}
