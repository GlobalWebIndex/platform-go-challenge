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
	globalLock     sync.RWMutex
	perUserLock    map[string]*sync.RWMutex
}

// NewInMemoryRepository creates a new instance of InMemoryRepository.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		assets:         make(map[string]entities.Asset),
		users:          make(map[string]*entities.User),
		userFavourites: make(map[string]map[string]entities.Asset),
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

	_, ok := r.users[userID]
	if !ok {
		return nil, UserNotFoundError
	}

	userLock := r.getUserLock(userID)
	userLock.RLock()
	defer userLock.RUnlock()

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
	r.globalLock.RLock()
	if _, ok := r.users[userID]; !ok {
		return UserNotFoundError
	}
	defer r.globalLock.RUnlock()

	asset, ok := r.assets[assetID]
	if !ok {
		return AssetNotFoundError
	}

	userLock := r.getUserLock(userID)
	userLock.Lock()
	defer userLock.Unlock()
	// Ensure the user's favourites map exists
	if _, ok := r.userFavourites[userID]; !ok {
		r.userFavourites[userID] = make(map[string]entities.Asset)
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
	r.globalLock.Lock()
	defer r.globalLock.Unlock()

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

func (r *InMemoryRepository) getUserLock(userID string) *sync.RWMutex {
	r.globalLock.Lock()
	defer r.globalLock.Unlock()

	if _, ok := r.perUserLock[userID]; !ok {
		r.perUserLock[userID] = &sync.RWMutex{}
	}
	return r.perUserLock[userID]
}
