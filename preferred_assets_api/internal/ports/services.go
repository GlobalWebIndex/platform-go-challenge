package ports

import "github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"

type UserService interface {
	CreateUser(user domain.User) error
	GetUserByID(id string) (*domain.User, error)
	GetAllUsers() ([]domain.User, error)
	UpdateUser(user domain.User) error
	DeleteUser(id string) error
	GetFavouritesByUser(id string) ([]domain.Favourite, error)
}

type AssetService interface {
	CreateAsset(asset domain.Asset) (domain.Asset, error)
	DeleteAsset(id string) error
}

type FavouriteService interface {
	CreateFavourite(favourite domain.Favourite) error
	DeleteFavourite(userId string, assetId string) error
}
