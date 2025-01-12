package types

import (
	"errors"
	"time"
)

var (
	ErrDescriptionTooLong = errors.New("description exceeds maximum length of 500 characters")
)

type AssetType string

const (
	AssetTypeChart    AssetType = "chart"
	AssetTypeInsight  AssetType = "insight"
	AssetTypeAudience AssetType = "audience"
)

type Asset interface {
	GetID() string
	GetType() AssetType
	GetDescription() string
	SetDescription(desc string) error
	GetCreatedAt() string
	GetUpdatedAt() string
}

type BaseAsset struct {
	ID          string    `json:"id"`
	Type        AssetType `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (a *BaseAsset) GetID() string { return a.ID }

func (a *BaseAsset) GetType() AssetType { return a.Type }

func (a *BaseAsset) GetDescription() string { return a.Description }

func (a *BaseAsset) GetCreatedAt() string { return a.CreatedAt.String() }

func (a *BaseAsset) GetUpdatedAt() string { return a.UpdatedAt.String() }

func (a *BaseAsset) SetDescription(desc string) error {
	if len(desc) > 500 {
		return ErrDescriptionTooLong
	}
	a.Description = desc
	a.UpdatedAt = time.Now()
	return nil
}

type Chart struct {
	BaseAsset
	Title     string                 `json:"title"`
	AxisX     string                 `json:"axis_x"`
	AxisY     string                 `json:"axis_y"`
	ChartData map[string]interface{} `json:"chart_data"`
}

type Insight struct {
	BaseAsset
	Content string `json:"content"`
}

type AgeGroup struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type Audience struct {
	BaseAsset
	Gender             string     `json:"gender"`
	BirthCountry       string     `json:"birth_country"`
	AgeGroups          []AgeGroup `json:"age_groups"`
	SocialMediaHours   float64    `json:"social_media_hours"`
	PurchasesLastMonth int        `json:"purchases_last_month"`
}

type Favorite struct {
	UserID    string    `json:"user_id"`
	AssetID   string    `json:"asset_id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}
