package storage

import (
	"errors"

	"github.com/real-splendid/url-shortener-practicum/internal"
)

type memoryStorage struct {
	data     map[string]string
	userURLs map[string][]internal.URLPair
}

func NewMemoryStorage() *memoryStorage {
	return &memoryStorage{
		data:     make(map[string]string),
		userURLs: make(map[string][]internal.URLPair),
	}
}

func (s *memoryStorage) Close() error {
	return nil
}

func (s *memoryStorage) Set(key string, value string, userID string) (string, error) {
	s.data[key] = value
	s.userURLs[userID] = append(s.userURLs[userID], internal.URLPair{ShortURL: key, OriginalURL: value})
	return "", nil
}

func (s *memoryStorage) Get(key string) (string, error) {
	v, ok := s.data[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return v, nil
}

func (s *memoryStorage) GetUserURLs(userID string) ([]internal.URLPair, error) {
	return s.userURLs[userID], nil
}
