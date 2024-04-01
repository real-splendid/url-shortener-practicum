package internal

import (
	"errors"
)

type MemoryStorage map[string]string

func NewMemoryStorage() MemoryStorage {
	return make(map[string]string)
}

func (s MemoryStorage) Set(key string, value string) {
	s[key] = value
}

func (s MemoryStorage) Get(key string) (string, error) {
	v, ok := s[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return v, nil
}
