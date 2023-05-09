package repositories

import (
	"context"
	"github.com/loukaspe/platform-go-challenge/internal/core/domain"
	"gorm.io/gorm"
)

type AudienceRepository struct {
	db *gorm.DB
}

func NewAudienceRepository(db *gorm.DB) *AudienceRepository {
	return &AudienceRepository{db: db}
}

func (repo *AudienceRepository) CreateAudience(
	ctx context.Context,
	audience *domain.Audience,
) (uint, error) {
	var err error

	modelAudience := Audience{
		Gender:             audience.Gender,
		BirthCountry:       audience.BirthCountry,
		AgeGroup:           audience.AgeGroup,
		HoursSpentOnSocial: audience.HoursSpentDaily,
		PurchasesLastMonth: audience.PurchasesLastMonth,
		Description:        audience.Description,
	}

	err = repo.db.WithContext(ctx).Create(&modelAudience).Error
	if err != nil {
		return modelAudience.ID, err
	}

	return modelAudience.ID, nil
}
