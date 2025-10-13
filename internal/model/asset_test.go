package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsset_GetBaseAsset(t *testing.T) {
	t.Run("returns Chart's BaseAsset when Chart is set", func(t *testing.T) {
		asset := &Asset{
			Chart: &Chart{
				BaseAsset: BaseAsset{
					ID:          "chart-1",
					Type:        AssetTypeChart,
					Description: "Test chart",
				},
				Title: "My Chart",
			},
		}

		base := asset.GetBaseAsset()

		assert.NotNil(t, base)
		assert.Equal(t, "chart-1", base.ID)
		assert.Equal(t, AssetTypeChart, base.Type)
		assert.Equal(t, "Test chart", base.Description)
	})

	t.Run("returns Insight's BaseAsset when Insight is set", func(t *testing.T) {
		asset := &Asset{
			Insight: &Insight{
				BaseAsset: BaseAsset{
					ID:          "insight-1",
					Type:        AssetTypeInsight,
					Description: "Test insight",
				},
				Title:   "My Insight",
				Content: "Some content",
			},
		}

		base := asset.GetBaseAsset()

		assert.NotNil(t, base)
		assert.Equal(t, "insight-1", base.ID)
		assert.Equal(t, AssetTypeInsight, base.Type)
	})

	t.Run("returns Audience's BaseAsset when Audience is set", func(t *testing.T) {
		asset := &Asset{
			Audience: &Audience{
				BaseAsset: BaseAsset{
					ID:   "audience-1",
					Type: AssetTypeAudience,
				},
				Gender:    "Male",
				Countries: []string{"US", "UK"},
				AgeMin:    18,
				AgeMax:    65,
			},
		}

		base := asset.GetBaseAsset()

		assert.NotNil(t, base)
		assert.Equal(t, "audience-1", base.ID)
		assert.Equal(t, AssetTypeAudience, base.Type)
	})

	t.Run("returns nil when no asset is set", func(t *testing.T) {
		asset := &Asset{}

		base := asset.GetBaseAsset()

		assert.Nil(t, base)
	})

	t.Run("prioritizes Chart when multiple assets set", func(t *testing.T) {
		asset := &Asset{
			Chart: &Chart{
				BaseAsset: BaseAsset{ID: "chart-1", Type: AssetTypeChart},
			},
			Insight: &Insight{
				BaseAsset: BaseAsset{ID: "insight-1", Type: AssetTypeInsight},
			},
		}

		base := asset.GetBaseAsset()

		assert.NotNil(t, base)
		assert.Equal(t, "chart-1", base.ID)
	})
}