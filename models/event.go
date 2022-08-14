package models

import (
	"time"
)

// UserEvent .
type UserEvent struct {
	EventType string    `json:"event_type,omitempty"`
	UID       string    `json:"uid,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
