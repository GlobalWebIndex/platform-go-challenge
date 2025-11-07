package services_test

import (
	"assetsApp/internal/models"
	favouriteServices "assetsApp/internal/services/favourite"
	"assetsApp/tests/mocks"
	"testing"

	"github.com/google/uuid"
)

func TestFavouriteService_GetFavourites(t *testing.T) {
	userID := uuid.New()
	mockStore := &mocks.MockAssetStore{
		GetFavouritesFunc: func(uid uuid.UUID) []models.Favourite {
			if uid == userID {
				return []models.Favourite{{Asset: &models.Chart{ID: "test-chart"}}}
			}
			return nil
		},
	}

	service := favouriteServices.NewFavouriteService(mockStore)

	assets := service.GetFavourites(userID)

	if len(assets) != 1 {
		t.Errorf("expected 1 asset, got %d", len(assets))
	}
}

func TestFavouriteService_AddFavourite(t *testing.T) {
	userID := uuid.New()
	assetID := "test-asset"
	assetType := "chart"
	mockStore := &mocks.MockAssetStore{
		AddFavouriteFunc: func(uid uuid.UUID, aid string, atype string) bool {
			return uid == userID && aid == assetID && atype == assetType
		},
	}

	service := favouriteServices.NewFavouriteService(mockStore)

	if !service.AddFavourite(userID, assetID, assetType) {
		t.Error("AddFavourite returned false")
	}
}

func TestFavouriteService_RemoveFavourite(t *testing.T) {
	userID := uuid.New()
	assetID := "test-asset"
	mockStore := &mocks.MockAssetStore{
		RemoveFavouriteFunc: func(uid uuid.UUID, aid string) bool {
			return uid == userID && aid == assetID
		},
	}

	service := favouriteServices.NewFavouriteService(mockStore)

	if !service.RemoveFavourite(userID, assetID) {
		t.Error("RemoveFavourite returned false")
	}
}

func TestFavouriteService_AddFavouriteForNonExistentAsset(t *testing.T) {
	userID := uuid.New()
	assetID := "non-existent-asset"
	assetType := "chart"
	mockStore := &mocks.MockAssetStore{
		AddFavouriteFunc: func(uid uuid.UUID, aid string, atype string) bool {
			return false // Simulate asset not found
		},
	}

	service := favouriteServices.NewFavouriteService(mockStore)

	if service.AddFavourite(userID, assetID, assetType) {
		t.Error("AddFavourite returned true for non-existent asset")
	}
}

func TestFavouriteService_RemoveNonExistentFavourite(t *testing.T) {
	userID := uuid.New()
	assetID := "non-existent-asset"
	mockStore := &mocks.MockAssetStore{
		RemoveFavouriteFunc: func(uid uuid.UUID, aid string) bool {
			return false // Simulate favourite not found
		},
	}

	service := favouriteServices.NewFavouriteService(mockStore)

	if service.RemoveFavourite(userID, assetID) {
		t.Error("RemoveFavourite returned true for non-existent favourite")
	}
}
