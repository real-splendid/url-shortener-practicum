package storage

import (
	"database/sql"
	"errors"
	"strings"

	_ "github.com/lib/pq"
	"github.com/real-splendid/url-shortener-practicum/internal"
)

type postgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(dsn string) (*postgresStorage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			short_url TEXT UNIQUE NOT NULL,
			original_url TEXT UNIQUE NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	return &postgresStorage{db: db}, nil
}

func (s *postgresStorage) Close() error {
	return s.db.Close()
}

func (s *postgresStorage) Set(key string, value string) (string, error) {
	_, err := s.db.Exec("INSERT INTO urls(short_url, original_url) VALUES($1, $2)", key, value)
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		var key string
		err = s.db.QueryRow("SELECT short_url FROM urls WHERE original_url = $1", value).Scan(&key)
		if err != nil {
			return "", err
		}
		return key, internal.ErrDuplicateKey
	}
	return "", err
}

func (s *postgresStorage) Get(key string) (string, error) {
	var originalURL string
	err := s.db.QueryRow("SELECT original_url FROM urls WHERE short_url = $1", key).Scan(&originalURL)
	if err == sql.ErrNoRows {
		return "", errors.New("key not found")
	}
	if err != nil {
		return "", err
	}
	return originalURL, nil
}
