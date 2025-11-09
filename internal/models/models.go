package models

import (
	"encoding/json"
	"time"
)

// AssetType denotes the kind of asset wrapped by a Favourite.
type AssetType string

const (
	AssetChart   AssetType = "chart"
	AssetInsight AssetType = "insight"
	AssetAudience AssetType = "audience"
)

// AssetBase holds fields common to all assets.
type AssetBase struct {
	Type        AssetType `json:"type"`
	Description string    `json:"description,omitempty"`
}

// Chart models a simple numeric chart.
type Chart struct {
	AssetBase
	Title      string    `json:"title"`
	AxisXTitle string    `json:"axis_x_title"`
	AxisYTitle string    `json:"axis_y_title"`
	Data       []float64 `json:"data"`
}

// Insight is a short textual insight.
type Insight struct {
	AssetBase
	Text string `json:"text"`
}

// Audience describes audience characteristics.
type Audience struct {
	AssetBase
	Gender             string   `json:"gender"`
	BirthCountry       string   `json:"birth_country"`
	AgeGroups          []string `json:"age_groups"`
	HoursSocialDaily   float64  `json:"hours_social_daily"`
	PurchasesLastMonth int      `json:"purchases_last_month"`
}

// Favourite is a user-saved asset with metadata.
// Asset keeps the raw JSON to allow payloads per type.
type Favourite struct {
	ID          string          `json:"id"`
	Type        AssetType       `json:"type"`
	Description string          `json:"description,omitempty"`
	Asset       json.RawMessage `json:"asset"`
	CreatedAt   time.Time       `json:"created_at"`
}
