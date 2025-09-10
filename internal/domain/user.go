package domain

import (
	"time"

	"github.com/google/uuid"
)

// User is the domain representation of a user.
type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
}
