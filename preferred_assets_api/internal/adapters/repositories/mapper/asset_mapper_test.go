package mapper_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/entities"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/mapper"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
)

func newAudienceEntity() *entities.AudienceEntity {
	return &entities.AudienceEntity{
		AssetBaseEntity: entities.AssetBaseEntity{
			ID:    "a1",
			Type:  entities.AssetTypeAudience,
			Title: "Audience Asset",
		},
		Gender: "female",
	}
}

func newChartEntity() *entities.ChartEntity {
	axesJSON, _ := json.Marshal([]string{"X", "Y"})
	dataJSON, _ := json.Marshal([][]float64{{1, 2}, {3, 4}})
	return &entities.ChartEntity{
		AssetBaseEntity: entities.AssetBaseEntity{
			ID:    "c1",
			Type:  entities.AssetTypeChart,
			Title: "Chart Asset",
		},
		AxesTitles: string(axesJSON),
		Data:       string(dataJSON),
	}
}

func newInsightEntity() *entities.InsightEntity {
	return &entities.InsightEntity{
		AssetBaseEntity: entities.AssetBaseEntity{
			ID:    "i1",
			Type:  entities.AssetTypeInsight,
			Title: "Insight Asset",
		},
	}
}

func TestAssetEntityToDomain_HappyPaths(t *testing.T) {
	tests := []struct {
		name     string
		entity   entities.AssetEntity
		wantType string
	}{
		{"Audience", newAudienceEntity(), "Audience"},
		{"Chart", newChartEntity(), "Chart"},
		{"Insight", newInsightEntity(), "Insight"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			domainAsset, err := mapper.AssetEntityToDomain(tt.entity)

			// Assert
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if domainAsset == nil {
				t.Fatalf("expected domain asset, got nil")
			}
		})
	}
}

func TestAssetEntityToDomain_UnhappyPaths(t *testing.T) {
	t.Run("nil entity", func(t *testing.T) {
		// Act
		domainAsset, err := mapper.AssetEntityToDomain(nil)

		// Assert
		if err != nil {
			t.Errorf("expected no error for nil entity, got %v", err)
		}
		if domainAsset != nil {
			t.Errorf("expected nil domain asset, got %+v", domainAsset)
		}
	})

	t.Run("unknown type", func(t *testing.T) {
		// Arrange
		e := &entities.AudienceEntity{}
		e.Type = 999 // unsupported type

		// Act
		_, err := mapper.AssetEntityToDomain(e)

		// Assert
		if err == nil {
			t.Errorf("expected error for unknown type, got nil")
		}
	})
}

func TestAssetEntityFromDomain_HappyPaths(t *testing.T) {
	aud := &domain.Audience{AssetBase: domain.AssetBase{ID: "a1", Type: domain.AssetTypeAudience, Title: "aud"}}
	chart := &domain.Chart{AssetBase: domain.AssetBase{ID: "c1", Type: domain.AssetTypeChart, Title: "chart"}}
	insight := &domain.Insight{AssetBase: domain.AssetBase{ID: "i1", Type: domain.AssetTypeInsight, Title: "insight"}}

	tests := []struct {
		name  string
		input domain.Asset
	}{
		{"Audience", aud},
		{"Chart", chart},
		{"Insight", insight},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			entity, err := mapper.AssetEntityFromDomain(tt.input)

			// Assert
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if entity == nil {
				t.Fatalf("expected entity, got nil")
			}
		})
	}
}

func TestAssetEntityFromDomain_UnhappyPaths(t *testing.T) {
	t.Run("nil domain asset", func(t *testing.T) {
		// Act
		entity, err := mapper.AssetEntityFromDomain(nil)

		// Assert
		if err != nil {
			t.Errorf("expected no error for nil domain, got %v", err)
		}
		if entity != nil {
			t.Errorf("expected nil entity, got %+v", entity)
		}
	})
}

