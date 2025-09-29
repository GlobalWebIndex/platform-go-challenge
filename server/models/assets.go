package models

import (
	"time"

	"github.com/google/uuid"
)

type AssetType string

const (
	AssetChart    AssetType = "chart"
	AssetInsight  AssetType = "insight"
	AssetAudience AssetType = "audience"
)

type Asset interface {
	GetID() uuid.UUID
	GetDescription() string
	SetDescription(desc string)
	GetType() AssetType
}

type BaseAsset struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (b *BaseAsset) GetID() uuid.UUID           { return b.ID }
func (b *BaseAsset) GetDescription() string     { return b.Description }
func (b *BaseAsset) SetDescription(desc string) { b.Description = desc }

type Chart struct {
	BaseAsset
	Title string   `json:"title"`
	XAxis string   `json:"xAxis"`
	YAxis string   `json:"yAxis"`
	Data  [][2]any `json:"data"` 
}

func (c *Chart) GetType() AssetType { return AssetChart }

type Insight struct {
	BaseAsset
	Text string `json:"text"`
}

func (i *Insight) GetType() AssetType { return AssetInsight }

type Audience struct {
	BaseAsset
	Gender             string `json:"gender"`
	BirthCountry       string `json:"birthCountry"`
	AgeGroup           string `json:"ageGroup"`
	HoursOnSocial      int    `json:"hoursOnSocial"`
	PurchasesLastMonth int    `json:"purchasesLastMonth"`
}

func (a *Audience) GetType() AssetType { return AssetAudience }

type Favourite struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"userId"`
	AssetID   uuid.UUID  `json:"assetId"`
	AssetType AssetType  `json:"assetType"` 
	Asset     Asset      `json:"asset,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
}
