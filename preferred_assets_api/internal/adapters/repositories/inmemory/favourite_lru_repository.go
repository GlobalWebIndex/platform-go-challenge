package inmemory

import (
	"sync"
	"time"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/entities"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/ports"
	lru "github.com/hashicorp/golang-lru/v2"
)

var _ ports.FavouriteRepository = (*LRUFavouriteRepositoryImpl)(nil)

type LRUFavouriteRepositoryImpl struct {
	userAssetsCache *lru.Cache[string, map[string]time.Time]
	existsCache     *lru.Cache[string, bool]

	mu sync.RWMutex
}

func NewFavouriteRepository(cache *lru.Cache[string, map[string]time.Time], excache *lru.Cache[string, bool]) *LRUFavouriteRepositoryImpl {
	return &LRUFavouriteRepositoryImpl{
		userAssetsCache: cache,
		existsCache:     excache,
	}
}

func (c *LRUFavouriteRepositoryImpl) Add(f entities.FavouriteEntity) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Update exists cache
	existsKey := c.generateExistsKey(f.UserId, f.AssetId)
	c.existsCache.Add(existsKey, true)

	// Update user assets cache
	if userAssets, ok := c.userAssetsCache.Get(f.UserId); ok {
		userAssets[f.AssetId] = f.CreatedAt
		c.userAssetsCache.Add(f.UserId, userAssets)
	} else {
		// Create new user entry
		userAssets := make(map[string]time.Time)
		userAssets[f.AssetId] = f.CreatedAt
		c.userAssetsCache.Add(f.UserId, userAssets)
	}

	return nil
}

func (c *LRUFavouriteRepositoryImpl) Delete(userID, assetID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Update exists cache
	existsKey := c.generateExistsKey(userID, assetID)
	c.existsCache.Add(existsKey, false)

	// Update user assets cache
	if userAssets, ok := c.userAssetsCache.Get(userID); ok {
		delete(userAssets, assetID)
		// If user has no more favourites, remove the entry entirely
		if len(userAssets) == 0 {
			c.userAssetsCache.Remove(userID)
		} else {
			c.userAssetsCache.Add(userID, userAssets)
		}
	}

	return nil
}
func (c *LRUFavouriteRepositoryImpl) GetByUserID(userID string) ([]entities.FavouriteEntity, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if userAssets, ok := c.userAssetsCache.Get(userID); ok {
		favourites := make([]entities.FavouriteEntity, 0, len(userAssets))
		for assetID, createdAt := range userAssets {
			favourites = append(favourites, entities.FavouriteEntity{
				UserId:    userID,
				AssetId:   assetID,
				CreatedAt: createdAt,
			})
		}
		return favourites, nil
	}

	// Return empty slice if user has no favourites
	return []entities.FavouriteEntity{}, nil
}

func (c *LRUFavouriteRepositoryImpl) Exists(userID, assetID string) (bool, error) {
	existsKey := c.generateExistsKey(userID, assetID)

	// Try exists cache first
	c.mu.RLock()
	if exists, ok := c.existsCache.Get(existsKey); ok {
		c.mu.RUnlock()
		return exists, nil
	}
	c.mu.RUnlock()

	// Fallback to checking user assets cache
	c.mu.RLock()
	defer c.mu.RUnlock()

	if userAssets, ok := c.userAssetsCache.Get(userID); ok {
		_, exists := userAssets[assetID]
		return exists, nil
	}

	return false, nil
}

func (c *LRUFavouriteRepositoryImpl) generateExistsKey(userID, assetID string) string {
	return userID + ":" + assetID
}

// Get individual favourite with creation time
func (c *LRUFavouriteRepositoryImpl) GetFavourite(userID, assetID string) (*entities.FavouriteEntity, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if userAssets, ok := c.userAssetsCache.Get(userID); ok {
		if createdAt, exists := userAssets[assetID]; exists {
			return &entities.FavouriteEntity{
				UserId:    userID,
				AssetId:   assetID,
				CreatedAt: createdAt,
			}, nil
		}
	}

	return nil, nil
}

// Utility methods
func (c *LRUFavouriteRepositoryImpl) InvalidateUserCache(userID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.userAssetsCache.Remove(userID)
}
