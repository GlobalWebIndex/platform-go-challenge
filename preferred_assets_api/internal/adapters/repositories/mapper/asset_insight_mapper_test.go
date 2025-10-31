package mapper_test

import (
	"testing"
	"time"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/entities"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/mapper"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
)

func TestInsightEntityMapping(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name      string
		domainObj *domain.Insight
		expectNil bool
	}{
		{
			name: "happy path",
			domainObj: &domain.Insight{
				AssetBase: domain.AssetBase{
					ID:          "i1",
					Type:        domain.AssetTypeInsight,
					Title:       "Insight Title",
					Description: "Insight Desc",
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				Text: "Some insight text",
			},
		},
		{
			name:      "nil domain input",
			domainObj: nil,
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// --- Domain -> Entity ---
			var entity *entities.InsightEntity
			if tt.domainObj != nil {
				entity = mapper.InsightEntityFromDomain(tt.domainObj)
			} else {
				entity = mapper.InsightEntityFromDomain(nil)
			}

			if tt.expectNil {
				if entity != nil {
					t.Errorf("expected nil entity, got %+v", entity)
				}
				return
			}

			if entity.ID != tt.domainObj.ID {
				t.Errorf("expected ID %v, got %v", tt.domainObj.ID, entity.ID)
			}
			if entity.Text != tt.domainObj.Text {
				t.Errorf("expected Text %v, got %v", tt.domainObj.Text, entity.Text)
			}

			// --- Entity -> Domain ---
			domainBack := mapper.InsightEntityToDomain(entity)
			if domainBack.ID != tt.domainObj.ID {
				t.Errorf("expected ID %v, got %v", tt.domainObj.ID, domainBack.ID)
			}
			if domainBack.Text != tt.domainObj.Text {
				t.Errorf("expected Text %v, got %v", tt.domainObj.Text, domainBack.Text)
			}
		})
	}
}
