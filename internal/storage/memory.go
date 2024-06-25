package storage

import (
	"errors"
)

type memoryStorage struct {
	data map[string]string
}

func NewMemoryStorage() *memoryStorage {
	return &memoryStorage{
		data: make(map[string]string),
	}
}

func (s *memoryStorage) Close() error {
	return nil
}

func (s *memoryStorage) Set(key string, value string) error {
	s.data[key] = value
	return nil
}

func (s *memoryStorage) Get(key string) (string, error) {
	v, ok := s.data[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return v, nil
}
