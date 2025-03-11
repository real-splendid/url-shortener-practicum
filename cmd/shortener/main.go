package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	"go.uber.org/zap"

	"github.com/real-splendid/url-shortener-practicum/internal"
	"github.com/real-splendid/url-shortener-practicum/internal/app"
	"github.com/real-splendid/url-shortener-practicum/internal/storage"
)

// address указывает адрес, на котором будет запущен сервер.
var address *string

// baseURL указывает базовый URL для коротких ссылок.
var baseURL *string

// fileStoragePath указывает путь к файлу для хранения данных (если используется файловое хранилище).
var fileStoragePath *string

// dDSN указывает DSN для подключения к базе данных PostgreSQL.
var dDSN *string

func init() { // coverage-ignore
	address = flag.String("a", ":8080", "server address")
	baseURL = flag.String("b", "http://localhost:8080", "base url")
	fileStoragePath = flag.String("f", "/tmp/short-url-db.json", "file to store results")
	dDSN = flag.String("d", "", "database dsn")
	memprofile := flag.String("memprofile", "", "write memory profile to file")
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

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not create memory profile: %v\n", err)
			os.Exit(1)
		}
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "could not close memory profile: %v\n", err)
			}
		}()
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "could not write memory profile: %v\n", err)
		}
	}
}

func main() { // coverage-ignore
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
	defer func() {
		if err := s.Close(); err != nil {
			logger.Errorf("Failed to close storage: %v", err)
		}
	}()

	app := app.NewApp(s, logger, *baseURL, *dDSN)
	if err := app.Serve(address); err != nil {
		logger.Fatalf("Failed to serve application: %v", err)
	}
}
