package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"

	_ "github.com/lib/pq"

	"github.com/real-splendid/url-shortener-practicum/internal"
)

type fileRecord struct {
	UUID        string `json:"uuid"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
}

type fileStorage struct {
	data     map[string]string
	userURLs map[string][]internal.URLPair
	file     *os.File
}

func NewFileStorage(path string) (*fileStorage, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	fs := &fileStorage{
		data:     make(map[string]string),
		userURLs: make(map[string][]internal.URLPair),
		file:     f,
	}

	scanner := bufio.NewScanner(f)
	var rr fileRecord
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &rr); err != nil {
			return nil, err
		}
		fs.data[rr.UUID] = rr.OriginalURL
		fs.userURLs[rr.UserID] = append(fs.userURLs[rr.UserID], internal.URLPair{
			ShortURL:    rr.UUID,
			OriginalURL: rr.OriginalURL,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return fs, nil
}

func (s *fileStorage) Close() error {
	return s.file.Close()
}

func (s *fileStorage) Set(key string, value string, userID string) (string, error) {
	if existingKey, exists := s.getKeyByValue(value); exists {
		return existingKey, internal.ErrDuplicateKey
	}

	rr := fileRecord{
		UUID:        key,
		OriginalURL: value,
		UserID:      userID,
	}

	jsonData, err := json.Marshal(rr)
	if err != nil {
		return "", err
	}

	if _, err := s.file.Write(append(jsonData, '\n')); err != nil {
		return "", err
	}

	s.data[key] = value
	s.userURLs[userID] = append(s.userURLs[userID], internal.URLPair{
		ShortURL:    key,
		OriginalURL: value,
	})

	return "", nil
}

func (s *fileStorage) Get(key string) (string, error) {
	v, ok := s.data[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return v, nil
}

func (s *fileStorage) GetUserURLs(userID string) ([]internal.URLPair, error) {
	return s.userURLs[userID], nil
}

func (s *fileStorage) getKeyByValue(value string) (string, bool) {
	for k, v := range s.data {
		if v == value {
			return k, true
		}
	}
	return "", false
}
