package main

import (
	"flag"
	"os"

	"go.uber.org/zap"

	"github.com/real-splendid/url-shortener-practicum/internal"
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
	rawLogger, _ := zap.NewDevelopment()
	logger := rawLogger.Sugar()
	var err error
	var s internal.Storage
	if *dDSN != "" {
		s, err = storage.NewPostgresStorage(*dDSN)
		if err != nil {
			logger.Fatalf("Failed to connect to PostgreSQL: %v", err)
		}
	} else if *fileStoragePath != "" {
		s, err = storage.NewFileStorage(*fileStoragePath)
		if err != nil {
			logger.Fatalf("Failed to create file storage: %v", err)
		}
	} else {
		s = storage.NewMemoryStorage()
	}
	defer s.Close()

	app := app.NewApp(s, logger, *baseURL, *dDSN)
	app.Serve(address)
}
