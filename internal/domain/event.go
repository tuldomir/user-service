package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserEvent .
type UserEvent struct {
	EventType string    `json:"event_type,omitempty"`
	UID       uuid.UUID `json:"uid,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
