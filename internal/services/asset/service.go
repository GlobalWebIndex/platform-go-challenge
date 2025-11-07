package assetServices

import (
	"assetsApp/internal/models"
	"assetsApp/internal/storage"

	"github.com/google/uuid"
)

type AssetService struct {
	store storage.AssetStore
}

func NewAssetService(store storage.AssetStore) *AssetService {
	return &AssetService{store: store}
}

func (s *AssetService) GetAssets(userID uuid.UUID) []models.Asset {
	return s.store.Get(userID)
}

func (s *AssetService) AddAsset(userID uuid.UUID, asset models.Asset) {
	s.store.Add(userID, asset)
}

func (s *AssetService) RemoveAsset(userID uuid.UUID, assetID string) bool {
	return s.store.Remove(userID, assetID)
}

func (s *AssetService) EditDescription(userID uuid.UUID, assetID, description string) bool {
	return s.store.EditDescription(userID, assetID, description)
}
