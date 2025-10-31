package domain_test

import (
	"testing"
	"time"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"

	"github.com/stretchr/testify/require"
)

// Helper for a valid AssetBase
func newValidAssetBase() domain.AssetBase {
	now := time.Now().UTC()
	return domain.AssetBase{
		ID:          "123",
		Type:        domain.AssetTypeChart,
		Title:       "My Chart",
		Description: "Description",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func TestParseAssetType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected domain.AssetType
		wantErr  bool
	}{
		{"audience type", "audience", domain.AssetTypeAudience, false},
		{"chart type", "chart", domain.AssetTypeChart, false},
		{"insight type", "insight", domain.AssetTypeInsight, false},
		{"invalid type", "unknown", -10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			got, err := domain.ParseAssetType(tt.input)

			// Assert
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestAssetBase_Validate(t *testing.T) {
	tests := []struct {
		name      string
		asset     domain.AssetBase
		wantError bool
	}{
		{"empty title", domain.AssetBase{Title: ""}, true},
		{"valid asset", newValidAssetBase(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.asset.Validate()
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAssetBase_Getters(t *testing.T) {
	// Arrange
	a := newValidAssetBase()

	// Act & Assert
	t.Run("GetID", func(t *testing.T) {
		require.Equal(t, "123", a.GetID())
	})
	t.Run("GetType", func(t *testing.T) {
		require.Equal(t, domain.AssetTypeChart, a.GetType())
	})
	t.Run("GetTitle", func(t *testing.T) {
		require.Equal(t, "My Chart", a.GetTitle())
	})
	t.Run("GetDescription", func(t *testing.T) {
		require.Equal(t, "Description", a.GetDescription())
	})
	t.Run("GetCreatedAt", func(t *testing.T) {
		require.WithinDuration(t, time.Now().UTC(), a.GetCreatedAt(), time.Second)
	})
	t.Run("GetUpdatedAt", func(t *testing.T) {
		require.WithinDuration(t, time.Now().UTC(), a.GetUpdatedAt(), time.Second)
	})
}
