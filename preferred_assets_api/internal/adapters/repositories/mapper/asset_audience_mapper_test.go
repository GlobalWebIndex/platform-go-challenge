package mapper_test

import (
	"testing"
	"time"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/entities"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/mapper"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
)

func TestAudienceEntityMapping(t *testing.T) {
	created := time.Now()
	updated := created.Add(time.Hour)

	t.Run("happy_path: domain to entity and back", func(t *testing.T) {
		// Arrange
		domainObj := &domain.Audience{
			AssetBase: domain.AssetBase{
				ID:          "aud1",
				Type:        domain.AssetTypeAudience,
				Title:       "Audience Title",
				Description: "Audience Description",
				CreatedAt:   created,
				UpdatedAt:   updated,
			},
			Gender:          "female",
			BirthCountry:    "Canada",
			AgeGroup:        "25-34",
			HoursSocial:     4,
			PurchasesLastMo: 12,
		}

		// Act
		entityObj := mapper.AudienceEntityFromDomain(domainObj)
		backToDomain := mapper.AudienceEntityToDomain(entityObj)

		// Assert
		if entityObj.ID != domainObj.ID || backToDomain.ID != domainObj.ID {
			t.Errorf("IDs do not match")
		}
		if backToDomain.Gender != domainObj.Gender || backToDomain.BirthCountry != domainObj.BirthCountry {
			t.Errorf("Fields did not map correctly")
		}
		if backToDomain.HoursSocial != domainObj.HoursSocial || backToDomain.PurchasesLastMo != domainObj.PurchasesLastMo {
			t.Errorf("Numeric fields did not map correctly")
		}
		if !backToDomain.CreatedAt.Equal(domainObj.CreatedAt) || !backToDomain.UpdatedAt.Equal(domainObj.UpdatedAt) {
			t.Errorf("Timestamps did not map correctly")
		}
	})

	t.Run("unhappy_path: nil domain input", func(t *testing.T) {
		// Arrange
		var domainObj *domain.Audience = nil

		// Act
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Recovered from panic as expected: %v", r)
			}
		}()
		entityObj := mapper.AudienceEntityFromDomain(domainObj)

		// Assert
		if entityObj != nil {
			t.Errorf("expected nil entity, got %+v", entityObj)
		}
	})

	t.Run("unhappy_path: nil entity input", func(t *testing.T) {
		// Arrange
		var entityObj *entities.AudienceEntity = nil

		// Act
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Recovered from panic as expected: %v", r)
			}
		}()
		domainObj := mapper.AudienceEntityToDomain(entityObj)

		// Assert
		if domainObj != nil {
			t.Errorf("expected nil domain, got %+v", domainObj)
		}
	})
}

func TestAudienceEntityMapping_List(t *testing.T) {
	created := time.Now()
	updated := created.Add(time.Hour)

	t.Run("happy_path: slice domain to entity and back", func(t *testing.T) {
		// Arrange
		domainSlice := []domain.Audience{
			{
				AssetBase: domain.AssetBase{
					ID:          "aud1",
					Type:        domain.AssetTypeAudience,
					Title:       "Audience 1",
					Description: "Desc 1",
					CreatedAt:   created,
					UpdatedAt:   updated,
				},
				Gender:          "female",
				BirthCountry:    "Canada",
				AgeGroup:        "25-34",
				HoursSocial:     4,
				PurchasesLastMo: 12,
			},
			{
				AssetBase: domain.AssetBase{
					ID:          "aud2",
					Type:        domain.AssetTypeAudience,
					Title:       "Audience 2",
					Description: "Desc 2",
					CreatedAt:   created,
					UpdatedAt:   updated,
				},
				Gender:          "male",
				BirthCountry:    "USA",
				AgeGroup:        "35-44",
				HoursSocial:     2,
				PurchasesLastMo: 5,
			},
		}

		// Act
		entitySlice := make([]*entities.AudienceEntity, len(domainSlice))
		for i, d := range domainSlice {
			entitySlice[i] = mapper.AudienceEntityFromDomain(&d)
		}

		backToDomain := make([]*domain.Audience, len(entitySlice))
		for i, e := range entitySlice {
			backToDomain[i] = mapper.AudienceEntityToDomain(e)
		}

		// Assert
		for i := range domainSlice {
			if backToDomain[i].ID != domainSlice[i].ID || backToDomain[i].Gender != domainSlice[i].Gender {
				t.Errorf("Mismatch in element %d: expected %+v got %+v", i, domainSlice[i], backToDomain[i])
			}
		}
	})

	t.Run("unhappy_path: empty slice", func(t *testing.T) {
		// Arrange
		var domainSlice []domain.Audience

		// Act
		entitySlice := make([]*entities.AudienceEntity, len(domainSlice))
		for i, d := range domainSlice {
			entitySlice[i] = mapper.AudienceEntityFromDomain(&d)
		}

		// Assert
		if len(entitySlice) != 0 {
			t.Errorf("expected empty slice, got %d elements", len(entitySlice))
		}
	})

	t.Run("unhappy_path: slice with nil elements", func(t *testing.T) {
		// Arrange
		var domainSlice []*domain.Audience
		domainSlice = append(domainSlice, nil)

		// Act
		entitySlice := make([]*entities.AudienceEntity, len(domainSlice))
		for i, d := range domainSlice {
			entitySlice[i] = mapper.AudienceEntityFromDomain(d)
		}

		// Assert
		if entitySlice[0] != nil {
			t.Errorf("expected nil entity, got %+v", entitySlice[0])
		}
	})
}
