package data

import "platform-go-challenge/models"

// UserDataService defines methods for accessing and modifying user data.
type UserDataService interface {
	// GetUser retrieves a user by ID.
	GetUser(userID string) (*models.User, error)

	// AddUser adds a new user.
	AddUser(user *models.User) error

	// AddAssetToFavourites adds an asset to the user's list of favourites.
	AddAssetToFavourites(userID string, asset models.Asset) error

	// RemoveAssetFromFavourites removes an asset from the user's list of favourites.
	RemoveAssetFromFavourites(userID string, assetID string) error

	// EditAssetDescription edits the description of an asset in the user's list of favourites.
	EditAssetDescription(userID string, assetID string, description string) error
}
