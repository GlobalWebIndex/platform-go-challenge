package storage_test

import (
	"testing"

	"github.com/artemis13/platform-go-challenge/models"
	"github.com/artemis13/platform-go-challenge/storage"

	"github.com/stretchr/testify/assert"
)

func setupTestData() {
	storage.ClearUsers()
}

func TestConvertStringToUint(t *testing.T) {
	id, err := storage.ConvertStringToUint("123")
	assert.NoError(t, err)
	assert.Equal(t, uint(123), id)

	_, err = storage.ConvertStringToUint("invalid")
	assert.Error(t, err)
}

func TestGetUser(t *testing.T) {
	setupTestData()

	// Add a test user
	testUser := &models.User{ID: 1}
	storage.AddUser(testUser)

	// Test retrieval of the user
	user, exists := storage.GetUser(1)
	assert.True(t, exists)
	assert.Equal(t, testUser, user)

	// Test retrieval of a non-existent user
	_, exists = storage.GetUser(2)
	assert.False(t, exists)
}

func TestAddUser(t *testing.T) {
	setupTestData()

	// Add a test user
	testUser := &models.User{ID: 1}
	storage.AddUser(testUser)

	// Verify that the user was added
	user, exists := storage.GetUser(1)
	assert.True(t, exists)
	assert.Equal(t, testUser, user)
}

func TestUpdateUserFavorites(t *testing.T) {
	setupTestData()

	// Add a test user
	testUser := &models.User{ID: 1}
	storage.AddUser(testUser)

	// Add a favorite asset
	asset := models.Asset{ID: "1", Description: "Favorite 1"}
	storage.UpdateUserFavorites(1, asset)

	// Verify that the asset was added
	user, exists := storage.GetUser(1)
	assert.True(t, exists)
	assert.Equal(t, 1, len(user.Favorites))
	assert.Equal(t, "Favorite 1", user.Favorites[0].Description)
}

func TestRemoveUserFavorite(t *testing.T) {
	setupTestData()

	// Add a test user with a favorite asset
	testUser := &models.User{
		ID: 1,
		Favorites: []models.Asset{
			{ID: "1", Description: "Favorite 1"},
		},
	}
	storage.AddUser(testUser)

	// Remove the favorite asset
	storage.RemoveUserFavorite(1, "1")

	// Verify that the asset was removed
	user, exists := storage.GetUser(1)
	assert.True(t, exists)
	assert.Equal(t, 0, len(user.Favorites))
}
