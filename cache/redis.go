package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

func (cache *RedisCache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return cache.client.Set(ctx, key, value, expiration).Err()
}

func (cache *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return val, nil
}

func (cache *RedisCache) Delete(ctx context.Context, keys ...string) error {
	return cache.client.Del(ctx, keys...).Err()
}
