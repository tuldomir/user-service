package domain

import (
	"time"

	"github.com/google/uuid"
)

// User .
type User struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
