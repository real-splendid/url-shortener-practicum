package internal

import (
	"errors"
)

type Storage interface {
	Set(key string, value string, userID string) (string, error)
	Get(key string) (string, error)
	GetUserURLs(userID string) ([]URLPair, error)
	DeleteUserURLs(userID string, shortURLs []string) error
	Close() error
}

type URLPair struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

var ErrDuplicateKey = errors.New("duplicate key")
var ErrURLDeleted = errors.New("url deleted")