func TestAssetEntityBaseMapping(t *testing.T) {
	now := time.Now()
	base := domain.AssetBase{
		ID:        "1",
		Type:      domain.AssetTypeAudience,
		Title:     "title",
		CreatedAt: now,
		UpdatedAt: now,
	}

	t.Run("Domain -> Entity", func(t *testing.T) {
		// Act
		ent := mapper.AssetBaseEntityFromDomain(base)

		// Assert
		if ent.ID != base.ID || ent.Title != base.Title {
			t.Errorf("mapping mismatch: got %+v", ent)
		}
	})

	t.Run("Entity -> Domain", func(t *testing.T) {
		ent := entities.AssetBaseEntity{
			ID:        "2",
			Type:      entities.AssetTypeChart,
			Title:     "chart",
			CreatedAt: now,
			UpdatedAt: now,
		}
		dom := mapper.AssetBaseEntityToDomain(ent)
		if dom.ID != ent.ID || dom.Title != ent.Title {
			t.Errorf("mapping mismatch: got %+v", dom)
		}
	})
}

func TestAssetEntityToDomain_TableDriven(t *testing.T) {

	aud := &entities.AudienceEntity{AssetBaseEntity: entities.AssetBaseEntity{ID: "a1", Type: entities.AssetTypeAudience, Title: "aud"}}
	chart := &entities.ChartEntity{AssetBaseEntity: entities.AssetBaseEntity{ID: "c1", Type: entities.AssetTypeChart, Title: "chart"}}
	insight := &entities.InsightEntity{AssetBaseEntity: entities.AssetBaseEntity{ID: "i1", Type: entities.AssetTypeInsight, Title: "insight"}}
	unknown := &entities.AudienceEntity{AssetBaseEntity: entities.AssetBaseEntity{ID: "u1", Type: 999, Title: "unknown"}}

	tests := []struct {
		name    string
		entity  entities.AssetEntity
		wantNil bool
		wantErr bool
	}{
		{"Audience", aud, false, false},
		{"Chart", chart, false, false},
		{"Insight", insight, false, false},
		{"Nil entity", nil, true, false},
		{"Unknown type", unknown, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			dom, err := mapper.AssetEntityToDomain(tt.entity)

			// Assert
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.wantNil {
				if dom != nil {
					t.Errorf("expected nil domain, got %+v", dom)
				}
			} else {
				if dom == nil {
					t.Errorf("expected non-nil domain")
				}
			}
		})
	}
}

func TestAssetEntityFromDomain_TableDriven(t *testing.T) {
	aud := &domain.Audience{AssetBase: domain.AssetBase{ID: "a1", Type: domain.AssetTypeAudience, Title: "aud"}}
	chart := &domain.Chart{AssetBase: domain.AssetBase{ID: "c1", Type: domain.AssetTypeChart, Title: "chart"}}
	insight := &domain.Insight{AssetBase: domain.AssetBase{ID: "i1", Type: domain.AssetTypeInsight, Title: "insight"}}
	var nilAsset domain.Asset

	tests := []struct {
		name    string
		domain  domain.Asset
		wantNil bool
	}{
		{"Audience", aud, false},
		{"Chart", chart, false},
		{"Insight", insight, false},
		{"Nil domain", nilAsset, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			ent, err := mapper.AssetEntityFromDomain(tt.domain)

			// Assert
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.wantNil {
				if ent != nil {
					t.Errorf("expected nil entity, got %+v", ent)
				}
			} else {
				if ent == nil {
					t.Errorf("expected non-nil entity")
				}
			}
		})
	}
}

func TestAssetBaseMapping(t *testing.T) {
	now := time.Now()
	base := domain.AssetBase{
		ID:        "1",
		Type:      domain.AssetTypeAudience,
		Title:     "title",
		CreatedAt: now,
		UpdatedAt: now,
	}

	t.Run("Domain -> Entity -> Domain", func(t *testing.T) {
		// Act
		ent := mapper.AssetBaseEntityFromDomain(base)
		dom := mapper.AssetBaseEntityToDomain(*ent)

		// Assert
		if dom.ID != base.ID || dom.Title != base.Title || dom.Type != base.Type {
			t.Errorf("roundtrip mismatch: got %+v", dom)
		}
	})
}

