package internal

type Storage interface {
	Set(key string, value string) error
	Get(key string) (string, error)
}
