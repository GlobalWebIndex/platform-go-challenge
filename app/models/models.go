package models

import (
	"app/models/constants"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Asset struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Description string             `bson:"description"`
	Type        int                `bson:"type"`
}

type Chart struct {
	Asset      `bson:",inline"`
	Title      string `bson:"title"`
	AxesTitles string `bson:"axesTitles"`
	PlotData   string `bson:"plotData"`
}

type Insight struct {
	Asset `bson:",inline"`
	Text  string `bson:"text"`
}

type Audience struct {
	Asset `bson:",inline"`
}

type AudienceCharacteristic struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty"`
	AssetID             string             `bson:"assetId"`
	CharacteristicID    int                `bson:"characteristicId"`
	CharacteristicValue int                `bson:"characteristicValue"`
}

type User struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
}

type UserFavorite struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	UserID  string             `bson:"userId"`
	AssetID string             `bson:"assetId"`
}

type AssetUpdateRequest struct {
	AssetID     string
	Description string
}

type AssetVO struct {
	ID              string
	Description     string
	Type            int
	Title           string
	AxesTitles      string
	PlotData        string
	Text            string
	Characteristics []AudienceCharacteristicVO
}

type AudienceCharacteristicVO struct {
	Key   int
	Value int
}

func NewChartAsset(assetId string, description string, title string, axesTitles string, plotData string) AssetVO {
	return AssetVO{
		ID: assetId,
		Description: description,
		Type:  constants.ASSET_TYPE_CHART,
		Title: title,
		AxesTitles: axesTitles,
		PlotData: plotData,
	}
}

func NewInsightAsset(assetId string, description string, text string) AssetVO {
	return AssetVO{
		ID: assetId,
		Description: description,
		Type:  constants.ASSET_TYPE_INSIGHT,
		Text: text,
	}
}

func NewAudienceAsset(assetId string, description string, audienceCharacteristics []AudienceCharacteristic) AssetVO {
	characteristics := make([]AudienceCharacteristicVO, 0, len(audienceCharacteristics))
	for _, au := range audienceCharacteristics {
		characteristics = append(characteristics, AudienceCharacteristicVO{Key: au.CharacteristicID, Value: au.CharacteristicValue})
	}
	return AssetVO{
		ID: assetId,
		Description: description,
		Type:  constants.ASSET_TYPE_AUDIENCE,
		Characteristics: characteristics,
	}
}
