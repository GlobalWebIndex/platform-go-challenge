package mapper

import (
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/entities"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
)

// InsightEntityFromDomain converts domain model to entity
func InsightEntityFromDomain(i *domain.Insight) *entities.InsightEntity {
	if i == nil {
		return nil
	}

	return &entities.InsightEntity{
		AssetBaseEntity: *AssetBaseEntityFromDomain(i.AssetBase),
		Text:            i.Text,
	}
}

// InsightEntityToDomain converts entity to domain model
func InsightEntityToDomain(i *entities.InsightEntity) *domain.Insight {
	if i == nil {
		return nil
	}

	return &domain.Insight{
		AssetBase: AssetBaseEntityToDomain(i.AssetBaseEntity),
		Text:      i.Text,
	}
}
