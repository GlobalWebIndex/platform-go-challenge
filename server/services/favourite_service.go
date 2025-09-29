package services

import (
	"time"

	"favourite_assets/server/errors"
	"favourite_assets/server/models"
	"favourite_assets/server/repositories"

	"github.com/google/uuid"
)

type FavouriteService struct {
	repo         *repositories.FavouriteRepository
	userService  *UserService
	assetService *AssetService
}

func NewFavouriteService(
	repo *repositories.FavouriteRepository,
	userService *UserService,
	assetService *AssetService,
) *FavouriteService {
	return &FavouriteService{
		repo:         repo,
		userService:  userService,
		assetService: assetService,
	}
}

func (s *FavouriteService) AddFavourite(userID, assetID uuid.UUID) (*models.Favourite, error) {

	if _, err := s.userService.GetUser(userID); err != nil {
		return nil, errors.ErrUserNotFound
	}

	asset, err := s.assetService.GetAsset(assetID)
	if err != nil {
		return nil, errors.ErrAssetNotFound
	}

	// Check if favourite already exists
	existingFavourites := s.repo.ListByUser(userID)
	for _, fav := range existingFavourites {
		if fav.AssetID == assetID {
			return nil, errors.ErrConflict
		}
	}

	fav := &models.Favourite{
		ID:        uuid.New(),
		UserID:    userID,
		AssetID:   assetID,
		AssetType: asset.GetType(),
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(fav); err != nil {
		return nil, err
	}
	return fav, nil
}

func (s *FavouriteService) RemoveFavourite(favID uuid.UUID) error {
	if err := s.repo.Delete(favID); err != nil {
		return errors.ErrNotFound
	}
	return nil
}

func (s *FavouriteService) ListFavouritesByUser(userID uuid.UUID) ([]*models.Favourite, error) {

	if _, err := s.userService.GetUser(userID); err != nil {
		return nil, errors.ErrUserNotFound
	}

	return s.repo.ListByUser(userID), nil
}

func (s *FavouriteService) GetFavourite(favID uuid.UUID) (*models.Favourite, error) {
	fav, err := s.repo.GetByID(favID)
	if err != nil {
		return nil, errors.ErrFavouriteNotFound
	}
	return fav, nil
}
