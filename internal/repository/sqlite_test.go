package repository

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"platform-go-challenge/internal/domain"
	"platform-go-challenge/internal/repository/testdata"
	"platform-go-challenge/internal/service/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	expectedFavsPerUser = 20
	maxDescriptionLen   = 500
)

func TestRepository(t *testing.T) {
	defer os.Remove(dbPath)

	repo, err := New()
	require.NoError(t, err, "failed to create repository")

	ctx := context.Background()

	// Load initial favorites for all users
	favs1, err := repo.ListFavorites(ctx, testdata.Users[0], 1, expectedFavsPerUser)
	require.NoError(t, err, "failed to get favorites for user 1")
	require.Len(t, favs1, expectedFavsPerUser, "invalid favs len for user 1")

	favs2, err := repo.ListFavorites(ctx, testdata.Users[1], 1, expectedFavsPerUser)
	require.NoError(t, err, "failed to get favorites for user 2")
	require.Len(t, favs2, expectedFavsPerUser, "invalid favs len for user 2")

	favs3, err := repo.ListFavorites(ctx, testdata.Users[2], 1, expectedFavsPerUser)
	require.NoError(t, err, "failed to get favorites for user 3")
	require.Len(t, favs3, expectedFavsPerUser, "invalid favs len for user 3")

	t.Run("GetAssets", func(t *testing.T) {
		t.Run("get and validate assets", func(t *testing.T) {
			var assetIDs []any
			for _, fav := range favs2 {
				assetIDs = append(assetIDs, fav.AssetID)
			}

			assets, err := repo.GetAssets(ctx, assetIDs)
			require.NoError(t, err, "failed to get assets")
			require.Len(t, assets, expectedFavsPerUser, "invalid assets len")

			for _, asset := range assets {
				// Test base asset interface methods
				assert.NotEmpty(t, asset.GetID(), "asset ID should not be empty")
				assert.NotEmpty(t, asset.GetDescription(), "asset description should not be empty")
				assert.NotEmpty(t, asset.GetCreatedAt(), "created at should not be empty")
				assert.NotEmpty(t, asset.GetUpdatedAt(), "updated at should not be empty")

				switch v := asset.(type) {
				case *types.Chart:
					assert.Equal(t, types.AssetTypeChart, v.GetType())
					assert.NotEmpty(t, v.Title)
					assert.NotNil(t, v.ChartData)

				case *types.Insight:
					assert.Equal(t, types.AssetTypeInsight, v.GetType())
					assert.NotEmpty(t, v.Content)

				case *types.Audience:
					assert.Equal(t, types.AssetTypeAudience, v.GetType())
					assert.NotEmpty(t, v.BirthCountry)
					assert.NotEmpty(t, v.Gender)
					require.NotEmpty(t, v.AgeGroups, "age groups should not be empty")

					for _, age := range v.AgeGroups {
						assert.GreaterOrEqual(t, age.Max, age.Min,
							"age group max should be >= min")
						assert.GreaterOrEqual(t, age.Min, 0,
							"age group min should be >= 0")
					}

					assert.GreaterOrEqual(t, v.SocialMediaHours, 0.0,
						"social media hours should be >= 0")
					assert.GreaterOrEqual(t, v.PurchasesLastMonth, 0,
						"purchases last month should be >= 0")

				default:
					t.Errorf("unexpected asset type: %T", v)
				}
			}
		})

		t.Run("with non-existent assets", func(t *testing.T) {
			nonExistentUUID1 := "123e4567-e89b-12d3-a456-426614174000"
			nonExistentUUID2 := "123e4567-e89b-12d3-a456-426614174001"

			assets, err := repo.GetAssets(ctx, []any{nonExistentUUID1, nonExistentUUID2})
			require.NoError(t, err)
			assert.Empty(t, assets)
		})
	})

	t.Run("UpdateAsset", func(t *testing.T) {
		t.Run("update with valid description", func(t *testing.T) {
			assets, err := repo.GetAssets(ctx, []any{favs1[0].AssetID})
			require.NoError(t, err, "failed to get asset")
			require.Len(t, assets, 1)

			asset := assets[0]
			originalDesc := asset.GetDescription()
			newDesc := "Updated description " + time.Now().String()

			err = asset.SetDescription(newDesc)
			require.NoError(t, err, "failed to set description")

			err = repo.UpdateAsset(ctx, asset)
			require.NoError(t, err, "failed to update asset")

			// Verify update and timestamps
			updated, err := repo.GetAssets(ctx, []any{asset.GetID()})
			require.NoError(t, err, "failed to get updated asset")
			require.Len(t, updated, 1)

			assert.Equal(t, newDesc, updated[0].GetDescription())
			assert.NotEqual(t, originalDesc, updated[0].GetDescription())
			assert.NotEqual(t, asset.GetUpdatedAt(), updated[0].GetUpdatedAt(),
				"UpdatedAt should change")
		})

		t.Run("update with too long description", func(t *testing.T) {
			assets, err := repo.GetAssets(ctx, []any{favs1[0].AssetID})
			require.NoError(t, err)
			require.Len(t, assets, 1)

			asset := assets[0]
			longDesc := strings.Repeat("a", maxDescriptionLen+1)

			err = asset.SetDescription(longDesc)
			assert.ErrorIs(t, err, types.ErrDescriptionTooLong)
		})
	})

	t.Run("ListFavorites", func(t *testing.T) {
		t.Run("pagination", func(t *testing.T) {
			testCases := []struct {
				name     string
				page     int
				limit    int
				expected int
			}{
				{"full page", 1, 10, 10},
				{"partial page", 2, 15, 5},
				{"exact limit", 1, expectedFavsPerUser, expectedFavsPerUser},
				{"beyond available", 3, 10, 0},
				{"zero limit", 1, 0, 0},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					favs, err := repo.ListFavorites(ctx, testdata.Users[0], tc.page, tc.limit)
					require.NoError(t, err)
					assert.Len(t, favs, tc.expected)

					for _, fav := range favs {
						assert.NotEmpty(t, fav.UserID)
						assert.NotEmpty(t, fav.AssetID)
						assert.NotEmpty(t, fav.CreatedAt)
					}
				})
			}
		})

		t.Run("with empty user ID", func(t *testing.T) {
			favs, err := repo.ListFavorites(ctx, "", 1, 10)
			require.NoError(t, err)
			assert.Empty(t, favs)
		})
	})

	t.Run("Favorites", func(t *testing.T) {
		t.Run("add favorite successfully", func(t *testing.T) {
			fav := favs2[5]
			fav.UserID = testdata.Users[0]
			fav.CreatedAt = time.Now()

			err := repo.AddFavorite(ctx, fav)
			require.NoError(t, err, "failed to add favorite")

			favs, err := repo.ListFavorites(ctx, testdata.Users[0], 1, 1000)
			require.NoError(t, err, "failed to list favorites")

			var found bool
			for _, v := range favs {
				if v.AssetID == fav.AssetID {
					found = true
					assert.Equal(t, fav.UserID, v.UserID)
					assert.WithinDuration(t, fav.CreatedAt, v.CreatedAt, time.Second)
					break
				}
			}
			assert.True(t, found, "favorite not found after adding")
		})

		t.Run("prevent duplicate favorites", func(t *testing.T) {
			err := repo.AddFavorite(ctx, favs3[1])
			assert.True(t, errors.Is(err, domain.ErrAlreadyExists))
		})

		t.Run("remove favorite successfully", func(t *testing.T) {
			err := repo.RemoveFavorite(ctx, favs2[5].UserID, favs2[5].AssetID)
			require.NoError(t, err, "failed to remove favorite")

			favs, err := repo.ListFavorites(ctx, favs2[5].UserID, 1, 1000)
			require.NoError(t, err, "failed to list favorites")

			for _, v := range favs {
				assert.NotEqual(t, favs2[5].AssetID, v.AssetID,
					"found removed favorite in list")
			}
		})

		t.Run("remove non-existent favorite", func(t *testing.T) {
			nonExistentUserID := "123e4567-e89b-12d3-a456-426614174000"
			nonExistentAssetID := "123e4567-e89b-12d3-a456-426614174001"

			err := repo.RemoveFavorite(ctx, nonExistentUserID, nonExistentAssetID)
			assert.ErrorIs(t, err, domain.ErrNotFound)
		})
	})
}
