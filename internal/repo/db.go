package repo

import (
	"context"
	"user-service/models"
)

// DB .
type DB interface {
	Add(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, uid string) error
	List(ctx context.Context) ([]*models.User, error)
}
