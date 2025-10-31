package domain_test

import (
	"testing"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper constructors for assets
func newAudience(id string) *domain.Audience {
	return &domain.Audience{AssetBase: domain.AssetBase{ID: id, Type: domain.AssetTypeAudience}}
}
func newChart(id string) *domain.Chart {
	return &domain.Chart{AssetBase: domain.AssetBase{ID: id, Type: domain.AssetTypeChart}}
}
func newInsight(id string) *domain.Insight {
	return &domain.Insight{AssetBase: domain.AssetBase{ID: id, Type: domain.AssetTypeInsight}}
}

func TestFavourite_GetAsset(t *testing.T) {
	tests := []struct {
		name      string
		favourite *domain.Favourite
		wantID    string
		wantType  domain.AssetType
		wantNil   bool
	}{
		{"audience asset", &domain.Favourite{AssetType: domain.AssetTypeAudience, Audience: newAudience("aud-1")}, "aud-1", domain.AssetTypeAudience, false},
		{"chart asset", &domain.Favourite{AssetType: domain.AssetTypeChart, Chart: newChart("chart-1")}, "chart-1", domain.AssetTypeChart, false},
		{"insight asset", &domain.Favourite{AssetType: domain.AssetTypeInsight, Insight: newInsight("insight-1")}, "insight-1", domain.AssetTypeInsight, false},
		{"unknown type", &domain.Favourite{AssetType: domain.AssetType(-12)}, "", 0, true},
		{"nil asset for type", &domain.Favourite{AssetType: domain.AssetTypeAudience, Audience: nil}, "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			got := tt.favourite.GetAsset()

			// Assert
			if tt.wantNil {
				assert.Nil(t, got)
			} else {
				require.NotNil(t, got)
				assert.Equal(t, tt.wantID, got.GetID())
				assert.Equal(t, tt.wantType, got.GetType())
			}
		})
	}
}

func TestFavourite_SetAsset(t *testing.T) {
	tests := []struct {
		name      string
		favourite *domain.Favourite
		asset     domain.Asset
		wantID    string
		wantType  domain.AssetType
	}{
		{"set audience", &domain.Favourite{UserID: "u1"}, newAudience("aud-123"), "aud-123", domain.AssetTypeAudience},
		{"set chart", &domain.Favourite{UserID: "u1"}, newChart("chart-456"), "chart-456", domain.AssetTypeChart},
		{"set insight", &domain.Favourite{UserID: "u1"}, newInsight("insight-789"), "insight-789", domain.AssetTypeInsight},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := tt.favourite.SetAsset(tt.asset)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, tt.wantID, tt.favourite.AssetID)
			assert.Equal(t, tt.wantType, tt.favourite.AssetType)

			// Ensure only the matching field is set
			switch tt.wantType {
			case domain.AssetTypeAudience:
				require.NotNil(t, tt.favourite.Audience)
				assert.Nil(t, tt.favourite.Chart)
				assert.Nil(t, tt.favourite.Insight)
			case domain.AssetTypeChart:
				require.NotNil(t, tt.favourite.Chart)
				assert.Nil(t, tt.favourite.Audience)
				assert.Nil(t, tt.favourite.Insight)
			case domain.AssetTypeInsight:
				require.NotNil(t, tt.favourite.Insight)
				assert.Nil(t, tt.favourite.Audience)
				assert.Nil(t, tt.favourite.Chart)
			}
		})
	}
}

func TestFavourite_Integration_GetAfterSet(t *testing.T) {
	assets := []domain.Asset{
		newAudience("aud-int"),
		newChart("chart-int"),
		newInsight("insight-int"),
	}

	for _, asset := range assets {
		t.Run(asset.GetID(), func(t *testing.T) {
			// Arrange
			fav := &domain.Favourite{UserID: "test-user"}

			// Act
			err := fav.SetAsset(asset)
			require.NoError(t, err)
			got := fav.GetAsset()

			// Assert
			require.NotNil(t, got)
			assert.Equal(t, asset.GetID(), got.GetID())
			assert.Equal(t, asset.GetType(), got.GetType())
			assert.Equal(t, asset, got)
		})
	}
}
