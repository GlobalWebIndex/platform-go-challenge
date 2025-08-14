package unit

import (
	"context"
	"testing"

	"gwi-favorites-service/internal/domain"
	"gwi-favorites-service/internal/repository/memory"
	"gwi-favorites-service/internal/service"
	"gwi-favorites-service/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFavoritesService_AddFavorite(t *testing.T) {
	// Setup
	repo := memory.NewRepository()
	log := logger.NewLogger()
	svc := service.NewFavoritesService(repo, log)
	ctx := context.Background()

	// Create test user
	user := domain.NewUser("user1", "test@example.com", "Test User")
	require.NoError(t, repo.CreateUser(user))

	// Create test asset
	asset := domain.NewChart("chart1", "Test Chart", "X", "Y", "Test Description", nil)

	// Test
	err := svc.AddFavorite(ctx, "user1", asset)
	assert.NoError(t, err)

	// Verify
	favorites, err := svc.GetUserFavorites(ctx, "user1", 10, 0)
	assert.NoError(t, err)
	assert.Len(t, favorites, 1)
	assert.Equal(t, "chart1", favorites[0].AssetID)
}

func TestFavoritesService_GetUserFavorites(t *testing.T) {
	// Setup
	repo := memory.NewRepository()
	log := logger.NewLogger()
	svc := service.NewFavoritesService(repo, log)
	ctx := context.Background()

	// Create test user
	user := domain.NewUser("user1", "test@example.com", "Test User")
	require.NoError(t, repo.CreateUser(user))

	// Test empty favorites
	favorites, err := svc.GetUserFavorites(ctx, "user1", 10, 0)
	assert.NoError(t, err)
	assert.Len(t, favorites, 0)

	// Add some favorites
	asset1 := domain.NewChart("chart1", "Chart 1", "X", "Y", "Description 1", nil)
	asset2 := domain.NewInsight("insight1", "Test insight", "Description 2", nil, "")

	require.NoError(t, svc.AddFavorite(ctx, "user1", asset1))
	require.NoError(t, svc.AddFavorite(ctx, "user1", asset2))

	// Test with favorites
	favorites, err = svc.GetUserFavorites(ctx, "user1", 10, 0)
	assert.NoError(t, err)
	assert.Len(t, favorites, 2)
}
