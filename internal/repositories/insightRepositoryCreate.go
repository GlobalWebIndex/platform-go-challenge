package repositories

import (
	"context"
	"github.com/loukaspe/platform-go-challenge/internal/core/domain"
	"gorm.io/gorm"
)

type InsightRepository struct {
	db *gorm.DB
}

func NewInsightRepository(db *gorm.DB) *InsightRepository {
	return &InsightRepository{db: db}
}

func (repo *InsightRepository) CreateInsight(
	ctx context.Context,
	insight *domain.Insight,
) (uint, error) {
	var err error

	modelInsight := Insight{
		Text:        insight.Text,
		Description: insight.Description,
	}

	err = repo.db.WithContext(ctx).Create(&modelInsight).Error
	if err != nil {
		return modelInsight.ID, err
	}

	return modelInsight.ID, nil
}
