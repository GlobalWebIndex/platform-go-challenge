package gwiexercise

import (
	"encoding/json"
	"time"
)

type AssetType string

const (
	AssetChart    AssetType = "chart"
	AssetInsight  AssetType = "insight"
	AssetAudience AssetType = "audience"
)

type Asset interface {
	GetID() string
	GetType() AssetType
	GetDescription() string
	SetDescription(desc string)
	json.Marshaler
}

type BaseAsset struct {
	ID          string    `json:"id"`
	Type        AssetType `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (b *BaseAsset) GetID() string           { return b.ID }
func (b *BaseAsset) GetType() AssetType      { return b.Type }
func (b *BaseAsset) GetDescription() string  { return b.Description }
func (b *BaseAsset) SetDescription(d string) { b.Description = d }

type ChartAsset struct {
	BaseAsset
	Title      string    `json:"title"`
	XAxisTitle string    `json:"xAxisTitle"`
	YAxisTitle string    `json:"yAxisTitle"`
	Data       []float64 `json:"data"`
}

type InsightAsset struct {
	BaseAsset
	Text string `json:"text"`
}

type AudienceAsset struct {
	BaseAsset
	Gender             string `json:"gender"`
	BirthCountry       string `json:"birthCountry"`
	AgeGroup           string `json:"ageGroup"`
	SocialHoursPerDay  string `json:"socialHoursPerDay"`
	PurchasesLastMonth int    `json:"purchasesLastMonth"`
}
