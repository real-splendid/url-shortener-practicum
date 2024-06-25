package main

import (
	"flag"
	"os"

	"github.com/real-splendid/url-shortener-practicum/internal/app"
	"github.com/real-splendid/url-shortener-practicum/internal/storage"
)

var (
	address         *string
	baseURL         *string
	fileStoragePath *string
	dDSN            *string
)

func init() {
	address = flag.String("a", ":8080", "server address")
	baseURL = flag.String("b", "http://localhost:8080", "base url")
	fileStoragePath = flag.String("f", "/tmp/short-url-db.json", "file to store results")
	dDSN = flag.String("d", "", "database dsn")
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
	if envDSN, ok := os.LookupEnv("DATABASE_DSN"); ok {
		*dDSN = envDSN
	}
}

func main() {
	s, _ := storage.NewFileStorage(*fileStoragePath)
	defer s.Close()

	app := app.NewApp(s, *baseURL, *dDSN)
	app.Serve(address)
}
