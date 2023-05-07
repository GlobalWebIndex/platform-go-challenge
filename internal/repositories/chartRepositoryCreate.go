package repositories

import (
	"context"
	"github.com/loukaspe/platform-go-challenge/internal/core/domain"
	"gorm.io/gorm"
)

type ChartRepository struct {
	db *gorm.DB
}

func NewChartRepository(db *gorm.DB) *ChartRepository {
	return &ChartRepository{db: db}
}

func (repo *ChartRepository) CreateChart(
	ctx context.Context,
	chart *domain.Chart,
) (uint, error) {
	var err error

	modelChart := Chart{
		Title:       chart.Title,
		XTitle:      chart.XAxisTitle,
		YTitle:      chart.YAxisTitle,
		Data:        chart.Data.(string),
		Description: chart.Description,
	}

	err = repo.db.WithContext(ctx).Create(&modelChart).Error
	if err != nil {
		return modelChart.ID, err
	}

	return modelChart.ID, nil
}
