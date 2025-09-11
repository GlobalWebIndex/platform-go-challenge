package service

import (
	"platform-go-challenge/internal/domain/entities"
	"platform-go-challenge/internal/repository"

	"go.uber.org/zap"
)

// FavouriteServiceIntf defines the interface for our FavouriteService.
type FavouriteServiceIntf interface {
	GetFavourites(userID string, page, size int) ([]entities.Asset, error)
	AddFavourite(userID string, assetID string) error
	UpdateFavouriteDescription(userID string, assetID string, description string) error
	RemoveFavourite(userID string, assetID string) error
}

type FavouriteService struct {
	repo   repository.Repository
	logger *zap.Logger
}

func NewFavouriteService(r repository.Repository, l *zap.Logger) *FavouriteService {
	return &FavouriteService{
		repo:   r,
		logger: l,
	}
}

// TODO check what Ill do with the ctx context.Context
// GetFavourites now returns a slice of a different type and an error.
func (s *FavouriteService) GetFavourites(userID string, page, size int) ([]entities.Asset, error) {
	// A real implementation would call the repository here.
	s.logger.Info("retrieving favourites", zap.String("userID", userID))

	result, err := s.repo.GetFavourites(userID, page, size)
	if err != nil {
		s.logger.Error("failed to get favourites", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (s *FavouriteService) AddFavourite(userID string, assetID string) error {
	s.logger.Info("adding favourite", zap.String("userID", userID), zap.String("assetID", ""))
	// A real implementation would call the repository here.
	if err := s.repo.AddFavourite(userID, assetID); err != nil {
		return err
	}
	return nil
}

func (s *FavouriteService) UpdateFavouriteDescription(userID string, assetID string, description string) error {
	s.logger.Info("updating favourite description", zap.String("userID", userID), zap.String("assetID", assetID))

	if err := s.repo.UpdateFavouriteDescription(userID, assetID, description); err != nil {
		return err
	}
	return nil
}

func (s *FavouriteService) RemoveFavourite(userID string, assetID string) error {
	s.logger.Info("removing favourite", zap.String("userID", userID), zap.String("assetID", assetID))

	if err := s.repo.RemoveFavourite(userID, assetID); err != nil {
		return err
	}
	return nil
}
