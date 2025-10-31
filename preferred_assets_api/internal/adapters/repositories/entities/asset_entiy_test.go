package entities_test

import (
	"testing"
	"time"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/entities"
)

func TestAssetBaseEntityValidate(t *testing.T) {
	tests := []struct {
		name      string
		entity    entities.AssetBaseEntity
		wantErr   bool
		errString string
	}{
		{
			name: "valid entity",
			entity: entities.AssetBaseEntity{
				ID:        "1",
				Type:      1,
				Title:     "Asset Title",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			entity: entities.AssetBaseEntity{
				ID:        "",
				Type:      1,
				Title:     "Asset Title",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr:   true,
			errString: "asset ID is required",
		},
		{
			name: "missing Title",
			entity: entities.AssetBaseEntity{
				ID:        "1",
				Type:      1,
				Title:     "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr:   true,
			errString: "asset title is required",
		},
		{
			name: "missing Type",
			entity: entities.AssetBaseEntity{
				ID:        "1",
				Type:      0,
				Title:     "Asset Title",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr:   true,
			errString: "asset type is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := tt.entity.Validate()

			// Assert
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
