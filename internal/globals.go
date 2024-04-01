package internal

import (
	"go.uber.org/zap"
)

const (
	DefaultBaseURL = "http://localhost:8080"
)

type Storage interface {
	Set(key string, value string)
	Get(key string) (string, error)
}

var (
	BaseURL         *string
	Logger          *zap.SugaredLogger
	internalStorage Storage
)

func init() {
	rawLogger, _ := zap.NewDevelopment()
	Logger = rawLogger.Sugar()

	baseURL := DefaultBaseURL
	BaseURL = &baseURL
}

func SetBaseURL(url string) {
	BaseURL = &url
}

func SetStorage(storage Storage) {
	internalStorage = storage
}
