package cache

import (
	"context"
	"user-service/internal/domain"
)

// Cache .
type Cache interface {
	Set(ctx context.Context, key string, users []*domain.User) error
	Get(ctx context.Context, key string) ([]*domain.User, bool, error)
	Clear(ctx context.Context, key string) error
}
