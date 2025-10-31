package mapper

import (
	"fmt"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/entities"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
)

func AssetEntityToDomain(entity entities.AssetEntity) (domain.Asset, error) {
	if entity == nil {
		return nil, nil
	}

	switch entity.GetType() {
	case entities.AssetTypeAudience:
		return AudienceEntityToDomain(entity.(*entities.AudienceEntity)), nil
	case entities.AssetTypeChart:
		return ChartEntityToDomain(entity.(*entities.ChartEntity)), nil
	case entities.AssetTypeInsight:
		return InsightEntityToDomain(entity.(*entities.InsightEntity)), nil
	default:
		return nil, fmt.Errorf("unknown asset type: %v", entity.GetType())
	}
}

func AssetEntityFromDomain(asset domain.Asset) (entities.AssetEntity, error) {
	if asset == nil {
		return nil, nil
	}
	switch asset.GetType() {
	case domain.AssetTypeAudience:
		return AudienceEntityFromDomain(asset.(*domain.Audience)), nil
	case domain.AssetTypeChart:
		return ChartEntityFromDomain(asset.(*domain.Chart)), nil
	case domain.AssetTypeInsight:
		return InsightEntityFromDomain(asset.(*domain.Insight)), nil
	default:
		return nil, fmt.Errorf("unknown asset type: %v", asset.GetType())
	}
}

// AssetBaseDomainToEntity domain to entity mapping
func AssetBaseEntityFromDomain(a domain.AssetBase) *entities.AssetBaseEntity {
	return &entities.AssetBaseEntity{
		ID:          a.ID,
		Type:        entities.AssetType(a.Type),
		Title:       a.Title,
		Description: a.Description,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

// AssetBaseEntityToDomain entity to domain mapping
func AssetBaseEntityToDomain(e entities.AssetBaseEntity) domain.AssetBase {
	return domain.AssetBase{
		ID:          e.ID,
		Type:        domain.AssetType(e.Type),
		Title:       e.Title,
		Description: e.Description,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
