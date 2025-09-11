package repository

import (
	"errors"
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
	userFavourites map[string]map[string]string
	globalLock     sync.RWMutex
	perUserLock    map[string]*sync.RWMutex
}

// NewInMemoryRepository creates a new instance of InMemoryRepository.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		assets:         make(map[string]entities.Asset),
		users:          make(map[string]*entities.User),
		userFavourites: make(map[string]map[string]string),
		globalLock:     sync.RWMutex{},
		perUserLock:    make(map[string]*sync.RWMutex),
	}
}

func (r *InMemoryRepository) SaveAsset(asset entities.Asset) error {
	r.globalLock.Lock()
	defer r.globalLock.Unlock()
	r.assets[asset.GetID()] = asset
	return nil
}

func (r *InMemoryRepository) GetAsset(assetID string) *entities.Asset {
	r.globalLock.RLock()
	defer r.globalLock.RUnlock()
	asset, ok := r.assets[assetID]
	if !ok {
		return nil
	}
	return &asset
}

// GetFavourites retrieves all favourite assets for a given user.
func (r *InMemoryRepository) GetFavourites(userID string, page, size int) ([]entities.Asset, error) {
	// First, acquire a read lock on global data to check for user existence.
	r.globalLock.RLock()
	if _, ok := r.users[userID]; !ok {
		r.globalLock.RUnlock()
		return nil, UserNotFoundError
	}
	r.globalLock.RUnlock()

	userLock := r.getUserLock(userID)
	userLock.RLock()
	defer userLock.RUnlock()

	favourites, ok := r.userFavourites[userID]
	if !ok {
		return make([]entities.Asset, 0), nil
	}

	assets := make([]entities.Asset, 0, len(favourites))
	r.globalLock.RLock()
	for _, assetID := range favourites {
		if asset, ok := r.assets[assetID]; ok {
			assets = append(assets, asset)
		}
	}
	r.globalLock.RUnlock()

	return utils.PaginatedResult(assets, page, size), nil
}

// AddFavourite adds an asset to a user's favourites.
func (r *InMemoryRepository) AddFavourite(userID string, assetID string) error {
	// Acquire a read lock to check for user and asset existence safely.
	r.globalLock.RLock()
	if _, ok := r.users[userID]; !ok {
		r.globalLock.RUnlock()
		return UserNotFoundError
	}
	_, ok := r.assets[assetID]
	if !ok {
		r.globalLock.RUnlock()
		return AssetNotFoundError
	}
	r.globalLock.RUnlock()

	// Now that we know the user and asset exist, acquire the per-user lock.
	userLock := r.getUserLock(userID)
	userLock.Lock()
	defer userLock.Unlock()

	// Ensure the user's favourites map exists
	if _, ok := r.userFavourites[userID]; !ok {
		r.userFavourites[userID] = make(map[string]string)
	}

	// Check if the asset is already a favourite
	if _, ok := r.userFavourites[userID][assetID]; ok {
		return errors.New("asset already in favourites")
	}

	r.userFavourites[userID][assetID] = assetID
	return nil
}

// UpdateFavouriteDescription updates the description of a specific favourite asset for a user.
func (r *InMemoryRepository) UpdateFavouriteDescription(userID string, assetID string, newDescription string) error {
	// The operation affects a global resource (the asset's description).
	// Therefore, we need to acquire a global write lock.
	r.globalLock.Lock()
	defer r.globalLock.Unlock()

	if _, ok := r.users[userID]; !ok {
		return UserNotFoundError
	}

	asset, ok := r.assets[assetID]
	if !ok {
		return AssetNotFoundError
	}

	asset.SetDescription(newDescription)
	r.assets[assetID] = asset

	return nil
}

// RemoveFavourite removes a specific favourite asset for a user.
func (r *InMemoryRepository) RemoveFavourite(userID string, assetID string) error {
	userLock := r.getUserLock(userID)
	userLock.Lock()
	defer userLock.Unlock()

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
	r.globalLock.RLock()
	defer r.globalLock.RUnlock()

	user, ok := r.users[userID]
	if !ok {
		return nil
	}
	return user
}

func (r *InMemoryRepository) SaveUser(user *entities.User) {
	r.globalLock.Lock()
	defer r.globalLock.Unlock()
	r.users[user.ID] = user
}

// getUserLock returns the per-user lock. It creates one if it doesn't exist.
func (r *InMemoryRepository) getUserLock(userID string) *sync.RWMutex {
	// Acquire a write lock to safely check and create the per-user lock.
	r.globalLock.Lock()
	defer r.globalLock.Unlock()

	if _, ok := r.perUserLock[userID]; !ok {
		r.perUserLock[userID] = &sync.RWMutex{}
	}
	return r.perUserLock[userID]
}
