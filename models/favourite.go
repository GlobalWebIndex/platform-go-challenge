package models

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Favourite A Favourite asset of the user
//
// swagger:model Favourite
type Favourite struct {
	ID          primitive.ObjectID `json:"id"         bson:"_id,omitempty"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	AssetType   string             `json:"asset_type"  bson:"asset_type" validate:"required,oneof=Chart Insight Audience"`
	AssetId     primitive.ObjectID `json:"asset_id"    bson:"asset_id" validate:"required"`
	Asset       AssetInterface     `json:"asset"`
	Role        string             `json:"-"           validate:"required"`
	CreatedAt   time.Time          `json:"created_at"  bson:"created_at"`
	ModifiedAt  time.Time          `json:"modified_at" bson:"modified_at"`
}

// Make sure these match the types of assets that exist
// And more importantly the Favourite.AssetType validation.oneof list
const (
	CHART_ASSET         = "Chart"
	INSIGHT_ASSET       = "Insight"
	AUDIENCE_ASSET      = "Audience"
	CHART_COLLECTION    = "charts"
	INSIGHT_COLLECTION  = "insights"
	AUDIENCE_COLLECTION = "audiences"
)

func (f *Favourite) GetAssetCollectionByType() string {

	switch f.AssetType {
	case CHART_ASSET:
		return CHART_COLLECTION
	case INSIGHT_ASSET:
		return INSIGHT_COLLECTION
	case AUDIENCE_ASSET:
		return AUDIENCE_COLLECTION
	default:
		panic("invalid asset type" + f.AssetType)
	}
	return ""
}

func (f *Favourite) EvaluateAssetType() {

	switch f.AssetType {
	case CHART_ASSET:
		f.Asset = new(Chart)
	case INSIGHT_ASSET:
		f.Asset = new(Insight)
	case AUDIENCE_ASSET:
		f.Asset = new(Audience)
	default:
		panic("invalid asset type" + f.AssetType)
	}
}

func (f *Favourite) UnmarshalBSON(data []byte) error {

	var raw map[string]interface{}
	err := bson.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	fmt.Println(raw)
	var ok bool

	f.AssetType, ok = raw["asset_type"].(string)
	f.ID, ok = raw["_id"].(primitive.ObjectID)
	f.Title, ok = raw["title"].(string)
	f.AssetId, ok = raw["asset_id"].(primitive.ObjectID)
	f.Role, ok = raw["role"].(string)
	fOn, ok := raw["created_at"].(primitive.DateTime)
	f.CreatedAt = fOn.Time()
	fOn, ok = raw["modified_at"].(primitive.DateTime)
	f.ModifiedAt = fOn.Time()
	f.Description, _ = raw["description"].(string)
	if !ok {
		return errors.New("invalid favourite asset")
	}

	assetBytes, err := bson.Marshal(raw["asset"])

	if err != nil {
		return err
	}
	switch f.AssetType {
	case CHART_ASSET:
		var a Chart
		err = bson.Unmarshal(assetBytes, &a)
		if err != nil {
			return err
		}
		f.Asset = a
	case INSIGHT_ASSET:
		var a Insight
		err = bson.Unmarshal(assetBytes, &a)
		if err != nil {
			return err
		}
		f.Asset = a
	case AUDIENCE_ASSET:
		var a Audience
		err = bson.Unmarshal(assetBytes, &a)
		if err != nil {
			return err
		}
		f.Asset = a
	default:
		return errors.New("invalid asset type" + f.AssetType)
	}

	return nil

}
