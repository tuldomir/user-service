package cache

import (
	"context"
	"user-service/models"
)

// Cache .
type Cache interface {
	Set(ctx context.Context, key string, users []*models.User) error
	Get(ctx context.Context, key string) ([]*models.User, bool, error)
	Clear(ctx context.Context, key string) error
}
