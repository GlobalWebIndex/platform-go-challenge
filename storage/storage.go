package storage

import (
	"challenge/models"
	"errors"
)

// Db is the storage access controller
var Db StorageHandler

// StorageHandler interface is impemented using all the required methods to
// interface with the various storage options
type StorageHandler interface {
	// initPrivilegedUser creates the original privileged account
	initPrivilegedUser(userID string, userPass string) error
	// AddAssetToUser adds the provided asset to the selected user if it doesn't exist
	AddAssetToUser(userID string, asset models.Asset) (bool, error)
	// RetrieveUserAssets returns the favourited assets for the selected user
	RetrieveUserAssets(userID string) (bool, []models.Asset, error)
	// UpdateUserAssetDescription updates the description of a favourited asset for
	// the selected user
	UpdateUserAssetDescription(userID string, assetID string, description string) (bool, error)
	// RemoveUser removes a user from the list of users
	RemoveUser(userID string) (bool, error)
	// RemoveUserAsset removes an asset from the selected user's favourites
	RemoveUserAsset(userID string, assetID string) (bool, error)
	// CheckUserExistence checks if a user exists
	CheckUserExistence(userID string) (bool, error)
	// CheckUserAssetExistence checks if an asset is favourited by the selected user
	CheckUserAssetExistence(userID string, assetID string) (bool, error)
	// CreateNewUser adds a new user to the list of users
	CreateNewUser(user models.User) (bool, error)
	// CheckAdmin checks if a user is an administrator
	CheckAdmin(id string) (bool, bool, error)
	// UpdateUserToken updates the user record with the latest cookie
	UpdateUserToken(user string, cookie string) (bool, error)
	// FetchUserToken retrieves the current user token
	FetchUserToken(user string) (bool, string, error)
	// FetchUserPass retrieves the encoded password for the selected user (if it exists)
	FetchUserPass(userID string) (bool, string, error)
	// FetchAllUsers retrieves all the currently registered users
	FetchAllUsers() ([]models.User, error)
	// Release closes the connection to the selected storage and releases resources
	Release() error
}

// ConnectoToStorage handles the initialization of the storage option
func ConnectToStorage(option string, adminCreds ...string) error {
	var err error
	storageOption := option
	switch storageOption {
	case "memory":
		Db = createMemoryHander()
		if len(adminCreds) == 2 {
			Db.initPrivilegedUser(adminCreds[0], adminCreds[1])
		}
	default:
		err = errors.New("unsupported storage option")
	}
	return err
}
