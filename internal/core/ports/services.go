package ports

import (
	"context"
	"iter"

	"go-favorites-app/internal/core/domain/favorites"
)

// AuthService defines the authentication service.
type AuthService interface {
	SignUp(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (token string, err error)
}

// Enricher defines an external service that enriches assets.
type Enricher interface {
	Enrich(ctx context.Context, asset favorites.Asset) error
}

// Cache defines the caching operations.
// We keep it simple and tailored to our needs.
type Cache interface {
	// AddToSet adds an asset ID with a score (timestamp) to the sorted set.
	AddToSet(ctx context.Context, id string, score float64) error

	// Set holds the asset data.
	Set(ctx context.Context, id string, data []byte) error

	// GetBatch retrieves multiple assets by ID.
	GetBatch(ctx context.Context, ids []string) (map[string][]byte, error)

	// GetIdsFromSet returns IDs from the sorted set for a range.
	GetIdsFromSet(ctx context.Context, start, stop int64) ([]string, error)

	// Remove removes an asset from cache.
	Remove(ctx context.Context, id string) error

	// Invalidate removes only the asset data, keeping the ID in the set.
	Invalidate(ctx context.Context, id string) error
}

// FavoriteService defines the application logic.
type FavoriteService interface {
	Save(ctx context.Context, asset favorites.Asset) error
	FindByID(ctx context.Context, id string) (favorites.Asset, error)
	FindAll(ctx context.Context, limit, offset int) (iter.Seq2[favorites.Asset, error], error)
	FindAllByUser(ctx context.Context, userID string, limit, offset int) (iter.Seq2[favorites.Asset, error], error)
	Delete(ctx context.Context, id, userID string) error
	UpdateDescription(ctx context.Context, id, description, userID string) (favorites.Asset, error)
}
