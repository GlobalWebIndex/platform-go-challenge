package repository

import (
	"errors"
	"fmt"
	"sync"

	"github.com/ceciivanov/platform-go-challenge/internal/models"
	"github.com/ceciivanov/platform-go-challenge/internal/repository/mock_data"
)

// InMemoryUserRepository is an in-memory implementation of the UserRepository interface
// InMemoryUserRepository contains a map of all Users with their favorite assets
type InMemoryUserRepository struct {
	Users map[int]models.User
	mu    sync.RWMutex // mu is a read-write mutex to protect the Users map from concurrent access
}

// NewInMemoryUserRepository creates a new instance of InMemoryUserRepository
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		Users: make(map[int]models.User),
	}
}

// GenerateSampleUsers generates sample users with sample assets
func (repo *InMemoryUserRepository) GenerateSampleUsers(NumberOfUsers, NumberOfAssets int) {
	// Lock the Users map for writing
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.Users = mock_data.GenerateMockData(NumberOfUsers, NumberOfAssets)
}

// GetUserFavorites returns a map of user's favorite assets
func (repo *InMemoryUserRepository) GetUserFavorites(userID int) (map[int]models.Asset, error) {
	// Lock the Users map for reading
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	user, ok := repo.Users[userID]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user.Favourites, nil
}

// AddUserFavorite adds an asset to the user's favorites
func (repo *InMemoryUserRepository) AddUserFavorite(userID int, asset models.Asset) error {
	// Lock the Users map for writing
	repo.mu.Lock()
	defer repo.mu.Unlock()

	user, ok := repo.Users[userID]
	if !ok {
		return errors.New("user not found")
	}

	if _, ok := user.Favourites[asset.GetID()]; ok {
		return errors.New("asset already exists")
	}

	user.Favourites[asset.GetID()] = asset
	repo.Users[userID] = user
	return nil
}

// DeleteUserFavorite deletes an asset from the user's favorites
func (repo *InMemoryUserRepository) DeleteUserFavorite(userID, assetID int) error {
	// Lock the Users map for writing
	repo.mu.Lock()
	defer repo.mu.Unlock()

	user, ok := repo.Users[userID]
	if !ok {
		return errors.New("user not found")
	}

	if _, ok := user.Favourites[assetID]; !ok {
		return errors.New("asset not found")
	}

	delete(user.Favourites, assetID)
	repo.Users[userID] = user
	return nil
}

// EditUserFavorite edits an asset in the user's favorites
func (repo *InMemoryUserRepository) EditUserFavorite(userID int, assetID int, asset models.Asset) error {
	// Lock the Users map for writing
	repo.mu.Lock()
	defer repo.mu.Unlock()

	user, ok := repo.Users[userID]
	if !ok {
		return errors.New("user not found")
	}

	if _, ok := user.Favourites[assetID]; !ok {
		return errors.New("asset not found")
	}

	// Validate asset type and ID match the existing asset
	existingAsset := user.Favourites[assetID]

	if existingAsset.GetID() != asset.GetID() {
		return fmt.Errorf("edited asset ID does not match existing asset ID")
	}

	if existingAsset.GetType() != asset.GetType() {
		return fmt.Errorf("edited asset type does not match existing asset type")
	}

	user.Favourites[asset.GetID()] = asset
	repo.Users[userID] = user
	return nil
}
