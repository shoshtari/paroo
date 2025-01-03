package cache

import "errors"

// errors
var (
	NotfoundError = errors.New("key not found")
)

type Cache[KeyType comparable, ValType any] interface {
	Set(KeyType, ValType) error
	Get(KeyType) (ValType, error)
	Exists(KeyType) (bool, error)
	Delete(KeyType) error
}
