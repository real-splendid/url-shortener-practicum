package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/real-splendid/url-shortener-practicum/internal"
)

func main() {
	defer internal.Logger.Sync()

	address := flag.String("a", ":8080", "server address")
	baseURL := flag.String("b", internal.DefaultBaseURL, "base url")
	fileStoragePath := flag.String("f", "/tmp/short-url-db.json", "file to store results")
	flag.Parse()
	if envAddress, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		*address = envAddress
	}
	if envBaseURL, ok := os.LookupEnv("BASE_URL"); ok {
		*baseURL = envBaseURL
	}
	if envFileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		*fileStoragePath = envFileStoragePath
	}
	internal.SetBaseURL(*baseURL)
	fStorage := internal.NewFileStorage(*fileStoragePath)
	defer fStorage.Close()
	internal.SetStorage(fStorage)

	r := chi.NewRouter()
	r.Use(internal.LoggingMiddleware)
	r.Use(internal.GzipMiddleware)
	r.Post("/", internal.HandleShorten)
	r.Post("/api/shorten", internal.HandleAPIShorten)
	r.Get("/{key}", internal.HandleRedirection)
	if err := http.ListenAndServe(*address, r); err != nil {
		internal.Logger.Error(err)
		os.Exit(1)
	}
}
