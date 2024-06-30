package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/real-splendid/url-shortener-practicum/internal"
)

type fileRecord struct {
	UUID        string `json:"uuid"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
	IsDeleted   bool   `json:"is_deleted"`
}

type fileStorage struct {
	records  map[string]fileRecord
	userURLs map[string][]string
	file     *os.File
	mu       sync.RWMutex
}

func NewFileStorage(path string) (*fileStorage, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	s := &fileStorage{
		records:  make(map[string]fileRecord),
		userURLs: make(map[string][]string),
		file:     f,
	}

	scanner := bufio.NewScanner(f)
	var record fileRecord
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return nil, err
		}

		s.records[record.UUID] = record
		if len(s.userURLs[record.UserID]) == 0 {
			s.userURLs[record.UserID] = make([]string, 1)
		}
		s.userURLs[record.UserID] = append(s.userURLs[record.UserID], record.UUID)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *fileStorage) Close() error {
	return s.file.Close()
}

func (s *fileStorage) Set(key string, value string, userID string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existingKey, err := s.getKeyByValue(value)
	if err == nil {
		return existingKey, internal.ErrDuplicateKey
	}

	record := fileRecord{
		UUID:        key,
		OriginalURL: value,
		UserID:      userID,
		IsDeleted:   false,
	}

	jsonData, err := json.Marshal(record)
	if err != nil {
		return "", err
	}

	_, err = s.file.Write(append(jsonData, '\n'))
	if err != nil {
		return "", err
	}

	s.records[key] = fileRecord{
		UUID:        key,
		OriginalURL: value,
		UserID:      userID,
		IsDeleted:   false,
	}
	if len(s.userURLs[record.UserID]) == 0 {
		s.userURLs[record.UserID] = make([]string, 1)
	}
	s.userURLs[userID] = append(s.userURLs[userID], key)

	return "", nil
}

func (s *fileStorage) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	record, ok := s.records[key]
	if !ok {
		return "", errors.New("key not found")
	}

	if record.IsDeleted {
		return "", internal.ErrURLDeleted
	}

	return record.OriginalURL, nil
}

func (s *fileStorage) GetUserURLs(userID string) ([]internal.URLPair, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	urls := make([]internal.URLPair, 0)
	for _, shortKey := range s.userURLs[userID] {
		record, ok := s.records[shortKey]
		if !ok {
			return urls, errors.New("corrupted url key")
		}
		if !record.IsDeleted {
			urls = append(urls, internal.URLPair{ShortURL: shortKey, OriginalURL: record.OriginalURL})
		}
	}
	return urls, nil
}

func (s *fileStorage) DeleteUserURLs(userID string, shortKeysToDelete []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, shortKey := range shortKeysToDelete {
		if _, ok := s.records[shortKey]; !ok {
			return errors.New("record not found by short key")
		}

		s.records[shortKey] = fileRecord{
			UUID:        shortKey,
			OriginalURL: s.records[shortKey].OriginalURL,
			UserID:      s.records[shortKey].UserID,
			IsDeleted:   true,
		}

		jsonData, err := json.Marshal(s.records[shortKey])
		if err != nil {
			return err
		}

		if _, err := s.file.Write(append(jsonData, '\n')); err != nil {
			return err
		}
	}

	return nil
}

func (s *fileStorage) getKeyByValue(value string) (string, error) {
	for _, v := range s.records {
		if v.OriginalURL == value && !v.IsDeleted {
			return v.UUID, nil
		}
	}
	return "", errors.New("key not found")
}
