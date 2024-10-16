package cache

import (
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, keys ...string) error
}

var ErrNotFound = NotFoundError{}

type NotFoundError struct{}

func (e NotFoundError) Error() string {
	return "Key not found"
}
