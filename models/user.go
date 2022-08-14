package models

import (
	"time"
)

// User .
type User struct {
	UID       string    `json:"id,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
