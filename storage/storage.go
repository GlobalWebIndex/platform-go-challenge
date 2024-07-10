package storage

import (
	"strconv"
	"sync"

	"github.com/artemis13/platform-go-challenge/models"
)

var (
	users = make(map[uint]*models.User)
	mu    sync.RWMutex
)

// ConvertStringToUint safely converts a string to uint
func ConvertStringToUint(s string) (uint, error) {
	id, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// GetUser safely retrieves a user from the map
func GetUser(userID uint) (*models.User, bool) {
	mu.RLock()
	defer mu.RUnlock()
	user, exists := users[userID]
	return user, exists
}

// UpdateUserFavorites safely updates the favorites of a user
func UpdateUserFavorites(userID uint, asset models.Asset) {
	mu.Lock()
	defer mu.Unlock()
	user, exists := users[userID]
	if !exists {
		user = &models.User{ID: userID, Favorites: []models.Asset{}}
		users[userID] = user
	}
	user.Favorites = append(user.Favorites, asset)
}

// For testing purposes, clear the users map
func ClearUsers() {
	mu.Lock()
	defer mu.Unlock()
	users = make(map[uint]*models.User)
}

// For testing purposes, AddUser safely adds a user to the map
func AddUser(user *models.User) {
	mu.Lock()
	defer mu.Unlock()
	users[user.ID] = user
}

// For testing purposes, remove a user favorite
func RemoveUserFavorite(userID uint, assetID string) {
	mu.Lock()
	defer mu.Unlock()
	user, exists := users[userID]
	if exists {
		for i, asset := range user.Favorites {
			if asset.ID == assetID {
				user.Favorites = append(user.Favorites[:i], user.Favorites[i+1:]...)
				break
			}
		}
	}
}
