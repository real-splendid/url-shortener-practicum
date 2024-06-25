package internal

import (
	"strconv"
	"time"
)

func MakeKey() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
