package service

import (
	"context"
	"testing"
	"time"

	"github.com/gwi/favourites-api/internal/cache"
	"github.com/gwi/favourites-api/internal/models"
	"github.com/gwi/favourites-api/internal/pagination"
	"github.com/gwi/favourites-api/internal/store"
	"github.com/stretchr/testify/require"
)

// spyCache is a tiny test just enough to verify behavior.
type spyCache struct {
	data     map[string]any
	deleted  []string
	setCalls int
	getCalls int
	lastTTL  time.Duration
	stats    cache.CacheStats
}

func newSpyCache() *spyCache {
	return &spyCache{data: make(map[string]any)}
}

func (s *spyCache) Get(ctx context.Context, key string) (any, bool) {
	s.getCalls++
	v, ok := s.data[key]
	if ok {
		s.stats.Hits++
		return v, true
	}
	s.stats.Misses++
	return nil, false
}

func (s *spyCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	s.setCalls++
	s.lastTTL = ttl
	s.stats.Sets++
	s.data[key] = value
	return nil
}

func (s *spyCache) Delete(ctx context.Context, key string) error {
	s.stats.Deletes++
	s.deleted = append(s.deleted, key)
	delete(s.data, key)
	return nil
}

func (s *spyCache) Clear(ctx context.Context) error {
	s.data = make(map[string]any)
	return nil
}

func (s *spyCache) GetStats() cache.CacheStats {
	out := s.stats
	out.Size = len(s.data)
	return out
}

func TestFavouritesService_AddFavourite_Success(t *testing.T) {
	ctx := context.Background()
	st := store.NewStorage()
	c := newSpyCache()
	ttl := 5 * time.Minute
	svc := NewFavouritesService(st, c, ttl)

	userID := "user-123"

	chart := &models.Chart{
		BaseAsset: models.BaseAsset{
			AssetID:     "asset-1",
			Type:        models.AssetTypeChart,
			Description: "desc",
			// FavouritedAt intentionally zero -> service should set
		},
		Title: "Test",
		Data:  []map[string]any{},
	}

	err := svc.AddFavourite(ctx, userID, chart)
	require.NoError(t, err)

	require.False(t, chart.BaseAsset.FavouritedAt.IsZero(), "service should set favourited_at if missing")
	require.Equal(t, time.UTC, chart.BaseAsset.FavouritedAt.Location(), "favourited_at should be UTC")

	// Ensure it was stored
	favs, err := st.Favourites.GetUserFavourites(ctx, userID)
	require.NoError(t, err)
	require.Len(t, favs, 1)

	// Ensure cache invalidation happened
	require.Contains(t, c.deleted, "favourites:user:"+userID)
}

func TestFavouritesService_AddFavourite_MissingUserID(t *testing.T) {
	ctx := context.Background()
	st := store.NewStorage()
	c := newSpyCache()
	ttl := 5 * time.Minute
	svc := NewFavouritesService(st, c, ttl)

	chart := &models.Chart{
		BaseAsset: models.BaseAsset{
			AssetID: "asset-1",
			Type:    models.AssetTypeChart,
		},
		Title: "Test",
		Data:  []map[string]any{},
	}

	err := svc.AddFavourite(ctx, "", chart)
	require.Error(t, err)
	require.Contains(t, err.Error(), "user ID is required")
}

func TestFavouritesService_AddFavourite_MissingAssetID(t *testing.T) {
	ctx := context.Background()
	st := store.NewStorage()
	c := newSpyCache()
	ttl := 5 * time.Minute
	svc := NewFavouritesService(st, c, ttl)

	userID := "user-123"

	chart := &models.Chart{
		BaseAsset: models.BaseAsset{
			AssetID: "", // missing external asset id
			Type:    models.AssetTypeChart,
		},
		Title: "Test",
		Data:  []map[string]any{},
	}

	err := svc.AddFavourite(ctx, userID, chart)
	require.Error(t, err)

	require.Contains(t, err.Error(), "asset_id")
}

