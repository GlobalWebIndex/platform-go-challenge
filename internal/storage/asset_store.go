package storage

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/google/uuid"

	"platform-go-challenge/internal/common"
	"platform-go-challenge/internal/models"
)

// FavoritesStore to manage user favorite assets
// I have chosen to use a map of maps to represent user favorites for two main reasons:
//  1. Avoid circular dependencies: By defining FavoritesStore in the storage package, we avoid circular import issues between models and storage packages. This keeps the codebase modular and easier to maintain.
//  2. In a relational database, favorites would typically be represented as a join table between users and assets. Using a map of maps simulates this relationship effectively,
//     allowing quick lookups and management of user favorites without complicating the User model itself.
type FavoritesStore struct {
	favorites map[uuid.UUID]map[string]models.Asset // Map of UserID to map of AssetID to Asset
	userLocks map[uuid.UUID]*sync.RWMutex           // Per-user locks for concurrency safety
	mutex     sync.RWMutex                          // Global lock for the FavoritesStore
}

func NewFavoritesStore() *FavoritesStore {
	return &FavoritesStore{
		favorites: make(map[uuid.UUID]map[string]models.Asset),
		userLocks: make(map[uuid.UUID]*sync.RWMutex),
	}
}

// ---------- Favortites management methods ----------

func (fs *FavoritesStore) AddFavorite(userID uuid.UUID, favorite models.Asset) error {
	fs.mutex.Lock()

	if _, exists := fs.favorites[userID]; !exists {
		fs.userLocks[userID] = &sync.RWMutex{}
		fs.favorites[userID] = make(map[string]models.Asset)
	}
	userLock := fs.userLocks[userID]
	fs.mutex.Unlock()

	// Lock per user
	userLock.Lock()
	defer userLock.Unlock()

	if _, exists := fs.favorites[userID][favorite.GetID()]; exists {
		return errors.New("asset already in favorites")
	}

	fs.favorites[userID][favorite.GetID()] = favorite

	return nil
}

func (fs *FavoritesStore) RemoveFavorite(userID uuid.UUID, assetID string) error {
	fs.mutex.RLock()
	userLock, exists := fs.userLocks[userID]
	fs.mutex.RUnlock()

	// If user has no favorites, asset cannot exist
	if !exists {
		return errors.New(common.ErrAssetNotFound)
	}

	// Lock per user
	userLock.Lock()
	defer userLock.Unlock()

	userFavorites, exists := fs.favorites[userID]
	if !exists || userFavorites == nil {
		return errors.New(common.ErrAssetNotFound)
	}

	if _, exists := fs.favorites[userID][assetID]; !exists {
		return errors.New(common.ErrAssetNotFound)
	}
	delete(fs.favorites[userID], assetID)

	return nil
}

func (fs *FavoritesStore) RemoveUserFavorites(userID uuid.UUID) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	// Delete favorites entry
	delete(fs.favorites, userID)

	// Also delete per-user lock
	delete(fs.userLocks, userID)
}

func (fs *FavoritesStore) GetUserFavorites(userID uuid.UUID) ([]models.Asset, error) {
	fs.mutex.RLock()
	userLock, exists := fs.userLocks[userID]
	fs.mutex.RUnlock()

	// If user has no favorites yet, return empty list (not an error)
	if !exists {
		return []models.Asset{}, nil
	}

	// Lock per user
	userLock.RLock()
	defer userLock.RUnlock()

	assetsMap, exists := fs.favorites[userID]
	if !exists {
		return []models.Asset{}, nil
	}

	// returning a copy of assets as slice to avoid external modification
	favorites := make([]models.Asset, 0, len(assetsMap))
	for _, asset := range assetsMap {
		favorites = append(favorites, asset)
	}

	return favorites, nil
}

func (fs *FavoritesStore) GetFavoriteByAssetID(userID uuid.UUID, assetID string) (models.Asset, error) {
	fs.mutex.RLock()
	userLock, exists := fs.userLocks[userID]
	fs.mutex.RUnlock()

	if !exists {
		return nil, errors.New(common.ErrUserNotFound)
	}

	// Lock per user
	userLock.RLock()
	defer userLock.RUnlock()

	asset, exists := fs.favorites[userID][assetID]
	if !exists {
		return nil, errors.New(common.ErrAssetNotFound)
	}

	return asset, nil
}

func (fs *FavoritesStore) UpdateFavorite(userID uuid.UUID, updatedAsset models.Asset) error {
	fs.mutex.RLock()
	userLock, exists := fs.userLocks[userID]
	defer fs.mutex.RUnlock()

	if !exists {
		return errors.New(common.ErrUserNotFound)
	}

	// Lock per user
	userLock.Lock()
	defer userLock.Unlock()

	userFavorites, exists := fs.favorites[userID]
	if !exists {
		return errors.New("user has no favorites")
	}

	currentAsset, exists := userFavorites[updatedAsset.GetID()]
	if !exists {
		return fmt.Errorf("asset %s not found in user %s favorites", updatedAsset.GetID(), userID)
	}

	currentVal := reflect.ValueOf(currentAsset).Elem()
	updatedVal := reflect.ValueOf(updatedAsset).Elem()
	typ := currentVal.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Skip fields tagged `merge:"ignore"`
		if field.Tag.Get("merge") == "ignore" {
			continue
		}

		currentField := currentVal.Field(i)
		updatedField := updatedVal.Field(i)

		// Only overwrite if updatedField is not zero
		if !updatedField.IsZero() {
			currentField.Set(updatedField)
		}
	}

	// Update the UpdatedAt timestamp explicitly
	if b, ok := currentAsset.(interface{ SetUpdatedTimestamp(time.Time) }); ok {
		b.SetUpdatedTimestamp(time.Now())
	}

	return nil
}