func TestAssetEntityToDomain_Batch(t *testing.T) {

	// Prepare a mix of entities
	aud := &entities.AudienceEntity{AssetBaseEntity: entities.AssetBaseEntity{ID: "a1", Type: entities.AssetTypeAudience, Title: "aud"}}
	chart := &entities.ChartEntity{AssetBaseEntity: entities.AssetBaseEntity{ID: "c1", Type: entities.AssetTypeChart, Title: "chart"}}
	insight := &entities.InsightEntity{AssetBaseEntity: entities.AssetBaseEntity{ID: "i1", Type: entities.AssetTypeInsight, Title: "insight"}}
	nilEntity := entities.AssetEntity(nil)
	unknown := &entities.AudienceEntity{AssetBaseEntity: entities.AssetBaseEntity{ID: "u1", Type: 999, Title: "unknown"}}

	entitiesSlice := []entities.AssetEntity{aud, chart, insight, nilEntity, unknown}
	expectedNil := []bool{false, false, false, true, true}
	expectedErr := []bool{false, false, false, false, true}

	for i, ent := range entitiesSlice {
		t.Run(fmt.Sprintf("Index %d", i), func(t *testing.T) {
			// Act
			dom, err := mapper.AssetEntityToDomain(ent)

			// Assert
			if expectedErr[i] {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if expectedNil[i] {
				if dom != nil {
					t.Errorf("expected nil domain, got %+v", dom)
				}
			} else {
				if dom == nil {
					t.Errorf("expected non-nil domain")
				}
			}
		})
	}
}

func TestAssetEntityFromDomain_Batch(t *testing.T) {
	now := time.Now()

	// Prepare a mix of domain assets
	aud := &domain.Audience{AssetBase: domain.AssetBase{ID: "a1", Type: domain.AssetTypeAudience, Title: "aud", CreatedAt: now}}
	chart := &domain.Chart{AssetBase: domain.AssetBase{ID: "c1", Type: domain.AssetTypeChart, Title: "chart", CreatedAt: now}}
	insight := &domain.Insight{AssetBase: domain.AssetBase{ID: "i1", Type: domain.AssetTypeInsight, Title: "insight", CreatedAt: now}}
	nilAsset := domain.Asset(nil)

	// Custom type to simulate unknown
	type UnknownAsset struct{ domain.AssetBase }
	unknown := &UnknownAsset{AssetBase: domain.AssetBase{ID: "u1", Type: 999, Title: "unknown", CreatedAt: now}}

	domainSlice := []domain.Asset{aud, chart, insight, nilAsset, unknown}
	expectedNil := []bool{false, false, false, true, true}
	expectedErr := []bool{false, false, false, false, true}

	for i, asset := range domainSlice {
		t.Run(fmt.Sprintf("Index %d", i), func(t *testing.T) {
			// Act
			entity, err := mapper.AssetEntityFromDomain(asset)

			// Assert
			if expectedErr[i] {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if expectedNil[i] {
				if entity != nil {
					t.Errorf("expected nil entity, got %+v", entity)
				}
			} else {
				if entity == nil {
					t.Errorf("expected non-nil entity")
				}
			}
		})
	}
}

func TestAssetBaseEntityMapping_Batch(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name       string
		domainBase domain.AssetBase
		entityBase entities.AssetBaseEntity
	}{
		{
			name: "normal domain to entity",
			domainBase: domain.AssetBase{
				ID:          "a1",
				Type:        domain.AssetTypeAudience,
				Title:       "Title1",
				Description: "Desc1",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
		{
			name:       "zero-value domain",
			domainBase: domain.AssetBase{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// --- Domain -> Entity ---
			entity := mapper.AssetBaseEntityFromDomain(tt.domainBase)

			if entity.ID != tt.domainBase.ID {
				t.Errorf("expected ID %v, got %v", tt.domainBase.ID, entity.ID)
			}
			if domain.AssetType(entity.Type) != tt.domainBase.Type {
				t.Errorf("expected Type %v, got %v", tt.domainBase.Type, entity.Type)
			}
			if entity.Title != tt.domainBase.Title {
				t.Errorf("expected Title %v, got %v", tt.domainBase.Title, entity.Title)
			}

			// --- Entity -> Domain ---
			domainBack := mapper.AssetBaseEntityToDomain(*entity)
			if domainBack.ID != tt.domainBase.ID {
				t.Errorf("expected ID %v, got %v", tt.domainBase.ID, domainBack.ID)
			}
			if domainBack.Type != tt.domainBase.Type {
				t.Errorf("expected Type %v, got %v", tt.domainBase.Type, domainBack.Type)
			}
			if domainBack.Title != tt.domainBase.Title {
				t.Errorf("expected Title %v, got %v", tt.domainBase.Title, domainBack.Title)
			}
		})
	}
}
