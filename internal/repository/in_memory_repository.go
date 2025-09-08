package repository

import (
	"errors"
	_ "go/ast"
	_ "net/netip"
	"platform-go-challenge/internal/domain/entities"
	"platform-go-challenge/internal/utils"
	"sync"
)

var FavouriteNotFoundError = errors.New("favourite not found")
var AssetNotFoundError = errors.New("asset not found")
var UserNotFoundError = errors.New("user not found")

// InMemoryRepository is a simple in-memory implementation of the Repository interface.
// It uses a mutex to ensure concurrent access is safe.
type InMemoryRepository struct {
	assets         map[string]entities.Asset
	users          map[string]*entities.User
	userFavourites map[string]map[string]entities.Asset
	mu             sync.RWMutex
}

// NewInMemoryRepository creates a new instance of InMemoryRepository.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		assets:         make(map[string]entities.Asset),
		users:          make(map[string]*entities.User),
		userFavourites: make(map[string]map[string]entities.Asset),
	}
}

func (r *InMemoryRepository) SaveAsset(asset entities.Asset) error {
	r.assets[asset.GetID()] = asset
	return nil
}

func (r *InMemoryRepository) GetAsset(assetID string) *entities.Asset {
	r.mu.RLock()
	defer r.mu.RUnlock()
	asset, ok := r.assets[assetID]
	if !ok {
		return nil
	}
	return &asset
}

// GetFavourites retrieves all favourite assets for a given user.
func (r *InMemoryRepository) GetFavourites(userID string, page, size int) ([]entities.Asset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	favourites, ok := r.userFavourites[userID]
	if !ok {
		// Return an empty map for a user with no favourites, rather than an error.
		return make([]entities.Asset, 0), nil
	}

	// Create a slice to hold the asset values.
	assets := make([]entities.Asset, 0, len(favourites))
	for _, asset := range favourites {
		assets = append(assets, asset)
	}

	return utils.PaginatedResult(assets, page, size), nil
}

// AddFavourite adds an asset to a user's favourites.
func (r *InMemoryRepository) AddFavourite(userID string, assetID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[userID]; !ok {
		return UserNotFoundError
	}

	// Ensure the user's favourites map exists
	if _, ok := r.userFavourites[userID]; !ok {
		r.userFavourites[userID] = make(map[string]entities.Asset)
	}

	asset, ok := r.assets[assetID]
	if !ok {
		return AssetNotFoundError // Or a more specific error for asset not found
	}

	// Check if the asset is already a favourite
	if _, ok := r.userFavourites[userID][assetID]; ok {
		// If it exists, return an error or handle accordingly. For now, it's a no-op.
		return errors.New("asset already in favourites")
	}

	r.userFavourites[userID][assetID] = asset
	return nil
}

// UpdateFavouriteDescription updates the description of a specific favourite asset for a user.
func (r *InMemoryRepository) UpdateFavouriteDescription(userID string, assetID string, newDescription string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[userID]; !ok {
		return UserNotFoundError
	}

	userFavs, ok := r.userFavourites[userID]
	if !ok {
		return FavouriteNotFoundError
	}

	asset, ok := userFavs[assetID]
	if !ok {
		return AssetNotFoundError
	}

	asset.SetDescription(newDescription)
	userFavs[assetID] = asset

	return nil
}

// RemoveFavourite removes a specific favourite asset for a user.
func (r *InMemoryRepository) RemoveFavourite(userID string, assetID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	userFavourites, ok := r.userFavourites[userID]
	if !ok {
		return FavouriteNotFoundError
	}

	if _, ok := userFavourites[assetID]; !ok {
		return AssetNotFoundError
	}

	delete(userFavourites, assetID)
	return nil
}

func (r *InMemoryRepository) GetUser(userID string) *entities.User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[userID]
	if !ok {
		return nil
	}

	return user
}

func (r *InMemoryRepository) SaveUser(user *entities.User) {
	r.users[user.ID] = user
}
