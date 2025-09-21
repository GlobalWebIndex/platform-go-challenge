package domain

import "time"

type Asset struct {
	Key         string    `json:"_key,omitempty"`
	Type        string    `json:"asset_type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Data        AssetData `json:"data"`
}

type AssetData struct {
	ChartData    *Chart    `json:"chart,omitempty"`
	InsightData  *Insight  `json:"insight,omitempty"`
	AudienceData *Audience `json:"audience,omitempty"`
}

type Chart struct {
	Title      string  `json:"title"`
	XAxisTitle string  `json:"x_axis_title"`
	YAxisTitle string  `json:"y_axis_title"`
	Data       []Point `json:"data"`
}

type Point struct {
	X string `json:"x"`
	Y string `json:"y"`
}

type Insight struct {
	Text string `json:"text"`
}

type Audience struct {
	Gender             string `json:"gender"`
	BirthCountry       string `json:"birth_country"`
	AgeGroup           string `json:"age_group"`
	SocialMediaHours   int    `json:"social_media_hours"`
	PurchasesLastMonth int    `json:"purchases_last_month"`
}

type AssetRepository interface {
	CreateAssets() ([]Asset, error)
	GetAsset(assetID string) (string, error)
}
