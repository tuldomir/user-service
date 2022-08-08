package domain

import (
	"time"

	"github.com/google/uuid"
)

// User .
type User struct {
	ID        uuid.UUID
	Email     string
	CreatedAt time.Time
}
