package dto

import (
	"errors"
	"fmt"
	"github.com/Kercyn/crud_template/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID      primitive.ObjectID
	Name    string
	Surname string
	Assets  []Asset
}

func (u *User) ToDomain() (domain.User, error) {
	var assets []domain.Asset

	for _, a := range u.Assets {
		domainAsset, err := a.ToDomain()
		if err != nil {
			return domain.User{}, err
		}
		assets = append(assets, domainAsset)
	}

	return domain.User{
		Name:    u.Name,
		Surname: u.Surname,
		Assets:  assets,
	}, nil
}

func (u *User) UnmarshalBSON(data []byte) error {
	auxiliary := new(struct {
		ID      primitive.ObjectID `bson:"_id"`
		Name    string             `bson:"name"`
		Surname string             `bson:"surname"`
	})
	if err := bson.Unmarshal(data, &auxiliary); err != nil {
		return err
	}

	u.ID = auxiliary.ID
	u.Name = auxiliary.Name
	u.Surname = auxiliary.Surname

	assets, err := u.getAssetsFromRawData(data)
	if err != nil {
		return err
	}

	u.Assets = assets

	return nil
}

func (u *User) getAssetsFromRawData(data []byte) ([]Asset, error) {
	assetData := new(struct {
		Assets []bson.Raw `bson:"assets,omitempty"`
	})
	if err := bson.Unmarshal(data, assetData); err != nil {
		return nil, err
	}

	var assets []Asset
	for _, asset := range assetData.Assets {
		assetType, err := u.getAssetType(asset)
		if err != nil {
			return nil, err
		}

		unmarshalledAsset, err := u.unmarshalAssetBasedOnType(asset, assetType)
		if err != nil {
			return nil, err
		}

		assets = append(assets, unmarshalledAsset)
	}

	return assets, nil
}

func (u *User) getAssetType(assetData bson.Raw) (domain.AssetType, error) {
	typeData := new(struct {
		Type string `bson:"type"`
	})

	if err := bson.Unmarshal(assetData, &typeData); err != nil {
		return domain.AssetTypeInvalid, err
	}

	return domain.NewAssetTypeFromString(typeData.Type)
}

func (u *User) unmarshalAssetBasedOnType(
	assetData bson.Raw,
	assetType domain.AssetType,
) (Asset, error) {
	baseAssetData, err := u.unmarshalBaseAssetData(assetData)
	if err != nil {
		return &AssetBase{}, err
	}
	switch assetType {
	case domain.AssetTypeAudience:
		var audience Audience
		if err := bson.Unmarshal(assetData, &audience); err != nil {
			return &Audience{}, err
		}
		audience.AssetBase = baseAssetData
		return &audience, nil
	case domain.AssetTypeBarChart:
		var chart BarChart

		if err := bson.Unmarshal(assetData, &chart); err != nil {
			return &BarChart{}, err
		}
		chart.AssetBase = baseAssetData
		return &chart, nil
	case domain.AssetTypeInsight:
		var insight Insight

		if err := bson.Unmarshal(assetData, &insight); err != nil {
			return &Insight{}, err
		}
		insight.AssetBase = baseAssetData
		return &insight, nil
	default:
		return &AssetBase{}, errors.New(fmt.Sprintf("Invalid asset type: %d", assetType))
	}
}

func (u *User) unmarshalBaseAssetData(assetData bson.Raw) (AssetBase, error) {
	var assetImpl AssetBase
	err := bson.Unmarshal(assetData, &assetImpl)

	if err != nil {
		return AssetBase{}, err
	}

	return assetImpl, nil
}
