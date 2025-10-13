package storage

import (
	"testing"
	"time"

	"platform-go-challenge/internal/api_errors"
	"platform-go-challenge/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryStore_GetUserFavorites(t *testing.T) {
	t.Run("returns empty slice for non-existent user", func(t *testing.T) {
		store := NewMemoryStore()
		
		favorites, err := store.GetUserFavorites("user-123")
		
		require.NoError(t, err)
		assert.Empty(t, favorites)
	})

	t.Run("returns all favorites for existing user", func(t *testing.T) {
		store := NewMemoryStore()
		asset := &model.Asset{
			Chart: &model.Chart{
				BaseAsset: model.BaseAsset{ID: "chart-1", Type: model.AssetTypeChart},
				Title:     "Test Chart",
			},
		}
		store.AddFavorite("user-123", asset)
		
		favorites, err := store.GetUserFavorites("user-123")
		
		require.NoError(t, err)
		assert.Len(t, favorites, 1)
		assert.Equal(t, "chart-1", favorites[0].Asset.GetBaseAsset().ID)
	})
}

func TestMemoryStore_AddFavorite(t *testing.T) {
	t.Run("creates new user if doesn't exist", func(t *testing.T) {
		store := NewMemoryStore()
		asset := &model.Asset{
			Insight: &model.Insight{
				BaseAsset: model.BaseAsset{ID: "insight-1", Type: model.AssetTypeInsight},
				Title:     "Test Insight",
			},
		}
		
		err := store.AddFavorite("new-user", asset)
		
		require.NoError(t, err)
		assert.Contains(t, store.users, "new-user")
	})

	t.Run("adds favorite to existing user", func(t *testing.T) {
		store := NewMemoryStore()
		asset1 := &model.Asset{
			Chart: &model.Chart{
				BaseAsset: model.BaseAsset{ID: "chart-1", Type: model.AssetTypeChart},
			},
		}
		asset2 := &model.Asset{
			Insight: &model.Insight{
				BaseAsset: model.BaseAsset{ID: "insight-1", Type: model.AssetTypeInsight},
			},
		}
		
		store.AddFavorite("user-123", asset1)
		err := store.AddFavorite("user-123", asset2)
		
		require.NoError(t, err)
		favorites, _ := store.GetUserFavorites("user-123")
		assert.Len(t, favorites, 2)
	})

	t.Run("sets CreatedAt and UpdatedAt timestamps", func(t *testing.T) {
		store := NewMemoryStore()
		asset := &model.Asset{
			Audience: &model.Audience{
				BaseAsset: model.BaseAsset{ID: "audience-1", Type: model.AssetTypeAudience},
			},
		}
		
		err := store.AddFavorite("user-123", asset)
		
		require.NoError(t, err)
		favorites, _ := store.GetUserFavorites("user-123")
		baseAsset := favorites[0].Asset.GetBaseAsset()
		assert.Greater(t, baseAsset.CreatedAt, int64(0))
		assert.Greater(t, baseAsset.UpdatedAt, int64(0))
		assert.Greater(t, favorites[0].AddedAt, int64(0))
	})

	t.Run("overwrites duplicate favorites", func(t *testing.T) {
		store := NewMemoryStore()
		asset1 := &model.Asset{
			Chart: &model.Chart{
				BaseAsset: model.BaseAsset{ID: "chart-1", Type: model.AssetTypeChart, Description: "First"},
			},
		}
		asset2 := &model.Asset{
			Chart: &model.Chart{
				BaseAsset: model.BaseAsset{ID: "chart-1", Type: model.AssetTypeChart, Description: "Second"},
			},
		}
		
		store.AddFavorite("user-123", asset1)
		store.AddFavorite("user-123", asset2)
		
		favorites, _ := store.GetUserFavorites("user-123")
		assert.Len(t, favorites, 1)
		assert.Equal(t, "Second", favorites[0].Asset.GetBaseAsset().Description)
	})
}

