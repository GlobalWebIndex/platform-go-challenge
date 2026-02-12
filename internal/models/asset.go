// Asset model definitions and methods

package models

import "time"

type AssetType string

const (
	AssetTypeChart    AssetType = "chart"
	AssetTypeInsight  AssetType = "insight"
	AssetTypeAudience AssetType = "audience"
)

type Asset interface {
	GetID() string
	GetType() AssetType
}

type FavouriteAsset struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	AssetType   AssetType `json:"asset_type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Data        Asset     `json:"data"`
}

func (f *FavouriteAsset) GetID() string      { return f.ID }
func (f *FavouriteAsset) GetType() AssetType { return f.AssetType }

type Chart struct {
	ID         string         `json:"id"`
	Title      string         `json:"title"`
	XAxisTitle string         `json:"x_axis_title"`
	YAxisTitle string         `json:"y_axis_title"`
	Data       map[string]any `json:"data"`
}

func (c *Chart) GetID() string      { return c.ID }
func (c *Chart) GetType() AssetType { return AssetTypeChart }

type Insight struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

func (i *Insight) GetID() string      { return i.ID }
func (i *Insight) GetType() AssetType { return AssetTypeInsight }

type Audience struct {
	ID                    string   `json:"id"`
	Gender                []string `json:"gender"`
	BirthCountry          []string `json:"birth_country"`
	AgeGroups             []string `json:"age_groups"`
	SocialMediaHoursDaily string   `json:"social_media_hours_daily"`
	PurchasesLastMonth    int      `json:"purchases_last_month"`
}

func (a *Audience) GetID() string      { return a.ID }
func (a *Audience) GetType() AssetType { return AssetTypeAudience }
