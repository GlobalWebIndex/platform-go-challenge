package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"assetsApp/internal/models"

	"github.com/google/uuid"
)

type CachedStore struct {
	db    AssetStore
	cache *RedisClient
}

type cachedFavourite struct {
	UserID    uuid.UUID       `json:"user_id"`
	AssetType string          `json:"asset_type"`
	AssetData json.RawMessage `json:"asset_data"`
}

func NewCachedStore(db AssetStore, cache *RedisClient) *CachedStore {
	return &CachedStore{
		db:    db,
		cache: cache,
	}
}

// ----- Asset methods without caching -----
func favsCacheKey(userID uuid.UUID) string {
	return fmt.Sprintf("favourites:%s", userID.String())
}

func (c *CachedStore) Get(userID uuid.UUID) []models.Asset {
	return c.db.Get(userID)
}

func (c *CachedStore) Add(userID uuid.UUID, asset models.Asset) {
	c.db.Add(userID, asset)
}

func (c *CachedStore) Remove(userID uuid.UUID, assetID string) bool {
	res := c.db.Remove(userID, assetID)
	if res && c.cache != nil {
		_ = c.cache.Del(context.Background(), favsCacheKey(userID))
	}
	return res
}

func (c *CachedStore) EditDescription(userID uuid.UUID, assetID, newDesc string) bool {
	res := c.db.EditDescription(userID, assetID, newDesc)
	if res && c.cache != nil {
		_ = c.cache.Del(context.Background(), favsCacheKey(userID))
	}
	return res
}

// ----- Favourites with caching -----

func (c *CachedStore) GetFavourites(userID uuid.UUID) []models.Favourite {
	ctx := context.Background()

	if c.cache != nil {
		if cached, err := c.cache.Get(ctx, favsCacheKey(userID)); err == nil && cached != "" {
			fmt.Printf("cached_store: cache hit for favourites of user %v\n", userID)

			var cachedFavs []cachedFavourite
			if err := json.Unmarshal([]byte(cached), &cachedFavs); err == nil {
				var favs []models.Favourite
				for _, cf := range cachedFavs {
					var asset models.Asset
					switch cf.AssetType {
					case "chart":
						var a models.Chart
						if err := json.Unmarshal(cf.AssetData, &a); err == nil {
							asset = &a
						}
					case "insight":
						var a models.Insight
						if err := json.Unmarshal(cf.AssetData, &a); err == nil {
							asset = &a
						}
					case "audience":
						var a models.Audience
						if err := json.Unmarshal(cf.AssetData, &a); err == nil {
							asset = &a
						}
					}
					if asset != nil {
						favs = append(favs, models.Favourite{
							UserID: cf.UserID,
							Asset:  asset,
						})
					}
				}
				return favs
			}
			log.Printf("cached_store: failed to unmarshal favourites cache for user %v: %v", userID, err)
		}
	}

	// Fallback to DB
	favs := c.db.GetFavourites(userID)

	// Write back to cache
	if c.cache != nil {
		var cachedFavs []cachedFavourite
		for _, f := range favs {
			assetJSON, _ := json.Marshal(f.Asset)
			var assetType string
			switch f.Asset.(type) {
			case *models.Chart:
				assetType = "chart"
			case *models.Insight:
				assetType = "insight"
			case *models.Audience:
				assetType = "audience"
			}
			cachedFavs = append(cachedFavs, cachedFavourite{
				UserID:    f.UserID,
				AssetType: assetType,
				AssetData: assetJSON,
			})
		}
		if b, err := json.Marshal(cachedFavs); err == nil {
			_ = c.cache.Set(ctx, favsCacheKey(userID), string(b))
		} else {
			log.Printf("cached_store: failed to marshal favourites for caching: %v", err)
		}
	}

	return favs
}

func (c *CachedStore) AddFavourite(userID uuid.UUID, assetID, assetType string) bool {
	res := c.db.AddFavourite(userID, assetID, assetType)
	if res && c.cache != nil {
		_ = c.cache.Del(context.Background(), favsCacheKey(userID))
	}
	return res
}

func (c *CachedStore) RemoveFavourite(userID uuid.UUID, assetID string) bool {
	res := c.db.RemoveFavourite(userID, assetID)
	if res && c.cache != nil {
		_ = c.cache.Del(context.Background(), favsCacheKey(userID))
	}
	return res
}
