package ports

import (
	"context"

	"github.com/SokratisChaimanas/platform-go-challenge/internal/domain"
	"github.com/google/uuid"
)

// FavouriteRepository stores and queries favourites.
type FavouriteRepository interface {
	// Create inserts a favourite. Duplicate should return domain.ErrFavouriteAlreadyExists.
	Create(ctx context.Context, f *domain.Favourite) error

	// Delete removes a favourite. Missing should return domain.ErrFavouriteNotFound.
	Delete(ctx context.Context, userID, assetID uuid.UUID) error

	// Exists checks whether (user, asset) is already favourited.
	Exists(ctx context.Context, userID, assetID uuid.UUID) (bool, error)

	// ListAssetsFavouritedByUserKeyset returns assets favourited by a user,
	// ordered by favourite.created_at,id, and an opaque next cursor.
	ListAssetsFavouritedByUserKeyset(ctx context.Context, userID uuid.UUID, limit int, after string) ([]domain.Asset, *string, error)
}
