package mapper

import (
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/entities"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
)

// AudienceEntityToDomain converts entity to domain model
func AudienceEntityToDomain(a *entities.AudienceEntity) *domain.Audience {
	if a == nil {
		return nil
	}
	return &domain.Audience{
		AssetBase:       AssetBaseEntityToDomain(a.AssetBaseEntity),
		Gender:          a.Gender,
		BirthCountry:    a.BirthCountry,
		AgeGroup:        a.AgeGroup,
		HoursSocial:     a.HoursSocial,
		PurchasesLastMo: a.PurchasesLastMo,
	}
}

// AudienceEntityFromDomain converts domain model to entity
func AudienceEntityFromDomain(a *domain.Audience) *entities.AudienceEntity {
	if a == nil {
		return nil
	}
	return &entities.AudienceEntity{
		AssetBaseEntity: *AssetBaseEntityFromDomain(a.AssetBase),
		Gender:          a.Gender,
		BirthCountry:    a.BirthCountry,
		AgeGroup:        a.AgeGroup,
		HoursSocial:     a.HoursSocial,
		PurchasesLastMo: a.PurchasesLastMo,
	}
}
