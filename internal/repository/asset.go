package repository

import (
	"gwi_api/internal/domain"
	"gwi_api/internal/dto"
)

type AssetQuery interface {
	AddAsset(asset dto.AssetDto) (*uint64, error)
	GetAssets(params domain.PaginationParams) ([]domain.Asset, error)
	DeleteAsset(assetId uint64) error
}

type assetQuery struct{}

// AddAsset implements AssetService.
func (a *assetQuery) AddAsset(asset dto.AssetDto) (*uint64, error) {
	return DB.AddAsset(asset)
}

// DeleteAsset implements AssetService.
func (a *assetQuery) DeleteAsset(assetId uint64) error {
	return DB.DeleteAsset(assetId)
}

// GetAsset implements AssetService.
func (a *assetQuery) GetAssets(params domain.PaginationParams) ([]domain.Asset, error) {
	return DB.GetAssets(params)
}

func NewAssetService() AssetQuery {
	return &assetQuery{}
}
