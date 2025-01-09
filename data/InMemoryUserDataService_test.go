package data

import (
	"errors"
	"platform-go-challenge/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

// #region GetUser tests

// TestGetUserSuccess tests the successful retrieval of a user.
func TestGetUserSuccess(t *testing.T) {
	svc := NewInMemoryUserDataService()
	user := &models.User{ID: "123"}
	_ = svc.AddUser(user)

	retrievedUser, err := svc.GetUser("123")
	assert.NoError(t, err)
	assert.NotNil(t, retrievedUser)
	assert.Equal(t, "123", retrievedUser.ID)
}

// TestGetUserNotFound tests retrieval of a non-existent user.
func TestGetUserNotFound(t *testing.T) {
	svc := NewInMemoryUserDataService()

	retrievedUser, err := svc.GetUser("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, retrievedUser)
	assert.Equal(t, errors.New("user not found"), err)
}

// #endregion

// #region AddUser tests

// TestAddUserSuccess tests the successful addition of a user.
func TestAddUserSuccess(t *testing.T) {
	svc := NewInMemoryUserDataService()
	user := &models.User{ID: "123"}

	err := svc.AddUser(user)
	assert.NoError(t, err)
	assert.Equal(t, user, svc.users["123"])
}

// #endregion

// #region AddAssetToFavourites tests

// TestAddAssetToFavouritesSuccess tests adding an asset to the user's favourites successfully.
func TestAddAssetToFavouritesSuccess(t *testing.T) {
	svc := NewInMemoryUserDataService()
	user := &models.User{ID: "123"}
	_ = svc.AddUser(user)
	asset := models.Asset{ID: "456"}

	err := svc.AddAssetToFavourites("123", asset)
	assert.NoError(t, err)
	assert.Contains(t, svc.users["123"].Favourites, asset)
}

// TestAddAssetToFavouritesDuplicate tests adding a duplicate asset to the user's favourites.
func TestAddAssetToFavouritesDuplicate(t *testing.T) {
	svc := NewInMemoryUserDataService()
	user := &models.User{ID: "123"}
	_ = svc.AddUser(user)
	asset := models.Asset{ID: "456"}
	_ = svc.AddAssetToFavourites("123", asset)

	err := svc.AddAssetToFavourites("123", asset)
	assert.Error(t, err)
	assert.Equal(t, errors.New("duplicate asset ID"), err)
}

// #endregion

// #region RemoveAssetFromFavourites tests

// TestRemoveAssetFromFavouritesSuccess tests the successful removal of an asset from the user's favourites.
func TestRemoveAssetFromFavouritesSuccess(t *testing.T) {
	svc := NewInMemoryUserDataService()
	user := &models.User{ID: "123"}
	_ = svc.AddUser(user)
	asset := models.Asset{ID: "456"}
	_ = svc.AddAssetToFavourites("123", asset)

	err := svc.RemoveAssetFromFavourites("123", "456")
	assert.NoError(t, err)
	assert.NotContains(t, svc.users["123"].Favourites, asset)
}

// TestRemoveAssetFromFavouritesNotFound tests removal of a non-existent asset from the user's favourites.
func TestRemoveAssetFromFavouritesNotFound(t *testing.T) {
	svc := NewInMemoryUserDataService()
	user := &models.User{ID: "123"}
	_ = svc.AddUser(user)
	asset := models.Asset{ID: "456"}
	_ = svc.AddAssetToFavourites("123", asset)

	err := svc.RemoveAssetFromFavourites("123", "nonexistent")
	assert.Error(t, err)
	assert.Equal(t, errors.New("asset not found"), err)
}

// #endregion

// #region EditAssetDescription tests

// TestEditAssetDescriptionSuccess tests the successful editing of an asset's description.
func TestEditAssetDescriptionSuccess(t *testing.T) {
	svc := NewInMemoryUserDataService()
	user := &models.User{ID: "123"}
	_ = svc.AddUser(user)
	asset := models.Asset{ID: "456"}
	_ = svc.AddAssetToFavourites("123", asset)

	err := svc.EditAssetDescription("123", "456", "New Description")
	assert.NoError(t, err)
	assert.Equal(t, "New Description", svc.users["123"].Favourites[0].Description)
}

// TestEditAssetDescriptionNotFound tests editing the description of a non-existent asset.
func TestEditAssetDescriptionNotFound(t *testing.T) {
	svc := NewInMemoryUserDataService()
	user := &models.User{ID: "123"}
	_ = svc.AddUser(user)
	asset := models.Asset{ID: "456"}
	_ = svc.AddAssetToFavourites("123", asset)

	err := svc.EditAssetDescription("123", "nonexistent", "New Description")
	assert.Error(t, err)
	assert.Equal(t, errors.New("asset not found"), err)
}

// #endregion
