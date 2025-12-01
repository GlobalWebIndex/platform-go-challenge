package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestCreateAsset tests full asset creation of all asset types
func TestCreateAsset(t *testing.T) {
	tests := []struct {
		name          string
		jsonData      string
		expectedID    string
		expectedTitle string // only relevant for Chart
	}{
		{
			name: "Chart asset",
			jsonData: `{
				"type": "chart",
				"id": "chart1",
				"title": "Sales",
				"x_axis_title": "Month",
				"y_axis_title": "Revenue",
				"data": []
			}`,
			expectedID:    "chart1",
			expectedTitle: "Sales",
		},
		{
			name: "Insight asset",
			jsonData: `{
				"type": "insight",
				"id": "insight1",
				"insight_text": "Important insight"
			}`,
			expectedID: "insight1",
		},
		{
			name: "Audience asset",
			jsonData: `{
				"type": "audience",
				"id": "aud1",
				"gender": "female",
				"birth_country": "USA",
				"age_group": "25-34",
				"daily_sm_hours": 2,
				"monthly_purchases": 5
			}`,
			expectedID: "aud1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asset, err := CreateAsset([]byte(tt.jsonData))
			assert.NoError(t, err, "CreateAsset should not return error")
			assert.Equal(t, tt.expectedID, asset.GetID(), "ID mismatch")
			assert.False(t, asset.GetCreatedTimestamp().IsZero(), "CreatedAt not set")
			assert.False(t, asset.GetUpdatedTimestamp().IsZero(), "UpdatedAt not set")

			if chart, ok := asset.(*Chart); ok {
				assert.Equal(t, tt.expectedTitle, chart.Title, "Chart title mismatch")
			}
		})
	}
}

// TestCreatePartialAsset tests partial asset updates
func TestCreatePartialAsset(t *testing.T) {
	tests := []struct {
		name          string
		jsonData      string
		partialUpdate string
		expectedID    string
		expectedTitle string // only for Chart
	}{
		{
			name: "Chart asset partial update",
			jsonData: `{
				"type": "chart",
				"id": "chart1",
				"title": "Sales",
				"x_axis_title": "Month",
				"y_axis_title": "Revenue",
				"data": [{"label":"Jan","value":100}]
			}`,
			partialUpdate: `{
				"title": "Updated Sales",
				"x_axis_title": "Quarter"
			}`,
			expectedID:    "chart1",
			expectedTitle: "Updated Sales",
		},
		{
			name: "Insight asset partial update",
			jsonData: `{
				"type": "insight",
				"id": "insight1",
				"insight_text": "Original insight"
			}`,
			partialUpdate: `{
				"insight_text": "Updated insight"
			}`,
			expectedID: "insight1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asset, err := CreateAsset([]byte(tt.jsonData))
			assert.NoError(t, err)

			created := asset.GetCreatedTimestamp()
			updated := asset.GetUpdatedTimestamp()

			time.Sleep(5 * time.Millisecond)

			err = json.Unmarshal([]byte(tt.partialUpdate), asset)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedID, asset.GetID(), "ID should remain unchanged")
			assert.Equal(t, created, asset.GetCreatedTimestamp(), "CreatedAt should not change")

			asset.SetUpdatedTimestamp(time.Now())
			assert.True(t, asset.GetUpdatedTimestamp().After(updated), "UpdatedAt should be refreshed")

			if chart, ok := asset.(*Chart); ok {
				assert.Equal(t, tt.expectedTitle, chart.Title, "Chart title mismatch")
			}
		})
	}
}

// TestCreateAssetRegistry verifies correct type creation
func TestCreateAssetRegistry(t *testing.T) {
	tests := []struct {
		typeStr   string
		expectVal Asset
		expectErr bool
	}{
		{"chart", &Chart{}, false},
		{"insight", &Insight{}, false},
		{"audience", &Audience{}, false},
		{"unknown", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.typeStr, func(t *testing.T) {
			data := []byte(`{"type":"` + tt.typeStr + `"}`)
			asset, err := CreateAsset(data)
			if tt.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.IsType(t, tt.expectVal, asset)
		})
	}
}

// TestCreateAssetInsufficientJSON ensures missing fields are handled gracefully
func TestCreateAssetInsufficientJSON(t *testing.T) {
	data := []byte(`{"type":"insight"}`)
	asset, err := CreateAsset(data)
	assert.NoError(t, err)
	assert.Equal(t, "", asset.GetDescription())
}
