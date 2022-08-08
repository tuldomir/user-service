package repo

import (
	"context"
	"user-service/domain"

	"github.com/google/uuid"
)

// DB .
type DB interface {
	Add(ctx context.Context, id uuid.UUID, email string) (*domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*domain.User, error)
}
