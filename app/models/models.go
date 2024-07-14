package models

import (
	"fmt"
)

// ID represents a unique identifier for various entities
type ID int

// AssetType represents the type of an asset
type AssetType string

const (
	// ChartAsset represents a chart asset type
	ChartAsset AssetType = "CHART"
	// InsightAsset represents an insight asset type
	InsightAsset AssetType = "INSIGHT"
	// AudienceAsset represents an audience asset type
	AudienceAsset AssetType = "AUDIENCE"
)

// Asset interface defines common methods for all asset types
type Asset interface {
	ID() ID
	Type() AssetType
}

// BaseAsset contains common fields for all asset types
type BaseAsset struct {
	AssetID   ID        `json:"id"`
	AssetType AssetType `json:"asset_type"`
}

func (b BaseAsset) ID() ID          { return b.AssetID }
func (b BaseAsset) Type() AssetType { return b.AssetType }

// FavoriteAsset represents a user's favorite asset
type FavoriteAsset struct {
	BaseAsset
	AssetID ID `json:"asset_id"`
}

// UserFavorites represents a user's favorite assets
type UserFavorites struct {
	UserID    ID               `json:"user_id"`
	Favorites *[]FavoriteAsset `json:"favorites"`
}

type FavoriteAssetDetailed struct {
	ID    ID    `json:"id"`
	Asset Asset `json:"asset"`
}

// UserFavoritesDetailed represents a user's favorite assets with detailed information
type UserFavoritesDetailed struct {
	UserID    ID                       `json:"user_id"`
	Favorites *[]FavoriteAssetDetailed `json:"favorites"`
}

// Chart represents a chart asset
type Chart struct {
	BaseAsset
	Title      string       `json:"title"`
	XAxisTitle string       `json:"x_axis_title"`
	YAxisTitle string       `json:"y_axis_title"`
	DataPoints *[]DataPoint `json:"data_points"`
}

// DataPoint represents a single data point in a chart
type DataPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Insight represents an insight asset
type Insight struct {
	BaseAsset
	Text string `json:"text"`
}

// Audience represents an audience asset
type Audience struct {
	BaseAsset
	Gender                  string  `json:"gender"`
	BirthCountry            string  `json:"birth_country"`
	AgeGroup                string  `json:"age_group"`
	DailyHoursOnSocialMedia float64 `json:"daily_hours_on_social_media"`
	PurchasesLastMonth      int     `json:"purchases_last_month"`
}

// User represents a user of the system
type User struct {
	UserID   ID     `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   string `json:"status"`
}

// Validate checks if the User struct is valid
func (u *User) Validate() error {
	if u.UserID <= 0 {
		return fmt.Errorf("invalid user ID")
	}
	if u.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if u.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	return nil
}
