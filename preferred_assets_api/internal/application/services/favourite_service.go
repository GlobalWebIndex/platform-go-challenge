package services

import (
	"fmt"
	"time"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/mapper"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/ports"
)

var _ ports.FavouriteService = (*FavouriteServiceImpl)(nil)

type FavouriteServiceImpl struct {
	repo ports.FavouriteRepository
}

func NewFavouriteService(r ports.FavouriteRepository) *FavouriteServiceImpl {
	return &FavouriteServiceImpl{repo: r}
}

func (s FavouriteServiceImpl) CreateFavourite(f domain.Favourite) error {
	existing, _ := s.repo.Exists(f.UserID, f.AssetID)
	if existing {
		return fmt.Errorf("asset already favourited")
	}

	fav := mapper.FavouriteEntityFromDomain(f)
	fav.CreatedAt = time.Now().UTC()
	return s.repo.Add(fav)
}

func (s FavouriteServiceImpl) DeleteFavourite(userID, assetID string) error {
	return s.repo.Delete(userID, assetID)
}
