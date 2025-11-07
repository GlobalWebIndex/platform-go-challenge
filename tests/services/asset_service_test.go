package services_test

import (
	"assetsApp/internal/models"
	assetServices "assetsApp/internal/services/asset"
	"assetsApp/tests/mocks"
	"testing"

	"github.com/google/uuid"
)

func TestAssetService_GetAssets(t *testing.T) {
	userID := uuid.New()
	mockStore := &mocks.MockAssetStore{
		GetFunc: func(uid uuid.UUID) []models.Asset {
			if uid == userID {
				return []models.Asset{&models.Chart{ID: "test-chart"}}
			}
			return nil
		},
	}

	service := assetServices.NewAssetService(mockStore)

	assets := service.GetAssets(userID)

	if len(assets) != 1 {
		t.Errorf("expected 1 asset, got %d", len(assets))
	}
}

func TestAssetService_AddAsset(t *testing.T) {
	userID := uuid.New()
	asset := &models.Chart{ID: "test-chart"}
	called := false
	mockStore := &mocks.MockAssetStore{
		AddFunc: func(uid uuid.UUID, a models.Asset) {
			if uid == userID && a.GetID() == asset.ID {
				called = true
			}
		},
	}

	service := assetServices.NewAssetService(mockStore)

	service.AddAsset(userID, asset)

	if !called {
		t.Error("AddAsset was not called on the store")
	}
}

func TestAssetService_EditDescription(t *testing.T) {
	userID := uuid.New()
	assetID := "test-asset"
	description := "New Description"
	mockStore := &mocks.MockAssetStore{
		EditDescriptionFunc: func(uid uuid.UUID, aid string, desc string) bool {
			return uid == userID && aid == assetID && desc == description
		},
	}

	service := assetServices.NewAssetService(mockStore)

	if !service.EditDescription(userID, assetID, description) {
		t.Error("EditDescription returned false")
	}
}

func TestAssetService_RemoveAsset(t *testing.T) {
	userID := uuid.New()
	assetID := "test-asset"
	mockStore := &mocks.MockAssetStore{
		RemoveFunc: func(uid uuid.UUID, aid string) bool {
			return uid == userID && aid == assetID
		},
	}

	service := assetServices.NewAssetService(mockStore)

	if !service.RemoveAsset(userID, assetID) {
		t.Error("RemoveAsset returned false")
	}
}

func TestAssetService_EditNonExistentAsset(t *testing.T) {
	userID := uuid.New()
	assetID := "non-existent-asset"
	description := "New Description"
	mockStore := &mocks.MockAssetStore{
		EditDescriptionFunc: func(uid uuid.UUID, aid string, desc string) bool {
			return false // Simulate asset not found
		},
	}

	service := assetServices.NewAssetService(mockStore)

	if service.EditDescription(userID, assetID, description) {
		t.Error("EditDescription returned true for non-existent asset")
	}
}

func TestAssetService_RemoveNonExistentAsset(t *testing.T) {
	userID := uuid.New()
	assetID := "non-existent-asset"
	mockStore := &mocks.MockAssetStore{
		RemoveFunc: func(uid uuid.UUID, aid string) bool {
			return false // Simulate asset not found
		},
	}

	service := assetServices.NewAssetService(mockStore)

	if service.RemoveAsset(userID, assetID) {
		t.Error("RemoveAsset returned true for non-existent asset")
	}
}
