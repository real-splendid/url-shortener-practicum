// Package storage предоставляет функциональность для хранения данных.
package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/real-splendid/url-shortener-practicum/internal"
)

// fileRecord представляет запись в файле хранилища.
type fileRecord struct {
	// UUID уникальный идентификатор записи.
	UUID string `json:"uuid"`
	// OriginalURL исходный URL.
	OriginalURL string `json:"original_url"`
	// UserID идентификатор пользователя, создавшего запись.
	UserID string `json:"user_id"`
	// IsDeleted флаг, указывающий, удалена ли запись.
	IsDeleted bool `json:"is_deleted"`
}

type fileStorage struct {
	records  map[string]fileRecord
	userURLs map[string][]string
	file     *os.File
	mu       sync.RWMutex
}

// NewFileStorage создает новое файловое хранилище.
// Принимает путь к файлу хранилища.
// Возвращает указатель на fileStorage и ошибку, если произошла ошибка при открытии файла.
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

// Close закрывает файл хранилища.
// Возвращает ошибку, если произошла ошибка при закрытии файла.
func (s *fileStorage) Close() error {
	return s.file.Close()
}

// Set сохраняет пару ключ-значение в хранилище.
// Принимает ключ, значение и ID пользователя.
// Возвращает ошибку, если произошла ошибка при записи в файл.
// Возвращает internal.ErrDuplicateKey, если ключ уже существует.
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

// Get возвращает значение по ключу.
// Принимает ключ.
// Возвращает значение и ошибку, если ключ не найден или запись удалена.
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

// GetUserURLs возвращает все URL-адреса, принадлежащие пользователю.
// Принимает ID пользователя.
// Возвращает слайс URLPair и ошибку, если произошла ошибка.
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

// DeleteUserURLs удаляет URL-адреса пользователя по коротким ключам.
// Принимает ID пользователя и слайс коротких ключей для удаления.
// Возвращает ошибку, если произошла ошибка при удалении.
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
