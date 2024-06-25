package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"strconv"
)

type ResultRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type FileStorage struct {
	data map[string]string
	file *os.File
}

func NewFileStorage(path string) (*FileStorage, error) {
	data := make(map[string]string)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return &FileStorage{}, err
	}
	f.Seek(0, 0)
	scanner := bufio.NewScanner(f)
	var rr ResultRecord

	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &rr)
		if err != nil {
			return &FileStorage{}, err
		}
		data[rr.ShortURL] = rr.OriginalURL
	}
	if err := scanner.Err(); err != nil {
		return &FileStorage{}, err
	}
	return &FileStorage{data: data, file: f}, nil
}

func (s *FileStorage) Close() error {
	return s.file.Close()
}

func (s *FileStorage) Set(key string, value string) error {
	rr := ResultRecord{
		UUID:        strconv.Itoa(len(s.data) + 1),
		ShortURL:    key,
		OriginalURL: value,
	}
	jsonBytes, err := json.Marshal(rr)
	if err != nil {
		return err
	}
	s.file.Write(jsonBytes)
	s.file.Write([]byte("\n"))

	s.data[key] = value
	return nil
}

func (s *FileStorage) Get(key string) (string, error) {
	v, ok := s.data[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return v, nil
}
