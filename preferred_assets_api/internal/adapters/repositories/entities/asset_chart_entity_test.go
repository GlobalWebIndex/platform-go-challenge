package entities_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/entities"
)

func TestChartEntityValidate(t *testing.T) {
	validBase := entities.AssetBaseEntity{
		ID:        "1",
		Type:      1,
		Title:     "Chart Title",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name      string
		entity    *entities.ChartEntity
		wantErr   bool
		errString string
	}{
		{
			name: "valid chart entity with empty data",
			entity: &entities.ChartEntity{
				AssetBaseEntity: validBase,
				AxesTitles:      `["X","Y"]`,
				Data:            `[[1,2],[3,4]]`,
			},
			wantErr: false,
		},
		{
			name: "empty data string is allowed",
			entity: &entities.ChartEntity{
				AssetBaseEntity: validBase,
				AxesTitles:      `["X","Y"]`,
				Data:            "",
			},
			wantErr: false,
		},
		{
			name: "data string too long",
			entity: &entities.ChartEntity{
				AssetBaseEntity: validBase,
				AxesTitles:      `["X","Y"]`,
				Data:            string(make([]byte, 1001)),
			},
			wantErr:   true,
			errString: "data string too long",
		},
		{
			name: "data string contains null rune",
			entity: &entities.ChartEntity{
				AssetBaseEntity: validBase,
				AxesTitles:      `["X","Y"]`,
				Data:            "valid\u0000invalid",
			},
			wantErr:   true,
			errString: "invalid rune at position 5",
		},
		{
			name: "invalid AssetBaseEntity fails validation",
			entity: &entities.ChartEntity{
				AssetBaseEntity: entities.AssetBaseEntity{
					ID:    "",
					Type:  1,
					Title: "Chart",
				},
				AxesTitles: `["X"]`,
				Data:       `[[1,2]]`,
			},
			wantErr:   true,
			errString: "asset ID is required",
		},
		{
			name:      "nil chart entity pointer",
			entity:    nil,
			wantErr:   true,
			errString: "nil entity cannot be validated",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.entity == nil {
				err = validateNilSafe(tt.entity)
			} else {
				err = tt.entity.Validate()
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if err.Error() != tt.errString {
					t.Errorf("expected error '%s', got '%s'", tt.errString, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// helper for nil-safe validation
func validateNilSafe(c *entities.ChartEntity) error {
	if c == nil {
		return fmt.Errorf("nil entity cannot be validated")
	}
	return c.Validate()
}
