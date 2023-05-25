package dto

import (
	"github.com/Kercyn/crud_template/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Asset interface {
	ToDomain() (domain.Asset, error)
}

type AssetBase struct {
	ID          primitive.ObjectID `bson:"_id"`
	Type        string             `bson:"type"`
	Description string             `bson:"description"`
	IsFavourite bool               `bson:"is_favourite"`
}

func (a AssetBase) ToDomain() (domain.Asset, error) {
	t, err := domain.NewAssetTypeFromString(a.Type)
	if err != nil {
		return domain.AssetBase{}, err
	}
	return domain.AssetBase{
		Type:        t,
		Description: a.Description,
		IsFavourite: a.IsFavourite,
	}, nil
}
