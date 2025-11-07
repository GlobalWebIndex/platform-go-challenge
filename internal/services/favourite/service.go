package favouriteServices

import (
	"assetsApp/internal/models"
	"assetsApp/internal/storage"

	"github.com/google/uuid"
)

type FavouriteService struct {
	store storage.AssetStore
}

func NewFavouriteService(store storage.AssetStore) *FavouriteService {
	return &FavouriteService{store: store}
}

func (s *FavouriteService) GetFavourites(userID uuid.UUID) []models.Favourite {
	return s.store.GetFavourites(userID)
}

func (s *FavouriteService) AddFavourite(userID uuid.UUID, assetID, assetType string) bool {
	return s.store.AddFavourite(userID, assetID, assetType)
}

func (s *FavouriteService) RemoveFavourite(userID uuid.UUID, assetID string) bool {
	return s.store.RemoveFavourite(userID, assetID)
}
