package favourites

import "context"

// Repository provides storage operations for assets and favourites. Operations
// are user-centric, keyed by userID for efficient per-user queries.
type Repository interface {
	// CreateFavourite creates an asset and marks it as a favourite
	// in a single atomic operation. If asset.ID is empty, a new ID is assigned.
	// If asset.ID is non-empty, ErrInvalidAssetPayload should be returned.
	CreateFavourite(ctx context.Context, userID string, asset Asset) (FavouriteEntry, error)

	// ListFavourites returns favourites for the given user, optionally filtered
	// by type. Should be proportional to that user's favourites only.
	ListFavourites(ctx context.Context, userID string, filter ListFilter) ([]FavouriteEntry, error)

	// RemoveFavourite removes the favourite link.
	// Idempotent: does not return an error if already absent.
	RemoveFavourite(ctx context.Context, userID, assetID string) error

	// UpdateDescriptionForUser updates an asset's description only if the given
	// user has this asset as a favourite. 
	// Returns ErrFavouriteNotFound if the favourite (user, asset) does not exist.
	UpdateDescriptionForUser(ctx context.Context, userID, assetID, description string) (Asset, error)
}
