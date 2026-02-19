package store

import (
	"context"
	"testing"

	"github.com/gwi/favourites-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFavouritesStore_AddFavourite(t *testing.T) {
	store := NewInMemoryFavouriteStore()
	ctx := context.Background()
	userID := "user1"

	chart := &models.Chart{
		BaseAsset: models.BaseAsset{AssetID: "chart1", Type: models.AssetTypeChart},
		Title:     "Sales Chart",
	}

	err := store.AddFavourite(ctx, userID, chart)
	require.NoError(t, err)
}

func TestFavouritesStore_AddFavourite_Duplicate(t *testing.T) {
	store := NewInMemoryFavouriteStore()
	ctx := context.Background()
	userID := "user1"

	chart := &models.Chart{
		BaseAsset: models.BaseAsset{AssetID: "chart1", Type: models.AssetTypeChart},
	}

	store.AddFavourite(ctx, userID, chart)
	err := store.AddFavourite(ctx, userID, chart) // Try again

	assert.ErrorIs(t, err, ErrConflict)
}

func TestFavouritesStore_GetUserFavourites(t *testing.T) {
	store := NewInMemoryFavouriteStore()
	ctx := context.Background()
	userID := "user1"

	chart := &models.Chart{BaseAsset: models.BaseAsset{AssetID: "chart1", Type: models.AssetTypeChart}}
	insight := &models.Insight{BaseAsset: models.BaseAsset{AssetID: "insight1", Type: models.AssetTypeInsight}}

	store.AddFavourite(ctx, userID, chart)
	store.AddFavourite(ctx, userID, insight)

	favourites, err := store.GetUserFavourites(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, favourites, 2)
}

func TestFavouritesStore_GetFavourite(t *testing.T) {
	store := NewInMemoryFavouriteStore()
	ctx := context.Background()
	userID := "user1"

	chart := &models.Chart{
		BaseAsset: models.BaseAsset{AssetID: "chart1", Type: models.AssetTypeChart},
		Title:     "Sales",
	}
	store.AddFavourite(ctx, userID, chart)

	favourite, err := store.GetFavourite(ctx, userID, "chart1")
	require.NoError(t, err)
	assert.Equal(t, "Sales", favourite.(*models.Chart).Title)
}

func TestFavouritesStore_RemoveFavourite(t *testing.T) {
	store := NewInMemoryFavouriteStore()
	ctx := context.Background()
	userID := "user1"

	chart := &models.Chart{BaseAsset: models.BaseAsset{AssetID: "chart1", Type: models.AssetTypeChart}}
	store.AddFavourite(ctx, userID, chart)

	err := store.RemoveFavourite(ctx, userID, "chart1")
	require.NoError(t, err)

	_, err = store.GetFavourite(ctx, userID, "chart1")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestFavouritesStore_UpdateDescription(t *testing.T) {
	store := NewInMemoryFavouriteStore()
	ctx := context.Background()
	userID := "user1"

	chart := &models.Chart{
		BaseAsset: models.BaseAsset{AssetID: "chart1", Type: models.AssetTypeChart, Description: "Old"},
	}
	store.AddFavourite(ctx, userID, chart)

	err := store.UpdateDescription(ctx, userID, "chart1", "New description")
	require.NoError(t, err)

	favourite, _ := store.GetFavourite(ctx, userID, "chart1")
	assert.Equal(t, "New description", favourite.GetDescription())
}

func TestFavouritesStore_UserIsolation(t *testing.T) {
	store := NewInMemoryFavouriteStore()
	ctx := context.Background()

	chart := &models.Chart{BaseAsset: models.BaseAsset{AssetID: "chart1", Type: models.AssetTypeChart}}

	store.AddFavourite(ctx, "user1", chart)
	store.AddFavourite(ctx, "user2", chart)

	favs1, _ := store.GetUserFavourites(ctx, "user1")
	favs2, _ := store.GetUserFavourites(ctx, "user2")

	assert.Len(t, favs1, 1)
	assert.Len(t, favs2, 1)
}
