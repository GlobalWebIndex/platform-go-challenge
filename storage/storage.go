package storage

import (
	"platform-go-challenge/models"
	"sync"
)

// Mutex for safe concurrent access
var (
	mu        sync.Mutex
	favorites = make(map[string][]models.AssetInterface) // UserID → List of Assets
)

// GetFavoritesByType retrieves all favorites of a specific type for a given user
func GetFavoritesByType(userID string, assetType string) []models.AssetInterface {
	mu.Lock()
	defer mu.Unlock()
	var filtered []models.AssetInterface
	for _, fav := range favorites[userID] {
		if fav.GetType() == models.AssetType(assetType) {
			filtered = append(filtered, fav)
		}
	}
	return filtered
}

// Add asset to favorites
func AddFavorite(userID string, asset models.AssetInterface) {
	mu.Lock()
	defer mu.Unlock()
	for _, fav := range favorites[userID] {
		if fav.GetID() == asset.GetID() {
			return // Asset already exists, don't add again
		}
	}
	favorites[userID] = append(favorites[userID], asset)
}

// Remove asset from favorites
func RemoveFavorite(userID, assetID string) {
	mu.Lock()
	defer mu.Unlock()
	userAssets := favorites[userID]
	for i, asset := range userAssets {
		if asset.GetID() == assetID {
			favorites[userID] = append(userAssets[:i], userAssets[i+1:]...)
			return
		}
	}
	// If asset not found, no deletion performed
	return
}

// Edit asset description
func EditFavorite(userID, assetID, newDescription string) {
	mu.Lock()
	defer mu.Unlock()
	userAssets, exists := favorites[userID]
	if !exists {
		return
	}
	for i, asset := range userAssets {
		if asset.GetID() == assetID {
			switch v := asset.(type) {
			case models.Chart:
				v.Description = newDescription
				favorites[userID][i] = v
			case models.Insight:
				v.Description = newDescription
				favorites[userID][i] = v
			case models.Audience:
				v.Description = newDescription
				favorites[userID][i] = v
			}
			return
		}
	}
}

func UserExists(userID string) bool {
	mu.Lock()
	defer mu.Unlock()
	_, exists := favorites[userID]
	return exists
}

func AssetExists(userID, assetID string) bool {
	mu.Lock()
	defer mu.Unlock()
	for _, asset := range favorites[userID] {
		if asset.GetID() == assetID {
			return true
		}
	}
	return false
}

func InitializeTestUser(userID string) {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := favorites[userID]; !exists {
		favorites[userID] = []models.AssetInterface{}
	}
}
