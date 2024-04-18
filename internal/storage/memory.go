package storage

import (
	"errors"
)

type memoryStorage map[string]string

func NewMemoryStorage() memoryStorage {
	return make(map[string]string)
}

func (s memoryStorage) Set(key string, value string) error {
	s[key] = value
	return nil
}

func (s memoryStorage) Get(key string) (string, error) {
	v, ok := s[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return v, nil
}
