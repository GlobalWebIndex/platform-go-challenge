package db_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/mr-polymath/platform-go-challenge/internal/db"
	"github.com/mr-polymath/platform-go-challenge/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestMemoryStore(t *testing.T) {
	s := db.NewMemoryStore()

	userID := uuid.New()
	asset := model.NewAsset(model.InsightType, "Insight 1", "Test insight")

	// Add
	err := s.AddFavourite(userID, &asset)
	assert.NoError(t, err)

	// Get
	items, err := s.GetAllFavourites(userID, db.PagedQuery{})
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, asset.Description, items[0].Description)

	// Update
	err = s.PatchFavourite(userID, asset.ID, "Updated")
	assert.NoError(t, err)

	updated, _ := s.GetAllFavourites(userID, db.PagedQuery{})
	assert.Equal(t, "Updated", updated[0].Description)

	// Remove
	err = s.RemoveFavourite(userID, asset.ID)
	assert.NoError(t, err)

	items, _ = s.GetAllFavourites(userID, db.PagedQuery{})
	assert.Len(t, items, 0)
}
