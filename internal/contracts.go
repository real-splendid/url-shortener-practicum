package internal

import (
	"errors"
)

// Storage определяет интерфейс для работы с хранилищем данных.
type Storage interface {
	// Set сохраняет пару ключ-значение в хранилище.  Возвращает существующий ключ, если значение уже существует, или пустую строку и ошибку в случае ошибки.
	Set(key string, value string, userID string) (string, error)
	// Get возвращает значение по ключу.  Возвращает пустую строку и ошибку, если ключ не найден.
	Get(key string) (string, error)
	// GetUserURLs возвращает список пар короткий URL - оригинальный URL для указанного пользователя.
	GetUserURLs(userID string) ([]URLPair, error)
	// DeleteUserURLs удаляет указанные короткие URL для пользователя.
	DeleteUserURLs(userID string, shortURLs []string) error
	// Close закрывает хранилище.
	Close() error
}

// URLPair представляет пару короткий URL - оригинальный URL.
type URLPair struct {
	// ShortURL короткий URL.
	ShortURL string `json:"short_url"`
	// OriginalURL оригинальный URL.
	OriginalURL string `json:"original_url"`
}

// ErrDuplicateKey возвращается, когда ключ уже существует.
var ErrDuplicateKey = errors.New("duplicate key")

// ErrURLDeleted возвращается, когда URL был удален.
var ErrURLDeleted = errors.New("url deleted")
