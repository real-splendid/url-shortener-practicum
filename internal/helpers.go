// Пакет internal предоставляет основные типы и интерфейсы для сервиса сокращения URL.
package internal

import (
	"strconv"
	"time"
)

// MakeKey генерирует уникальный ключ, представляющий собой строковое представление времени в формате base36.
func MakeKey() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
