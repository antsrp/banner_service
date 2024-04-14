package cache

type Storager[T any] interface {
	Set(key string, value T) error
	Get(key string) (T, error)
	Delete(key string) error
	Close() error
}
