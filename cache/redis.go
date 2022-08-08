package cache

import (
	"context"
	"encoding/json"
	"time"
	"user-service/domain"

	"github.com/go-redis/redis/v8"
)

const (
	expirationTime = 1 * time.Minute
)

// RedisCache .
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache .
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

// Set .
func (c *RedisCache) Set(
	ctx context.Context, key string, val []*domain.User) error {

	bs, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, bs, expirationTime).Err()
}

// Get .
func (c *RedisCache) Get(
	ctx context.Context, key string) ([]*domain.User, bool, error) {

	str, err := c.client.Get(ctx, key).Bytes()
	switch {
	case err == redis.Nil:
		return nil, false, nil

	case err != nil:
		return nil, false, err
	}

	var users []*domain.User
	err = json.Unmarshal([]byte(str), &users)
	if err != nil {
		return nil, false, err
	}

	return users, true, nil
}

// Clear .
func (c *RedisCache) Clear(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
