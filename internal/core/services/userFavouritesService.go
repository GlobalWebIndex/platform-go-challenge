package services

import (
	"context"
	"github.com/loukaspe/platform-go-challenge/internal/core/domain"
	"github.com/loukaspe/platform-go-challenge/internal/core/ports"
	apierrors "github.com/loukaspe/platform-go-challenge/pkg/errors"
	"github.com/loukaspe/platform-go-challenge/pkg/logger"
	"net/http"
)

type UserFavouriteServiceInterface interface {
	AddAsset(
		ctx context.Context,
		userId uint,
		assetId uint,
		assetType string,
	) error

	EditAssetDescription(
		ctx context.Context,
		userId uint,
		assetId uint,
		assetType string,
		description string,
	) error

	RemoveAsset(
		ctx context.Context,
		userId uint,
		assetId uint,
		assetType string,
	) error

	GetAssets(
		ctx context.Context,
		userId uint,
	) ([]domain.Asset, error)
}

type UserFavouriteService struct {
	logger     logger.LoggerInterface
	repository ports.UserRepositoryInterface
}

func NewUserFavouriteService(
	logger logger.LoggerInterface,
	repository ports.UserRepositoryInterface,
) *UserFavouriteService {
	return &UserFavouriteService{
		logger:     logger,
		repository: repository,
	}
}

func (u UserFavouriteService) AddAsset(
	ctx context.Context,
	userId uint,
	assetId uint,
	assetType string,
) error {
	return u.repository.AddFavouriteAsset(ctx, uint32(userId), uint32(assetId), assetType)
}

func (u UserFavouriteService) EditAssetDescription(
	ctx context.Context,
	userId uint,
	assetId uint,
	assetType string,
	description string,
) error {
	return u.repository.EditFavouriteAssetDescription(
		ctx,
		uint32(userId),
		uint32(assetId),
		assetType,
		description,
	)
}

func (u UserFavouriteService) RemoveAsset(
	ctx context.Context,
	userId uint,
	assetId uint,
	assetType string,
) error {
	return u.repository.RemoveFavouriteAsset(ctx, uint32(userId), uint32(assetId), assetType)
}

func (u UserFavouriteService) GetAssets(
	ctx context.Context,
	userId uint,
) ([]domain.Asset, error) {
	user, err := u.repository.GetUser(ctx, uint32(userId))
	if err != nil {
		return nil, err
	}

	if len(user.FavouriteAssets) == 0 {
		return nil, apierrors.NoFavouriteAssetsErrorWrapper{
			ReturnedStatusCode: http.StatusNoContent,
		}
	}

	return user.FavouriteAssets, nil
}
