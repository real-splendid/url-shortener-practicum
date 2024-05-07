package storage

import (
	"database/sql"
	"errors"

	// "strings"

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
			original_url TEXT NOT NULL,
			user_id TEXT NOT NULL
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

func (s *postgresStorage) Set(key string, value string, userID string) (string, error) {
	var existingKey string
	err := s.db.QueryRow("SELECT short_url FROM urls WHERE original_url = $1", value).Scan(&existingKey)
	if err == nil {
		return existingKey, internal.ErrDuplicateKey
	} else if err != sql.ErrNoRows {
		return "", err
	}

	_, err = s.db.Exec("INSERT INTO urls(short_url, original_url, user_id) VALUES($1, $2, $3)", key, value, userID)
	if err != nil {
		return "", err
	}

	return "", nil
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

func (s *postgresStorage) GetUserURLs(userID string) ([]internal.URLPair, error) {
	rows, err := s.db.Query("SELECT short_url, original_url FROM urls WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []internal.URLPair
	for rows.Next() {
		var pair internal.URLPair
		if err := rows.Scan(&pair.ShortURL, &pair.OriginalURL); err != nil {
			return nil, err
		}
		urls = append(urls, pair)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}
