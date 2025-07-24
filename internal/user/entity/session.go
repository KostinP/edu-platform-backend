package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserSession struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Token        string
	UserAgent    string
	IPAddress    string
	Country      string
	City         string
	CreatedAt    time.Time
	LastActiveAt time.Time
	ExpiresAt    time.Time
}
