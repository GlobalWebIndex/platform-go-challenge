package ports

import (
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/entities"
)

type UserRepository interface {
	GetByID(id string) (entities.UserEntity, error)
	Save(user entities.UserEntity) error
	GetAll() ([]entities.UserEntity, error)
	Delete(id string) error
	Update(user entities.UserEntity) error
	GetFavouritesByID(id string) ([]entities.FavouriteEntity, error)
}

type AssetRepository interface {
	Save(asset entities.AssetEntity) (entities.AssetEntity, error)
	GetByID(id string) (entities.AssetEntity, error)
	GetByIDs(ids []string) ([]entities.AssetEntity, error)
	GetAll() ([]entities.AssetEntity, error)
	GetByType(assetType entities.AssetType) ([]entities.AssetEntity, error)
	Update(asset entities.AssetEntity) error
	Delete(id string) error
	Exists(id string) (bool, error)
}

type FavouriteRepository interface {
	Add(f entities.FavouriteEntity) error
	Delete(userID, assetID string) error
	GetByUserID(userID string) ([]entities.FavouriteEntity, error)
	Exists(userID, assetID string) (bool, error)
}
