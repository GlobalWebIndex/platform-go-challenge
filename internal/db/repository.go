// Package store defines an abstract Store interface and implementations used
// by the HTTP handlers. Implementations must be safe for concurrent use.
package db

import (
	"github.com/google/uuid"
	"github.com/mr-polymath/platform-go-challenge/internal/model"
)

type PagedQuery struct {
	Offset int
	Limit  int
	SortBy string
	Order  string // "asc" or "desc"
}

// FavouriteRepository defines storage operations needed by the application. Keep the
// interface small so adapters (in-memory, SQL, Redis) are easy to implement.
type FavouriteRepository interface {
	// GetAllFavourites returns all assets favourited by a user.
	GetAllFavourites(userID uuid.UUID, q PagedQuery) ([]model.Asset, error)
	// Add inserts an asset into the user's favourites.
	AddFavourite(userID uuid.UUID, a model.Asset) error
	// Remove deletes a specific asset from a user's favourites.
	RemoveFavourite(userID uuid.UUID, assetID uuid.UUID) error
	// UpdateDescription updates only the description field of an asset.
	PatchFavourite(userID uuid.UUID, assetID uuid.UUID, description string) error
}
