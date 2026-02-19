package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type AssetType string

const (
	AssetTypeChart    AssetType = "chart"
	AssetTypeInsight  AssetType = "insight"
	AssetTypeAudience AssetType = "audience"
)

type Asset interface {
	GetAssetID() string
	GetType() AssetType
	GetDescription() string
	SetDescription(string)
	GetFavouritedAt() time.Time
}

type BaseAsset struct {
	AssetID      string    `json:"asset_id"`
	Type         AssetType `json:"type"`
	Description  string    `json:"description"`
	FavouritedAt time.Time `json:"favourited_at"`
}

func (b *BaseAsset) GetAssetID() string         { return b.AssetID }
func (b *BaseAsset) GetType() AssetType         { return b.Type }
func (b *BaseAsset) GetDescription() string     { return b.Description }
func (b *BaseAsset) SetDescription(desc string) { b.Description = desc }
func (b *BaseAsset) GetFavouritedAt() time.Time { return b.FavouritedAt }

type Chart struct {
	BaseAsset
	Title      string                   `json:"title"`
	XAxisTitle string                   `json:"x_axis_title"`
	YAxisTitle string                   `json:"y_axis_title"`
	Data       []map[string]interface{} `json:"data"`
}

type Insight struct {
	BaseAsset
	Text string `json:"text"`
}

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

type AgeGroup struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type Audience struct {
	BaseAsset
	Name                  string   `json:"name"`
	Gender                Gender   `json:"gender"`
	BirthCountry          string   `json:"birth_country"`
	AgeGroup              AgeGroup `json:"age_group"`
	DailySocialMediaHours int      `json:"daily_social_media_hours"`
	PurchasesLastMonth    int      `json:"purchases_last_month"`
}

// JSON marshalling
func MarshalAsset(asset Asset) ([]byte, error) {
	return json.Marshal(asset)
}

// JSON unmarshalling
func UnmarshalAsset(data []byte) (Asset, error) {
	var base struct {
		Type AssetType `json:"type"`
	}

	if err := json.Unmarshal(data, &base); err != nil {
		return nil, fmt.Errorf("failed to determine asset type: %w", err)
	}

	if base.Type == "" {
		return nil, fmt.Errorf("missing required field: type")
	}

	var asset Asset
	switch base.Type {
	case AssetTypeChart:
		asset = &Chart{}
	case AssetTypeInsight:
		asset = &Insight{}
	case AssetTypeAudience:
		asset = &Audience{}
	default:
		return nil, fmt.Errorf("unknown asset type: %s", base.Type)
	}

	if err := json.Unmarshal(data, asset); err != nil {
		return nil, fmt.Errorf("failed to unmarshal asset: %w", err)
	}

	if asset.GetAssetID() == "" {
		return nil, fmt.Errorf("missing required field: asset_id")
	}

	return asset, nil
}
