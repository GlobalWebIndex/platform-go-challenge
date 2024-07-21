package database

import (
	"context"
	"favourites/models"
)

type AssetService interface {
	GetAll(ctx context.Context) (models.AssetCollection, error)
	GetByIdAndType(ctx context.Context, assetId string, assetType string) (models.AssetCollection, error)
}

type assetService struct {
	C ChartService
	I InsightService
	A AudienceService
}

var _ AssetService = (*assetService)(nil)

func NewAssetService(charts ChartService, insights InsightService, audiences AudienceService) AssetService {

	return &assetService{C: charts, I: insights, A: audiences}
}

func (s *assetService) GetAll(ctx context.Context) (models.AssetCollection, error) {

	assets := new(models.AssetCollection)

	charts, err := s.C.GetAll(nil)
	if err != nil {
		charts = make([]models.Chart, 0)
	}

	insights, err := s.I.GetAll(nil)
	if err != nil {
		insights = make([]models.Insight, 0)
	}
	audiences, err := s.A.GetAll(nil)
	if err != nil {
		audiences = make([]models.Audience, 0)
	}

	assets.Charts = charts
	assets.Insights = insights
	assets.Audiences = audiences

	return *assets, err

}

func (s *assetService) GetByIdAndType(ctx context.Context, assetId string, assetType string) (models.AssetCollection, error) {

	assets := new(models.AssetCollection)
	var e error
	switch assetType {
	case models.CHART_ASSET:
		chart, err := s.C.GetByID(ctx, assetId)
		if err != nil {
			e = err
		}
		assets.Charts = []models.Chart{chart}
	case models.INSIGHT_ASSET:
		insight, err := s.I.GetByID(ctx, assetId)
		if err != nil {
			e = err
		}
		assets.Insights = []models.Insight{insight}
	case models.AUDIENCE_ASSET:
		audience, err := s.A.GetByID(ctx, assetId)
		if err != nil {
			e = err
		}
		assets.Audiences = []models.Audience{audience}
	}

	return *assets, e
}
