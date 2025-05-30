package domain

import (
	"context"
	"errors"
)

//go:generate mockgen -destination=../persistence/mocks/mock_favorite.go -package=mocks -typed . FavoriteRepository
//go:generate mockgen -destination=../app/mocks/mock_favorite.go -package=mocks -typed . FavoriteService

var (
	// ErrInvalidUserID indicates that the provided user ID is invalid.
	ErrInvalidUserID = errors.New("invalid user ID")
	// ErrInvalidAssetID indicates that the provided asset ID is invalid.
	ErrInvalidAssetID = errors.New("invalid asset ID")
	// ErrInvalidFavoriteID indicates that the provided favorite ID is invalid.
	ErrInvalidFavoriteID = errors.New("invalid favorite ID")
	// ErrFavoriteNotFound indicates that the requested favorite asset was not found.
	ErrFavoriteNotFound = errors.New("favorite not found")
	// ErrInternalServerError indicates that an internal server error occurred.
	ErrInternalServerError = errors.New("internal server error")
	// ErrDuplicateFavorite indicates that the favorite asset already exists.
	ErrDuplicateFavorite = errors.New("favorite already exists")
)

// Favorite represents a user's favorite asset in the system.
// It contains information about the favorite's ID, associated user, asset, and description.
type Favorite struct {
	ID          int    // Unique identifier for the favorite.
	UserID      int    // ID of the user who marked the asset as favorite.
	AssetID     int    // ID of the asset marked as favorite.
	Description string // Optional description of the favorite.
}

// FavoriteService is an interface for managing user favorite assets.
// It defines methods for retrieving, adding, editing, and deleting favorites.
type FavoriteService interface {
	GetUserFavorites(ctx context.Context, userID int) ([]Favorite, error)                   // Retrieves all favorites for a given user.
	Add(ctx context.Context, favorite Favorite) (Favorite, error)                           // Adds a new favorite asset for a user.
	Edit(ctx context.Context, favoriteID, userID int, description string) (Favorite, error) // Edits the description of an existing favorite.
	Delete(ctx context.Context, favoriteID, userID int) error                               // Deletes a favorite asset for a user.
}

// FavoriteRepository is an interface for managing user favorite assets in the database.
// It defines methods for searching, inserting, updating, and removing favorites.
type FavoriteRepository interface {
	SearchByUserID(ctx context.Context, userID int) ([]Favorite, error)                                  // Searches for all favorites by user ID.
	Insert(ctx context.Context, favorite Favorite) (Favorite, error)                                     // Inserts a new favorite into the database.
	UpdateDescription(ctx context.Context, favoriteID, userID int, description string) (Favorite, error) // Updates the description of a favorite in the database.
	Remove(ctx context.Context, favoriteID, userID int) error                                            // Removes a favorite from the database.
}
