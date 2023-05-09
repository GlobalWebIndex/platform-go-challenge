package ports

import (
	"context"
	"github.com/loukaspe/platform-go-challenge/internal/core/domain"
)

type UserRepositoryInterface interface {
	GetUser(ctx context.Context, userId uint32) (*domain.User, error)
	CreateUser(context.Context, *domain.User) (uint, error)
	UpdateUser(context.Context, uint32, *domain.User) error
	DeleteUser(context.Context, uint32) error
	AddFavouriteAsset(
		ctx context.Context,
		userId uint32,
		assetId uint32,
		assetType string,
	) error
	EditFavouriteAssetDescription(
		ctx context.Context,
		userId uint32,
		assetId uint32,
		assetType,
		description string,
	) error
	RemoveFavouriteAsset(
		ctx context.Context,
		userId uint32,
		assetId uint32,
		assetType string,
	) error
}
