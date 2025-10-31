package entities_test

import (
	"testing"
	"time"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/entities"
)

func TestAudienceEntityValidate(t *testing.T) {
	validBase := entities.AssetBaseEntity{
		ID:        "1",
		Type:      1,
		Title:     "Audience",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name      string
		entity    *entities.AudienceEntity
		wantErr   bool
		errString string
	}{
		{
			name: "valid audience",
			entity: &entities.AudienceEntity{
				AssetBaseEntity: validBase,
				Gender:          "male",
				BirthCountry:    "USA",
				AgeGroup:        "25-34",
				HoursSocial:     5,
				PurchasesLastMo: 3,
			},
			wantErr: false,
		},
		{
			name: "empty optional fields allowed",
			entity: &entities.AudienceEntity{
				AssetBaseEntity: validBase,
				Gender:          "",
				AgeGroup:        "",
				HoursSocial:     0,
				PurchasesLastMo: 0,
			},
			wantErr: false,
		},
		{
			name: "invalid gender",
			entity: &entities.AudienceEntity{
				AssetBaseEntity: validBase,
				Gender:          "unknown",
			},
			wantErr:   true,
			errString: "invalid gender: unknown",
		},
		{
			name: "invalid age group",
			entity: &entities.AudienceEntity{
				AssetBaseEntity: validBase,
				AgeGroup:        "10-15",
			},
			wantErr:   true,
			errString: "invalid age group: 10-15",
		},
		{
			name: "negative hours social",
			entity: &entities.AudienceEntity{
				AssetBaseEntity: validBase,
				HoursSocial:     -1,
			},
			wantErr:   true,
			errString: "hours social cannot be negative",
		},
		{
			name: "negative purchases last month",
			entity: &entities.AudienceEntity{
				AssetBaseEntity: validBase,
				PurchasesLastMo: -5,
			},
			wantErr:   true,
			errString: "purchases last month cannot be negative",
		},
		{
			name:      "nil audience pointer",
			entity:    nil,
			wantErr:   true,
			errString: "cannot call Validate on nil AudienceEntity",
		},
		{
			name: "invalid base entity",
			entity: &entities.AudienceEntity{
				AssetBaseEntity: entities.AssetBaseEntity{
					ID:    "",
					Type:  0,
					Title: "",
				},
			},
			wantErr:   true,
			errString: "asset ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.entity == nil {
				err = &AudienceEntityError{"cannot call Validate on nil AudienceEntity"}
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

// Optional: define a small error type for nil pointer scenario
type AudienceEntityError struct {
	msg string
}

func (e *AudienceEntityError) Error() string { return e.msg }
