package services

import (
	"errors"
	"fmt"

	"github.com/alexandros-georgantas/platform-go-challenge/internal/helpers"
	"github.com/alexandros-georgantas/platform-go-challenge/internal/models"
	"github.com/alexandros-georgantas/platform-go-challenge/internal/serializers"
	"gorm.io/gorm"
)

type AssetService interface {
	GetAsset(aid uint) (*models.AssetResponse, error)
	GetAssets(p int, pS int, t string) (serializers.AssetsResponse, error)
	GetCharts(p int, pS int) (serializers.ChartsResponse, error)
	GetAudiences(p int, pS int) (serializers.AudienceResponse, error)
	GetInsights(p int, pS int) (serializers.InsightResponse, error)
	UpdateDescription(aId uint, d string) (*models.AssetResponse, error)
}

type assetService struct {
	db gorm.DB
}

func NewAssetService(db gorm.DB) (AssetService, error) {
	return &assetService{db: db}, nil
}

func (as *assetService) GetAssets(p int, pS int, t string) (serializers.AssetsResponse, error) {
	var assets []models.Asset
	var assetResponses []models.AssetResponse
	pr := serializers.AssetsResponse{}
	var count int64
	pr.Assets = nil
	pr.TotalCount = 0

	if t == "" {
		as.db.Model(&models.Asset{}).Count(&count)
		if err := as.db.Order("id asc").Order("updated_at desc").Scopes(helpers.Paginate(p, pS)).Find(&assets).Error; err != nil {
			return pr, errors.New("something went wrong while fetching assets")
		}
	} else {
		as.db.Where("related_type = ?", t).Model(&models.Asset{}).Count(&count)
		if err := as.db.Where("related_type = ?", t).Order("id asc").Order("updated_at desc").Scopes(helpers.Paginate(p, pS)).Find(&assets).Error; err != nil {
			return pr, errors.New("something went wrong while fetching assets")
		}
	}

	for _, asset := range assets {
		var assetResponse models.AssetResponse
		result, err := helpers.AggregateAssetType(asset, as.db)

		if err != nil {
			return pr, err
		}

		assetResponse = *result

		assetResponses = append(assetResponses, assetResponse)
	}
	pr.Assets = &assetResponses
	pr.TotalCount = int(count)
	return pr, nil
}

func (as *assetService) GetCharts(p int, pS int) (serializers.ChartsResponse, error) {
	var charts []models.Chart
	pr := serializers.ChartsResponse{}
	var count int64

	pr.Charts = nil
	pr.TotalCount = 0

	as.db.Model(&models.Chart{}).Count(&count)

	if err := as.db.Preload("Asset").Scopes(helpers.Paginate(p, pS)).Find(&charts).Error; err != nil {
		return pr, errors.New("something went wrong while fetching charts")

	}
	pr.Charts = &charts
	pr.TotalCount = int(count)
	return pr, nil
}

func (as *assetService) GetAudiences(p int, pS int) (serializers.AudienceResponse, error) {
	var audiences []models.Audience
	var count int64
	pr := serializers.AudienceResponse{}

	as.db.Model(&models.Audience{}).Count(&count)

	if err := as.db.Preload("Asset").Scopes(helpers.Paginate(p, pS)).Find(&audiences).Error; err != nil {
		return pr, errors.New("something went wrong while fetching audiences")

	}

	pr.Audiences = &audiences
	pr.TotalCount = int(count)

	return pr, nil
}

func (as *assetService) GetInsights(p int, pS int) (serializers.InsightResponse, error) {
	var insights []models.Insight
	var count int64
	pr := serializers.InsightResponse{}

	as.db.Model(&models.Insight{}).Count(&count)

	if err := as.db.Preload("Asset").Scopes(helpers.Paginate(p, pS)).Find(&insights).Error; err != nil {
		return pr, errors.New("something went wrong while fetching insights")

	}

	pr.Insights = &insights
	pr.TotalCount = int(count)

	return pr, nil
}

func (as *assetService) UpdateDescription(aId uint, d string) (*models.AssetResponse, error) {
	var asset models.Asset

	if fErr := as.db.First(&asset, aId).Error; fErr != nil {
		return nil, fmt.Errorf("something went wrong while fetching  asset with id %v", aId)
	}

	if pErr := as.db.Model(&asset).Update("description", d).Error; pErr != nil {
		return nil, fmt.Errorf("something went wrong while updating description of asset with id %v", aId)
	}

	result, aErr := helpers.AggregateAssetType(asset, as.db)

	// TODO: Here invalidation of user's favorites cache should take place

	if aErr != nil {
		return nil, aErr
	}

	return result, nil
}

func (as *assetService) GetAsset(aid uint) (*models.AssetResponse, error) {
	var asset models.Asset
	if err := as.db.First(&asset, aid).Error; err != nil {
		return nil, fmt.Errorf("something went wrong while fetching asset with id %v", aid)
	}
	result, aErr := helpers.AggregateAssetType(asset, as.db)

	if aErr != nil {
		return nil, aErr
	}

	return result, nil
}
