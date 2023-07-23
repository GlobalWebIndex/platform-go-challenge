package service

import (
	"gwi_api/internal/domain"
	"gwi_api/internal/dto"
	"gwi_api/internal/repository"
)

type AssetService interface {
	AddAsset(asset dto.AssetDto) (*uint64, error)
	GetAsset(params domain.PaginationParams) ([]domain.Asset, error)
	DeleteAsset(assetId uint64) error
}

type assetService struct {
	dbHandler repository.DBHandler
}

// AddAsset implements AssetService.
func (a *assetService) AddAsset(asset dto.AssetDto) (*uint64, error) {
	panic("unimplemented")
}

// DeleteAsset implements AssetService.
func (a *assetService) DeleteAsset(assetId uint64) error {
	panic("unimplemented")
}

// GetAsset implements AssetService.
func (a *assetService) GetAsset(params domain.PaginationParams) ([]domain.Asset, error) {
	panic("unimplemented")
}

func NewAssetService(dbHandler repository.DBHandler) AssetService {
	return &assetService{dbHandler}
}
