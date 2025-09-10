package ports

import (
	"context"

	"github.com/SokratisChaimanas/platform-go-challenge/internal/domain"
	"github.com/google/uuid"
)

// UserRepository loads users.
type UserRepository interface {
	// Get returns the user or domain.ErrUserNotFound.
	Get(ctx context.Context, id uuid.UUID) (*domain.User, error)

	// Exists reports whether a user exists.
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}
