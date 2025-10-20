package db_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/mr-polymath/platform-go-challenge/internal/db"
	"github.com/mr-polymath/platform-go-challenge/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestMemoryStore_BasicCRUD(t *testing.T) {
	s := db.NewMemoryStore()

	userID := uuid.New()
	asset := model.NewAsset(model.InsightType, "Insight 1", "Test insight")

	// Add
	err := s.AddFavourite(userID, asset)
	assert.NoError(t, err)

	// Get
	items, err := s.GetAllFavourites(userID, db.PagedQuery{})
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, asset.Description, items[0].Description)

	// Update
	err = s.PatchFavourite(userID, items[0].ID, "Updated")
	assert.NoError(t, err)

	updated, _ := s.GetAllFavourites(userID, db.PagedQuery{})
	assert.Equal(t, "Updated", updated[0].Description)

	// Remove
	err = s.RemoveFavourite(userID, items[0].ID)
	assert.NoError(t, err)

	items, _ = s.GetAllFavourites(userID, db.PagedQuery{})
	assert.Len(t, items, 0)
}
