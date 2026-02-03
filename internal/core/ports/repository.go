package ports

import (
	"context"
	"iter"

	"go-favorites-app/internal/core/domain/auth"
	"go-favorites-app/internal/core/domain/favorites"
)

// UserRepository defines storage for users.
type UserRepository interface {
	Save(ctx context.Context, user auth.User) error
	FindByEmail(ctx context.Context, email string) (auth.User, error)
}

// FavoriteRepository defines the interface for favorite asset storage.
type FavoriteRepository interface {
	// Save persists a generic Asset.
	Save(ctx context.Context, asset favorites.Asset) error

	// FindByID retrieves an asset by its ID.
	FindByID(ctx context.Context, id string) (favorites.Asset, error)

	// FindAll returns an iterator of Assets to stream results.
	// limit and offset determine pagination.
	FindAll(ctx context.Context, limit, offset int) (iter.Seq2[favorites.Asset, error], error)

	// FindByUser returns an iterator of Assets for a specific user.
	FindByUser(ctx context.Context, userID string, limit, offset int) (iter.Seq2[favorites.Asset, error], error)

	// Delete removes an asset by ID.
	Delete(ctx context.Context, id string) error

	// UpdateDescription updates just the description of an asset.
	UpdateDescription(ctx context.Context, id, description string) (favorites.Asset, error)
}
