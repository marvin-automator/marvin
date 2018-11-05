package db

import "fmt"

type KeyNotFoundError struct {
	key string
}

func (e KeyNotFoundError) Error() string {
	return fmt.Sprint("Key not found", e.key)
}
