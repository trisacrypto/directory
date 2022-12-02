package cache

import "errors"

var (
	ErrKeyNotFound  = errors.New("specified key not found")
	ErrValueExpired = errors.New("value specified by key has expired")
)
