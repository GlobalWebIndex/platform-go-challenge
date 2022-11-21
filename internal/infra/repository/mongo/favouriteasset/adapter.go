package favouriteasset

import (
	"encoding/json"
	"errors"
	assetsrv "platform-go-challenge/internal/app/favouriteasset"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// TypeChart is a Chart type.
	TypeChart = "Chart"
	// TypeInsight is an Insight type.
	TypeInsight = "Insight"
	// TypeAudience is an Audience type.
	TypeAudience = "Audience"
)

// repoFavouriteAsset represents the favourite asset in the repository.
type repoFavouriteAsset struct {
	ID          primitive.ObjectID `bson:"_id"`
	UserID      string             `bson:"user_id"`
	Description string             `bson:"description"`
	Asset       []repoAsset        `bson:"user_favourite_assets"`
}

type repoAsset struct {
	ID   primitive.ObjectID `bson:"_id"`
	Type string             `bson:"type"`
	Data interface{}        `bson:"data"`
}

// ChartAsset contains the Chart asset type details.
type ChartAsset struct {
	Title  string      `json:"title" ,bson:"title"`
	XTitle string      `json:"x_title" ,bson:"x_title"`
	YTitle string      `json:"y_title" ,bson:"y_title"`
	Data   interface{} `json:"data" ,bson:"data"`
}

// InsightAsset contains the Insight asset type details.
type InsightAsset struct {
	Text string `json:"text" ,bson:"text"`
}

// AudienceAsset contains the Audience asset type details.
type AudienceAsset struct {
	Gender                       string `json:"gender" ,bson:"gender"`
	BirthCountry                 string `json:"birth_country" ,bson:"birth_country"`
	AgeGroups                    string `json:"age_groups" ,bson:"age_groups"`
	HoursSpentDailyOnSocialMedia string `json:"hours_spent_daily_on_social_media" ,bson:"hours_spent_daily_on_social_media"`
	NumberOfPurchasesLastMonth   string `json:"number_of_purchases_last_month" ,bson:"number_of_purchases_last_month"`
}

func (fa repoFavouriteAsset) adaptToModel() (*assetsrv.FavouriteAsset, error) {
	faID := fa.ID.Hex()
	asset, err := transformAsset(fa.Asset)
	if err != nil {
		return nil, err
	}
	return &assetsrv.FavouriteAsset{
		ID:          &faID,
		Description: fa.Description,
		Asset:       asset,
	}, nil
}

func transformAsset(assets []repoAsset) (*assetsrv.Asset, error) {
	if len(assets) == 0 {
		return nil, nil
	}
	trData, err := transformAssetDataByType(assets[0].Type, assets[0].Data)
	if err != nil {
		return nil, err
	}

	return &assetsrv.Asset{
		ID:   assets[0].ID.Hex(),
		Type: assets[0].Type,
		Data: trData,
	}, nil
}

func transformAssetDataByType(t string, data interface{}) (interface{}, error) {
	var res interface{}

	temporaryBytes, err := bson.MarshalExtJSON(data, true, true)
	if err != nil {
		return nil, err
	}

	switch t {
	case TypeChart:
		chAsset := ChartAsset{}
		err = json.Unmarshal(temporaryBytes, &chAsset)
		if err != nil {
			return nil, err
		}
		res = chAsset
		break
	case TypeInsight:
		inAsset := InsightAsset{}
		err = json.Unmarshal(temporaryBytes, &inAsset)
		if err != nil {
			return nil, err
		}
		res = inAsset
		break
	case TypeAudience:
		auAsset := AudienceAsset{}
		err = json.Unmarshal(temporaryBytes, &auAsset)
		if err != nil {
			return nil, err
		}
		res = auAsset
		break
	default:
		return nil, errors.New("invalid asset type")
	}

	return res, nil
}
