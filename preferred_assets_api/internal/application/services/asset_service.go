package services

import (
	"time"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/mapper"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/ports"
)

var _ ports.AssetService = (*AssetServiceImpl)(nil)

type AssetServiceImpl struct {
	assetRepo ports.AssetRepository
}

func NewAssetService(assetRepo ports.AssetRepository) *AssetServiceImpl {
	return &AssetServiceImpl{
		assetRepo: assetRepo}
}

// CreateAsset implements ports.AssetService.
func (assetService *AssetServiceImpl) CreateAsset(asset domain.Asset) (domain.Asset, error) {
	asset.SetCreatedAt(time.Now().UTC())
	assetEntity, err := mapper.AssetEntityFromDomain(asset)
	if err != nil {
		return nil, err
	}

	createdAsset, err := assetService.assetRepo.Save(assetEntity)
	if err != nil {
		return nil, err
	}

	createdAssetDomain, err := mapper.AssetEntityToDomain(createdAsset)
	if err != nil {
		return nil, err
	}
	return createdAssetDomain, nil
}

// DeleteAsset implements ports.AssetService.
func (assetService *AssetServiceImpl) DeleteAsset(id string) error {
	return assetService.assetRepo.Delete(id)
}
