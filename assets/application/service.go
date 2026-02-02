package application

import (
	"context"
	"errors"
	"log"

	"github.com/manos/favourites/assets/domain"
)

type AssetService struct {
	repo AssetRepository
}

func NewAssetService(repo AssetRepository) *AssetService {
	return &AssetService{repo: repo}
}

var ErrAssetNotFound = errors.New("asset not found")

// --- Facade ---

type AssetsService interface {
	GetInsight(ctx context.Context, userId string, assetId string) (domain.InsightAsset, error)
	GetAudience(ctx context.Context, userId string, assetId string) (domain.AudienceAsset, error)
	GetChart(ctx context.Context, userId string, assetId string) (domain.ChartAsset, error)
	GetInsights(ctx context.Context, userId string, assetIds []string) (map[string]domain.InsightAsset, error)
	GetAudiences(ctx context.Context, userId string, assetIds []string) (map[string]domain.AudienceAsset, error)
	GetCharts(ctx context.Context, userId string, assetIds []string) (map[string]domain.ChartAsset, error)
}

func (f *AssetService) GetInsight(ctx context.Context, userId string, assetId string) (domain.InsightAsset, error) {
	log.Printf("svc assets: get insight user=%s asset_id=%s", userId, assetId)
	a, ok := f.repo.GetInsight(ctx, userId, assetId)
	if !ok {
		log.Printf("svc assets: get insight user=%s asset_id=%s not found", userId, assetId)
		return domain.InsightAsset{}, ErrAssetNotFound
	}
	return a, nil
}

func (f *AssetService) GetAudience(ctx context.Context, userId string, assetId string) (domain.AudienceAsset, error) {
	log.Printf("svc assets: get audience user=%s asset_id=%s", userId, assetId)
	a, ok := f.repo.GetAudience(ctx, userId, assetId)
	if !ok {
		log.Printf("svc assets: get audience user=%s asset_id=%s not found", userId, assetId)
		return domain.AudienceAsset{}, ErrAssetNotFound
	}
	return a, nil
}

func (f *AssetService) GetChart(ctx context.Context, userId string, assetId string) (domain.ChartAsset, error) {
	log.Printf("svc assets: get chart user=%s asset_id=%s", userId, assetId)
	a, ok := f.repo.GetChart(ctx, userId, assetId)
	if !ok {
		log.Printf("svc assets: get chart user=%s asset_id=%s not found", userId, assetId)
		return domain.ChartAsset{}, ErrAssetNotFound
	}
	return a, nil
}

func (f *AssetService) GetInsights(ctx context.Context, userId string, assetIds []string) (map[string]domain.InsightAsset, error) {
	log.Printf("svc assets: get insights batch user=%s count=%d", userId, len(assetIds))
	assets := f.repo.GetInsights(ctx, userId, assetIds)
	if len(assets) != len(assetIds) {
		log.Printf("svc assets: get insights batch user=%s missing assets", userId)
		return nil, ErrAssetNotFound
	}
	return assets, nil
}

func (f *AssetService) GetAudiences(ctx context.Context, userId string, assetIds []string) (map[string]domain.AudienceAsset, error) {
	log.Printf("svc assets: get audiences batch user=%s count=%d", userId, len(assetIds))
	assets := f.repo.GetAudiences(ctx, userId, assetIds)
	if len(assets) != len(assetIds) {
		log.Printf("svc assets: get audiences batch user=%s missing assets", userId)
		return nil, ErrAssetNotFound
	}
	return assets, nil
}

func (f *AssetService) GetCharts(ctx context.Context, userId string, assetIds []string) (map[string]domain.ChartAsset, error) {
	log.Printf("svc assets: get charts batch user=%s count=%d", userId, len(assetIds))
	assets := f.repo.GetCharts(ctx, userId, assetIds)
	if len(assets) != len(assetIds) {
		log.Printf("svc assets: get charts batch user=%s missing assets", userId)
		return nil, ErrAssetNotFound
	}
	return assets, nil
}
