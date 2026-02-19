package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/gwi/favourites-api/internal/cache"
	"github.com/gwi/favourites-api/internal/models"
	"github.com/gwi/favourites-api/internal/pagination"
	"github.com/gwi/favourites-api/internal/store"
)

type FavouritesService struct {
	storage  *store.Storage
	cache    cache.Cache
	cacheTTL time.Duration
}

func NewFavouritesService(storage *store.Storage, cache cache.Cache, ttl time.Duration) *FavouritesService {
	return &FavouritesService{
		storage:  storage,
		cache:    cache,
		cacheTTL: ttl,
	}
}

// AddFavourite adds an asset to a user's favourites
func (s *FavouritesService) AddFavourite(ctx context.Context, userID string, asset models.Asset) error {
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}
	if asset == nil {
		return fmt.Errorf("asset is required")
	}
	if asset.GetAssetID() == "" {
		return fmt.Errorf("asset_id is required")
	}

	// Ensure favourited_at is set (UTC)
	now := time.Now().UTC()
	if asset.GetFavouritedAt().IsZero() {
		switch a := asset.(type) {
		case *models.Chart:
			a.BaseAsset.FavouritedAt = now
		case *models.Insight:
			a.BaseAsset.FavouritedAt = now
		case *models.Audience:
			a.BaseAsset.FavouritedAt = now
		default:
			return fmt.Errorf("unsupported asset type")
		}
	}

	if err := s.storage.Favourites.AddFavourite(ctx, userID, asset); err != nil {
		return err
	}

	// Invalidate cache for this user
	_ = s.cache.Delete(ctx, "favourites:user:"+userID)

	return nil
}

// RemoveFavourite removes an asset from a user's favourites
func (s *FavouritesService) RemoveFavourite(ctx context.Context, userID, assetID string) error {
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}
	if assetID == "" {
		return fmt.Errorf("asset ID is required")
	}

	if err := s.storage.Favourites.RemoveFavourite(ctx, userID, assetID); err != nil {
		return err
	}

	// Invalidate cache for this user
	cacheKey := fmt.Sprintf("favourites:user:%s", userID)
	s.cache.Delete(ctx, cacheKey)

	return nil
}

// GetFavourites retrieves all favourites for a user (used internally)
func (s *FavouritesService) GetFavourites(ctx context.Context, userID string) ([]models.Asset, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Try cache first
	cacheKey := fmt.Sprintf("favourites:user:%s", userID)
	if cached, found := s.cache.Get(ctx, cacheKey); found {
		if favourites, ok := cached.([]models.Asset); ok {
			return favourites, nil
		}
	}

	// Cache miss - get from storage
	favourites, err := s.storage.Favourites.GetUserFavourites(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Cache the result (5 minutes TTL)
	s.cache.Set(ctx, cacheKey, favourites, 5*time.Minute)

	return favourites, nil
}

// GetFavouritesPaginated retrieves paginated favourites for a user
func (s *FavouritesService) GetFavouritesPaginated(
	ctx context.Context,
	userID string,
	params pagination.PaginationParams,
) ([]models.Asset, string, bool, error) {

	if userID == "" {
		return nil, "", false, fmt.Errorf("user ID is required")
	}

	// Parse cursor if provided
	var (
		cur    pagination.Cursor
		hasCur bool
	)
	if params.Cursor != "" {
		parsed, err := pagination.ParseCursor(params.Cursor)
		if err != nil {
			return nil, "", false, fmt.Errorf("invalid cursor: %w", err)
		}
		cur = parsed
		hasCur = true
	}

	// Get all favourites (cached list)
	allFavourites, err := s.GetFavourites(ctx, userID)
	if err != nil {
		return nil, "", false, err
	}

	// Sort newest-first: favourited_at DESC, asset_id DESC
	sort.Slice(allFavourites, func(i, j int) bool {
		ti := allFavourites[i].GetFavouritedAt()
		tj := allFavourites[j].GetFavouritedAt()
		if !ti.Equal(tj) {
			return ti.After(tj)
		}
		return allFavourites[i].GetAssetID() > allFavourites[j].GetAssetID()
	})

	// Apply keyset filter if cursor exists
	filtered := allFavourites
	if hasCur {
		tmp := make([]models.Asset, 0, len(allFavourites))
		for _, fav := range allFavourites {
			t := fav.GetFavouritedAt()

			// For DESC ordering, return items strictly "after" the cursor:
			// t < cur.t OR (t == cur.t AND id < cur.id)
			if t.Before(cur.FavouritedAt) ||
				(t.Equal(cur.FavouritedAt) && fav.GetAssetID() < cur.AssetID) {
				tmp = append(tmp, fav)
			}
		}
		filtered = tmp
	}

	// limit enforcement
	limit := params.Limit
	if limit <= 0 {
		limit = pagination.DefaultLimit
	}
	if limit > pagination.MaxLimit {
		limit = pagination.MaxLimit
	}

	if len(filtered) == 0 {
		return []models.Asset{}, "", false, nil
	}

	// Fetch limit+1 to compute hasMore
	hasMore := len(filtered) > limit
	end := limit
	if end > len(filtered) {
		end = len(filtered)
	}
	page := filtered[:end]

	// next cursor is based on last item in page
	var nextCursor string
	if hasMore && len(page) > 0 {
		last := page[len(page)-1]
		nextCursor = pagination.EncodeCursor(last.GetFavouritedAt(), last.GetAssetID())
	}

	return page, nextCursor, hasMore, nil
}

// UpdateDescription updates the description of a favourited asset
func (s *FavouritesService) UpdateDescription(ctx context.Context, userID, assetID, description string) error {
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}
	if assetID == "" {
		return fmt.Errorf("asset ID is required")
	}

	if err := s.storage.Favourites.UpdateDescription(ctx, userID, assetID, description); err != nil {
		return err
	}

	// Invalidate cache for this user
	cacheKey := fmt.Sprintf("favourites:user:%s", userID)
	s.cache.Delete(ctx, cacheKey)

	return nil
}

// GetFavourite retrieves a specific favourite
func (s *FavouritesService) GetFavourite(ctx context.Context, userID, assetID string) (models.Asset, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}
	if assetID == "" {
		return nil, fmt.Errorf("asset ID is required")
	}

	return s.storage.Favourites.GetFavourite(ctx, userID, assetID)
}

// GetCacheStats returns cache statistics for monitoring
func (s *FavouritesService) GetCacheStats() cache.CacheStats {
	return s.cache.GetStats()
}
