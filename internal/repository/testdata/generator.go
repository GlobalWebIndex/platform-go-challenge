package testdata

import (
	"encoding/json"
	"platform-go-challenge/internal/service/types"
	"time"

	"github.com/google/uuid"
)

type TestDataGenerator struct{}

func NewGenerator() *TestDataGenerator {
	return &TestDataGenerator{}
}

func (g *TestDataGenerator) GenerateChart() *types.Chart {
	chart := &types.Chart{
		BaseAsset: types.BaseAsset{
			ID:          uuid.New().String(),
			Type:        types.AssetTypeChart,
			Description: "Test chart description",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Title: "Sample Chart",
		AxisX: "Time",
		AxisY: "Value",
		ChartData: map[string]interface{}{
			"values": []float64{1.0, 2.0, 3.0, 4.0},
			"labels": []string{"Q1", "Q2", "Q3", "Q4"},
		},
	}
	return chart
}

func (g *TestDataGenerator) GenerateInsight() *types.Insight {
	return &types.Insight{
		BaseAsset: types.BaseAsset{
			ID:          uuid.New().String(),
			Type:        types.AssetTypeInsight,
			Description: "Test insight description",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Content: "40% of millennials spend more than 3 hours on social media daily",
	}
}

func (g *TestDataGenerator) GenerateAudience() *types.Audience {
	return &types.Audience{
		BaseAsset: types.BaseAsset{
			ID:          uuid.New().String(),
			Type:        types.AssetTypeAudience,
			Description: "Test audience description",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Gender:             "Male",
		BirthCountry:       "USA",
		AgeGroups:          []types.AgeGroup{{Min: 24, Max: 35}},
		SocialMediaHours:   3.5,
		PurchasesLastMonth: 5,
	}
}

func (g *TestDataGenerator) GenerateAssetJSON(asset types.Asset) (string, error) {
	data, err := json.Marshal(asset)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
