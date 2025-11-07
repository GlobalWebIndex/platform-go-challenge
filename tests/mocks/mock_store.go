
package mocks

import (
	"assetsApp/internal/models"
	"github.com/google/uuid"
)

// MockAssetStore is a mock implementation of the AssetStore interface.
type MockAssetStore struct {
	GetFunc             func(userID uuid.UUID) []models.Asset
	AddFunc             func(userID uuid.UUID, asset models.Asset)
	RemoveFunc          func(userID uuid.UUID, assetID string) bool
	EditDescriptionFunc func(userID uuid.UUID, assetID, newDesc string) bool

	GetFavouritesFunc   func(userID uuid.UUID) []models.Favourite
	AddFavouriteFunc    func(userID uuid.UUID, assetID, assetType string) bool
	RemoveFavouriteFunc func(userID uuid.UUID, assetID string) bool
}

func (m *MockAssetStore) Get(userID uuid.UUID) []models.Asset {
	if m.GetFunc != nil {
		return m.GetFunc(userID)
	}
	return nil
}

func (m *MockAssetStore) Add(userID uuid.UUID, asset models.Asset) {
	if m.AddFunc != nil {
		m.AddFunc(userID, asset)
	}
}

func (m *MockAssetStore) Remove(userID uuid.UUID, assetID string) bool {
	if m.RemoveFunc != nil {
		return m.RemoveFunc(userID, assetID)
	}
	return false
}

func (m *MockAssetStore) EditDescription(userID uuid.UUID, assetID, newDesc string) bool {
	if m.EditDescriptionFunc != nil {
		return m.EditDescriptionFunc(userID, assetID, newDesc)
	}
	return false
}

func (m *MockAssetStore) GetFavourites(userID uuid.UUID) []models.Favourite {
	if m.GetFavouritesFunc != nil {
		return m.GetFavouritesFunc(userID)
	}
	return nil
}

func (m *MockAssetStore) AddFavourite(userID uuid.UUID, assetID, assetType string) bool {
	if m.AddFavouriteFunc != nil {
		return m.AddFavouriteFunc(userID, assetID, assetType)
	}
	return false
}

func (m *MockAssetStore) RemoveFavourite(userID uuid.UUID, assetID string) bool {
	if m.RemoveFavouriteFunc != nil {
		return m.RemoveFavouriteFunc(userID, assetID)
	}
	return false
}
