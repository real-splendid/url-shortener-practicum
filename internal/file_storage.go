package internal

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

func NewFileStorage(path string) FileStorage {
	// FIXME: обрабатывать ошибку
	Logger.Info("Opening file " + path)
	data := make(map[string]string)
	f, _ := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	f.Seek(0, 0)
	scanner := bufio.NewScanner(f)
	var rr ResultRecord

	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &rr)
		if err != nil {
			Logger.Error(err)
			continue
		}
		data[rr.ShortURL] = rr.OriginalURL
	}
	if err := scanner.Err(); err != nil {
		Logger.Error(err)
	}
	return FileStorage{data: data, file: f}
}

func (s FileStorage) Close() error {
	return s.file.Close()
}

func (s FileStorage) Set(key string, value string) {
	rr := ResultRecord{
		UUID:        strconv.Itoa(len(s.data) + 1),
		ShortURL:    key,
		OriginalURL: value,
	}
	// FIXME: обрабатывать ошибку
	jsonBytes, _ := json.Marshal(rr)
	s.file.Write(jsonBytes)
	s.file.Write([]byte("\n"))

	s.data[key] = value
}

func (s FileStorage) Get(key string) (string, error) {
	v, ok := s.data[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return v, nil
}