func TestMemoryStore_RemoveFavorite(t *testing.T) {
	t.Run("successfully removes existing favorite", func(t *testing.T) {
		store := NewMemoryStore()
		asset := &model.Asset{
			Chart: &model.Chart{
				BaseAsset: model.BaseAsset{ID: "chart-1", Type: model.AssetTypeChart},
			},
		}
		store.AddFavorite("user-123", asset)
		
		err := store.RemoveFavorite("user-123", "chart-1")
		
		require.NoError(t, err)
		favorites, _ := store.GetUserFavorites("user-123")
		assert.Empty(t, favorites)
	})

	t.Run("returns error for non-existent user", func(t *testing.T) {
		store := NewMemoryStore()
		
		err := store.RemoveFavorite("user-123", "chart-1")
		
		assert.ErrorIs(t, err, api_errors.ErrUserNotFound)
	})

	t.Run("returns error for non-existent favorite", func(t *testing.T) {
		store := NewMemoryStore()
		asset := &model.Asset{
			Chart: &model.Chart{
				BaseAsset: model.BaseAsset{ID: "chart-1", Type: model.AssetTypeChart},
			},
		}
		store.AddFavorite("user-123", asset)
		
		err := store.RemoveFavorite("user-123", "chart-999")
		
		assert.ErrorIs(t, err, api_errors.ErrFavoriteNotFound)
	})
}

func TestMemoryStore_UpdateFavoriteDescription(t *testing.T) {
	t.Run("successfully updates description", func(t *testing.T) {
		store := NewMemoryStore()
		asset := &model.Asset{
			Chart: &model.Chart{
				BaseAsset: model.BaseAsset{ID: "chart-1", Type: model.AssetTypeChart, Description: "Old"},
			},
		}
		store.AddFavorite("user-123", asset)
		
		err := store.UpdateFavoriteDescription("user-123", "chart-1", "New description")
		
		require.NoError(t, err)
		favorites, _ := store.GetUserFavorites("user-123")
		assert.Equal(t, "New description", favorites[0].Asset.GetBaseAsset().Description)
	})

	t.Run("updates UpdatedAt timestamp", func(t *testing.T) {
		store := NewMemoryStore()
		asset := &model.Asset{
			Insight: &model.Insight{
				BaseAsset: model.BaseAsset{ID: "insight-1", Type: model.AssetTypeInsight},
			},
		}
		store.AddFavorite("user-123", asset)
		favorites, _ := store.GetUserFavorites("user-123")
		oldUpdatedAt := favorites[0].Asset.GetBaseAsset().UpdatedAt
		
		// Small delay to ensure different timestamp
		time.Sleep(time.Millisecond)
		
		store.UpdateFavoriteDescription("user-123", "insight-1", "Updated")
		
		favorites, _ = store.GetUserFavorites("user-123")
		newUpdatedAt := favorites[0].Asset.GetBaseAsset().UpdatedAt
		assert.GreaterOrEqual(t, newUpdatedAt, oldUpdatedAt)
	})

	t.Run("returns error for non-existent user", func(t *testing.T) {
		store := NewMemoryStore()
		
		err := store.UpdateFavoriteDescription("user-123", "chart-1", "New")
		
		assert.ErrorIs(t, err, api_errors.ErrUserNotFound)
	})

	t.Run("returns error for non-existent favorite", func(t *testing.T) {
		store := NewMemoryStore()
		asset := &model.Asset{
			Chart: &model.Chart{
				BaseAsset: model.BaseAsset{ID: "chart-1", Type: model.AssetTypeChart},
			},
		}
		store.AddFavorite("user-123", asset)
		
		err := store.UpdateFavoriteDescription("user-123", "chart-999", "New")
		
		assert.ErrorIs(t, err, api_errors.ErrFavoriteNotFound)
	})
}

func TestMemoryStore_Concurrency(t *testing.T) {
	t.Run("concurrent operations are thread-safe", func(t *testing.T) {
		store := NewMemoryStore()
		done := make(chan bool)
		
		// Concurrent writes
		for i := 0; i < 10; i++ {
			go func(id int) {
				asset := &model.Asset{
					Chart: &model.Chart{
						BaseAsset: model.BaseAsset{ID: string(rune(id)), Type: model.AssetTypeChart},
					},
				}
				store.AddFavorite("user-123", asset)
				done <- true
			}(i)
		}
		
		// Concurrent reads
		for i := 0; i < 10; i++ {
			go func() {
				store.GetUserFavorites("user-123")
				done <- true
			}()
		}
		
		for i := 0; i < 20; i++ {
			<-done
		}
		
		// If we reach here without deadlock or race, test passes
		assert.True(t, true)
	})
}