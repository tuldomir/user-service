package repo

import (
	"context"
	"user-service/models"

	"github.com/google/uuid"
)

// DB .
type DB interface {
	Add(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*models.User, error)
}
