package entities_test

import (
	"strings"
	"testing"
	"time"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/entities"
)

func TestInsightEntityValidate(t *testing.T) {
	tests := []struct {
		name      string
		entity    entities.InsightEntity
		wantError bool
	}{
		{
			name: "valid insight",
			entity: entities.InsightEntity{
				AssetBaseEntity: entities.AssetBaseEntity{
					ID:          "1",
					Type:        1,
					Title:       "Test Insight",
					Description: "Valid description",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				Text: "This is valid insight text.",
			},
			wantError: false,
		},
		{
			name: "missing ID",
			entity: entities.InsightEntity{
				AssetBaseEntity: entities.AssetBaseEntity{
					ID:          "",
					Type:        1,
					Title:       "Test Insight",
					Description: "Valid description",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				Text: "Some text",
			},
			wantError: true,
		},
		{
			name: "missing title",
			entity: entities.InsightEntity{
				AssetBaseEntity: entities.AssetBaseEntity{
					ID:          "1",
					Type:        1,
					Title:       "",
					Description: "Valid description",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				Text: "Some text",
			},
			wantError: true,
		},
		{
			name: "missing type",
			entity: entities.InsightEntity{
				AssetBaseEntity: entities.AssetBaseEntity{
					ID:          "1",
					Type:        0,
					Title:       "Test Insight",
					Description: "Valid description",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				Text: "Some text",
			},
			wantError: true,
		},
		{
			name: "text too long",
			entity: entities.InsightEntity{
				AssetBaseEntity: entities.AssetBaseEntity{
					ID:          "1",
					Type:        1,
					Title:       "Test Insight",
					Description: "Valid description",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				Text: strings.Repeat("a", 1001),
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := tt.entity.Validate()

			// Assert
			if tt.wantError && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
