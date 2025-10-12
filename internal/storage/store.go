package storage

import (
	"sync"
	"time"

	"platform-go-challenge/internal/api_errors"
	"platform-go-challenge/internal/model"
)

// Store defines the interface for favorite storage
type Store interface {
	GetUserFavorites(userID string) ([]*model.Favorite, error)
	AddFavorite(userID string, asset *model.Asset) error
	RemoveFavorite(userID string, assetID string) error
	UpdateFavoriteDescription(userID string, assetID string, description string) error
}

// In-memory storage & RWMutex for thread safety
type MemoryStore struct {
	mu    sync.RWMutex
	users map[string]*model.User // key is user ID
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users: make(map[string]*model.User),
	}
}

// Returns all favorites for a user
func (ms *MemoryStore) GetUserFavorites(userID string) ([]*model.Favorite, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	user, exists := ms.users[userID]
	if !exists {
		return []*model.Favorite{}, nil
	}

	favorites := make([]*model.Favorite, 0, len(user.Favorites))
	for _, fav := range user.Favorites {
		favorites = append(favorites, fav)
	}

	return favorites, nil
}

// Adds a new asset to user's favorites
func (ms *MemoryStore) AddFavorite(userID string, asset *model.Asset) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	user, exists := ms.users[userID]
	if !exists {
		user = &model.User{
			ID:        userID,
			Favorites: make(map[string]*model.Favorite),
		}
		ms.users[userID] = user
	}

	baseAsset := asset.GetBaseAsset()
	now := time.Now().Unix()
	baseAsset.CreatedAt = now
	baseAsset.UpdatedAt = now

	user.Favorites[baseAsset.ID] = &model.Favorite{
		Asset:   *asset,
		AddedAt: now,
	}

	return nil
}

// Removes an asset from user's favorites
func (ms *MemoryStore) RemoveFavorite(userID string, assetID string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	user, exists := ms.users[userID]
	if !exists {
		return api_errors.ErrUserNotFound
	}

	if _, exists := user.Favorites[assetID]; !exists {
		return api_errors.ErrFavoriteNotFound
	}

	delete(user.Favorites, assetID)
	return nil
}

// Updates the description of a favorited asset
func (ms *MemoryStore) UpdateFavoriteDescription(userID string, assetID string, description string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	user, exists := ms.users[userID]
	if !exists {
		return api_errors.ErrUserNotFound
	}

	fav, exists := user.Favorites[assetID]
	if !exists {
		return api_errors.ErrFavoriteNotFound
	}

	baseAsset := fav.Asset.GetBaseAsset()
	if baseAsset != nil {
		baseAsset.Description = description
		baseAsset.UpdatedAt = time.Now().Unix()
	}

	return nil
}