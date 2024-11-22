// Package storage предоставляет функциональность для хранения данных.
package storage

import (
	"errors"

	"github.com/real-splendid/url-shortener-practicum/internal"
)

// memoryStorage реализует хранилище в памяти.
type memoryStorage struct {
	// data хранит пары ключ-значение.
	data map[string]string
	// userURLs хранит URL-адреса, связанные с пользователями.
	userURLs map[string][]internal.URLPair
	// deleted хранит информацию об удаленных URL-адресах.
	deleted map[string]bool
}

// NewMemoryStorage создает новое хранилище в памяти.
// Возвращает указатель на memoryStorage.
func NewMemoryStorage() *memoryStorage {
	return &memoryStorage{
		data:     make(map[string]string),
		userURLs: make(map[string][]internal.URLPair),
		deleted:  make(map[string]bool),
	}
}

// Close закрывает хранилище (для memoryStorage не выполняет никаких действий).
// Возвращает nil.
func (s *memoryStorage) Close() error {
	return nil
}

// Set сохраняет пару ключ-значение в хранилище.
// Принимает ключ, значение и ID пользователя.
// Возвращает пустую строку и nil, если запись прошла успешно.
func (s *memoryStorage) Set(key string, value string, userID string) (string, error) {
	s.data[key] = value
	s.userURLs[userID] = append(s.userURLs[userID], internal.URLPair{ShortURL: key, OriginalURL: value})
	return "", nil
}

// Get возвращает значение по ключу.
// Принимает ключ.
// Возвращает значение и nil, если ключ найден.
// Возвращает пустую строку и ошибку "key not found", если ключ не найден.
// Возвращает пустую строку и internal.ErrURLDeleted, если URL был удален.
func (s *memoryStorage) Get(key string) (string, error) {
	if s.deleted[key] {
		return "", internal.ErrURLDeleted
	}

	v, ok := s.data[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return v, nil
}

// GetUserURLs возвращает все URL-адреса, принадлежащие пользователю.
// Принимает ID пользователя.
// Возвращает слайс URLPair и nil.
func (s *memoryStorage) GetUserURLs(userID string) ([]internal.URLPair, error) {
	urls := make([]internal.URLPair, 0)
	for _, pair := range s.userURLs[userID] {
		if !s.deleted[pair.ShortURL] {
			urls = append(urls, pair)
		}
	}
	return urls, nil
}

// DeleteUserURLs удаляет URL-адреса пользователя по коротким URL.
// Принимает ID пользователя и слайс коротких URL для удаления.
// Возвращает nil.
func (s *memoryStorage) DeleteUserURLs(userID string, shortURLs []string) error {
	for _, shortURL := range shortURLs {
		s.deleted[shortURL] = true
	}
	return nil
}
