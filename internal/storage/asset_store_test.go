package storage

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"platform-go-challenge/internal/common"
)

// MockAsset implements the Asset interface for testing
type MockAsset struct {
	ID          string
	Description string
	Created     time.Time
	Updated     time.Time
}

func (m *MockAsset) GetID() string                   { return m.ID }
func (m *MockAsset) GetDescription() string          { return m.Description }
func (m *MockAsset) GetCreatedTimestamp() time.Time  { return m.Created }
func (m *MockAsset) GetUpdatedTimestamp() time.Time  { return m.Updated }
func (m *MockAsset) SetUpdatedTimestamp(t time.Time) { m.Updated = t }

func newMockAsset(id string) *MockAsset {
	now := time.Now()
	return &MockAsset{
		ID:      id,
		Created: now,
		Updated: now,
	}
}

func TestFavoritesStore(t *testing.T) {
	store := NewFavoritesStore()
	userID := uuid.New()

	t.Run("AddFavorite", func(t *testing.T) {
		asset := newMockAsset("asset1")
		err := store.AddFavorite(userID, asset)
		assert.NoError(t, err)

		// Adding the same asset should fail
		err = store.AddFavorite(userID, asset)
		assert.Error(t, err)
		assert.EqualError(t, err, "asset already in favorites")
	})

	t.Run("GetUserFavorites", func(t *testing.T) {
		assets, err := store.GetUserFavorites(userID)
		assert.NoError(t, err)
		assert.Len(t, assets, 1)
		assert.Equal(t, "asset1", assets[0].GetID())

		// Non-existent user should return empty list
		assets, err = store.GetUserFavorites(uuid.New())
		assert.NoError(t, err)
		assert.Empty(t, assets)
	})

	t.Run("GetFavoriteByAssetID", func(t *testing.T) {
		asset, err := store.GetFavoriteByAssetID(userID, "asset1")
		assert.NoError(t, err)
		assert.Equal(t, "asset1", asset.GetID())

		// Non-existent asset
		_, err = store.GetFavoriteByAssetID(userID, "unknown")
		assert.Error(t, err)
		assert.EqualError(t, err, common.ErrAssetNotFound)
	})

	t.Run("RemoveFavorite", func(t *testing.T) {
		err := store.RemoveFavorite(userID, "asset1")
		assert.NoError(t, err)

		// Removing again should fail
		err = store.RemoveFavorite(userID, "asset1")
		assert.Error(t, err)
		assert.EqualError(t, err, common.ErrAssetNotFound)
	})

	t.Run("RemoveUserFavorites", func(t *testing.T) {
		// Add two assets
		a1 := newMockAsset("a1")
		a2 := newMockAsset("a2")
		store.AddFavorite(userID, a1)
		store.AddFavorite(userID, a2)

		store.RemoveUserFavorites(userID)

		// Should return empty favorites
		assets, err := store.GetUserFavorites(userID)
		assert.NoError(t, err)
		assert.Empty(t, assets)
	})

	t.Run("UpdateFavorite", func(t *testing.T) {
		asset := newMockAsset("assetX")
		store.AddFavorite(userID, asset)

		updated := &MockAsset{
			ID:          "assetX",
			Description: "updated description",
		}

		err := store.UpdateFavorite(userID, updated)
		assert.NoError(t, err)

		assetAfter, _ := store.GetFavoriteByAssetID(userID, "assetX")
		assert.Equal(t, "updated description", assetAfter.GetDescription())
		assert.NotZero(t, assetAfter.GetUpdatedTimestamp())
	})

	t.Run("UpdateFavorite errors", func(t *testing.T) {
		// Non-existent user
		err := store.UpdateFavorite(uuid.New(), &MockAsset{ID: "x"})
		assert.Error(t, err)
		assert.EqualError(t, err, common.ErrUserNotFound)

		// Non-existent asset
		err = store.UpdateFavorite(userID, &MockAsset{ID: "unknown"})
		assert.Error(t, err)
	})
}

func TestFavoritesStoreConcurrency(t *testing.T) {
	store := NewFavoritesStore()
	userID := uuid.New()
	numAssets := 50

	var wg sync.WaitGroup

	// Concurrently add favorites
	for i := 0; i < numAssets; i++ {
		i := i // capture loop variable
		wg.Add(1)
		go func() {
			defer wg.Done()
			assetID := "asset" + strconv.Itoa(i)
			asset := newMockAsset(assetID)
			_ = store.AddFavorite(userID, asset)
		}()
	}

	wg.Wait()

	// Verify all assets were added
	assets, err := store.GetUserFavorites(userID)
	assert.NoError(t, err)
	assert.Len(t, assets, numAssets)

	// Concurrently remove favorites
	wg = sync.WaitGroup{}
	for i := 0; i < numAssets; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			assetID := "asset" + strconv.Itoa(i)
			_ = store.RemoveFavorite(userID, assetID)
		}()
	}

	wg.Wait()

	// Verify store is empty for this user
	assets, err = store.GetUserFavorites(userID)
	assert.NoError(t, err)
	assert.Len(t, assets, 0)
}
