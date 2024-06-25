package internal

import (
	"errors"
)

type Storage interface {
	Set(key string, value string) (string, error)
	Get(key string) (string, error)
	Close() error
}

var ErrDuplicateKey = errors.New("duplicate key")
