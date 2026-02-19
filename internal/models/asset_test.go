package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalAsset_Chart(t *testing.T) {
	data := []byte(`{
		"asset_id": "asset-123",
		"type": "chart",
		"description": "Sales chart",
		"title": "Q4 Sales",
		"x_axis_title": "Month",
		"y_axis_title": "Revenue",
		"data": [{"month": "Jan", "revenue": 1000}]
	}`)

	asset, err := UnmarshalAsset(data)
	require.NoError(t, err)

	chart, ok := asset.(*Chart)
	require.True(t, ok)

	assert.Equal(t, "asset-123", chart.GetAssetID())
	assert.Equal(t, AssetTypeChart, chart.GetType())
	assert.Equal(t, "Sales chart", chart.GetDescription())
	assert.Equal(t, "Q4 Sales", chart.Title)
	assert.Equal(t, "Month", chart.XAxisTitle)
	assert.Equal(t, "Revenue", chart.YAxisTitle)
}

func TestUnmarshalAsset_Insight(t *testing.T) {
	data := []byte(`{
		"asset_id": "asset-456",
		"type": "insight",
		"description": "An insight",
		"text": "40% of millennials spend more than 3 hours on social media"
	}`)

	asset, err := UnmarshalAsset(data)
	require.NoError(t, err)

	insight, ok := asset.(*Insight)
	require.True(t, ok)

	assert.Equal(t, "asset-456", insight.GetAssetID())
	assert.Equal(t, AssetTypeInsight, insight.GetType())
	assert.Equal(t, "An insight", insight.GetDescription())
	assert.Equal(t, "40% of millennials spend more than 3 hours on social media", insight.Text)
}

func TestUnmarshalAsset_Audience(t *testing.T) {
	data := []byte(`{
		"asset_id": "asset-789",
		"type": "audience",
		"description": "Audience segment",
		"name": "Young Male Social Media Users",
		"gender": "male",
		"birth_country": "US",
		"age_group": {"min": 24, "max": 35},
		"daily_social_media_hours": 3,
		"purchases_last_month": 5
	}`)

	asset, err := UnmarshalAsset(data)
	require.NoError(t, err)

	audience, ok := asset.(*Audience)
	require.True(t, ok)

	assert.Equal(t, "asset-789", audience.GetAssetID())
	assert.Equal(t, AssetTypeAudience, audience.GetType())
	assert.Equal(t, "Audience segment", audience.GetDescription())
	assert.Equal(t, "Young Male Social Media Users", audience.Name)
	assert.Equal(t, GenderMale, audience.Gender)
	assert.Equal(t, 24, audience.AgeGroup.Min)
	assert.Equal(t, 35, audience.AgeGroup.Max)
}

func TestUnmarshalAsset_InvalidType(t *testing.T) {
	data := []byte(`{"asset_id":"x","type":"unknown"}`)
	_, err := UnmarshalAsset(data)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown asset type")
}

func TestUnmarshalAsset_MissingAssetID(t *testing.T) {
	data := []byte(`{
		"type": "chart",
		"title": "Q4 Sales",
		"data": []
	}`)
	_, err := UnmarshalAsset(data)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field: asset_id")
}

func TestUnmarshalAsset_MissingType(t *testing.T) {
	data := []byte(`{"asset_id":"x"}`)
	_, err := UnmarshalAsset(data)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field: type")
}

func TestMarshalAsset_RoundTrip(t *testing.T) {
	original := &Chart{
		BaseAsset: BaseAsset{
			AssetID:     "asset-123",
			Type:        AssetTypeChart,
			Description: "desc",
		},
		Title: "Sales Chart",
		Data:  []map[string]interface{}{},
	}

	// Marshal
	data, err := MarshalAsset(original)
	require.NoError(t, err)

	// Unmarshal
	asset, err := UnmarshalAsset(data)
	require.NoError(t, err)

	chart, ok := asset.(*Chart)
	require.True(t, ok)

	assert.Equal(t, original.AssetID, chart.AssetID)
	assert.Equal(t, original.Title, chart.Title)
	assert.Equal(t, original.Description, chart.Description)
}