func TestFavouritesService_GetFavouritesPaginated_KeysetPagination(t *testing.T) {
	ctx := context.Background()
	st := store.NewStorage()
	c := newSpyCache()
	ttl := 5 * time.Minute
	svc := NewFavouritesService(st, c, ttl)

	userID := "user-123"

	// Two timestamps to force tie-breaker behavior
	t1 := time.Date(2026, 2, 18, 10, 0, 0, 0, time.UTC) // newer
	t2 := time.Date(2026, 2, 18, 9, 0, 0, 0, time.UTC)  // older

	// Create favourites with controlled favourited_at + asset_id
	addInsight := func(assetID string, favAt time.Time) {
		ins := &models.Insight{
			BaseAsset: models.BaseAsset{
				AssetID:      assetID,
				Type:         models.AssetTypeInsight,
				Description:  "d",
				FavouritedAt: favAt,
			},
			Text: "x",
		}
		require.NoError(t, svc.AddFavourite(ctx, userID, ins))
	}

	// Expected order (DESC by favourited_at, then asset_id):
	// (t1, "c"), (t1, "b"), (t1, "a"), (t2, "b2"), (t2, "a2")
	addInsight("a", t1)
	addInsight("b", t1)
	addInsight("c", t1)
	addInsight("a2", t2)
	addInsight("b2", t2)

	// Page 1: limit 2 => ["c","b"]
	p1, next1, hasMore1, err := svc.GetFavouritesPaginated(ctx, userID, pagination.PaginationParams{
		Cursor: "",
		Limit:  2,
	})
	require.NoError(t, err)
	require.True(t, hasMore1)
	require.Len(t, p1, 2)
	require.Equal(t, "c", p1[0].GetAssetID())
	require.Equal(t, "b", p1[1].GetAssetID())
	require.NotEmpty(t, next1)

	// Cursor should encode last element ("b") of page 1
	cur1, err := pagination.ParseCursor(next1)
	require.NoError(t, err)
	require.True(t, cur1.FavouritedAt.Equal(t1))
	require.Equal(t, "b", cur1.AssetID)

	// Page 2 using cursor => ["a","b2"]
	p2, next2, hasMore2, err := svc.GetFavouritesPaginated(ctx, userID, pagination.PaginationParams{
		Cursor: next1,
		Limit:  2,
	})
	require.NoError(t, err)
	require.True(t, hasMore2)
	require.Len(t, p2, 2)
	require.Equal(t, "a", p2[0].GetAssetID())
	require.Equal(t, "b2", p2[1].GetAssetID())
	require.NotEmpty(t, next2)

	// Page 3 => ["a2"] and then end
	p3, next3, hasMore3, err := svc.GetFavouritesPaginated(ctx, userID, pagination.PaginationParams{
		Cursor: next2,
		Limit:  2,
	})
	require.NoError(t, err)
	require.False(t, hasMore3)
	require.Len(t, p3, 1)
	require.Equal(t, "a2", p3[0].GetAssetID())
	require.Empty(t, next3)

	// Ensure no duplicates across pages
	seen := map[string]bool{}
	for _, page := range [][]models.Asset{p1, p2, p3} {
		for _, a := range page {
			if seen[a.GetAssetID()] {
				t.Fatalf("duplicate across pages: %s", a.GetAssetID())
			}
			seen[a.GetAssetID()] = true
		}
	}
}

func TestFavouritesService_GetFavouritesPaginated_InvalidCursor(t *testing.T) {
	ctx := context.Background()
	st := store.NewStorage()
	c := newSpyCache()
	ttl := 5 * time.Minute
	svc := NewFavouritesService(st, c, ttl)

	userID := "user-123"

	// Add one favourite so list isn't empty
	ins := &models.Insight{
		BaseAsset: models.BaseAsset{
			AssetID:      "a",
			Type:         models.AssetTypeInsight,
			Description:  "d",
			FavouritedAt: time.Now().UTC(),
		},
		Text: "x",
	}
	require.NoError(t, svc.AddFavourite(ctx, userID, ins))

	_, _, _, err := svc.GetFavouritesPaginated(ctx, userID, pagination.PaginationParams{
		Cursor: "not-a-valid-cursor",
		Limit:  2,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid cursor")
}

func TestFavouritesService_RemoveFavourite_InvalidatesCache(t *testing.T) {
	ctx := context.Background()
	st := store.NewStorage()
	c := newSpyCache()
	ttl := 5 * time.Minute
	svc := NewFavouritesService(st, c, ttl)

	userID := "user-123"

	chart := &models.Chart{
		BaseAsset: models.BaseAsset{
			AssetID:      "asset-1",
			Type:         models.AssetTypeChart,
			Description:  "desc",
			FavouritedAt: time.Now().UTC(),
		},
		Title: "Chart",
		Data:  []map[string]any{},
	}

	require.NoError(t, svc.AddFavourite(ctx, userID, chart))

	// Prime cache to make invalidation meaningful
	_, _ = svc.GetFavourites(ctx, userID)
	require.Equal(t, 1, c.setCalls)
	require.Equal(t, ttl, c.lastTTL)

	err := svc.RemoveFavourite(ctx, userID, "asset-1")
	require.NoError(t, err)

	require.Contains(t, c.deleted, "favourites:user:"+userID, "remove should invalidate cache")
}
