package dto

import (
	"time"

	"github.com/GlobalWebIndex/platform-go-challenge/internal/domain/favourites"
)

// AssetType re-exports the domain AssetType for the public API.
type AssetType = favourites.AssetType

const (
	AssetTypeChart    = favourites.AssetTypeChart
	AssetTypeInsight  = favourites.AssetTypeInsight
	AssetTypeAudience = favourites.AssetTypeAudience
)

// ChartPayload represents the chart-specific payload exposed over HTTP.
type ChartPayload struct {
	Title string    `json:"title"`
	XAxis string    `json:"xAxis"`
	YAxis string    `json:"yAxis"`
	Data  []float64 `json:"data"`
}

// InsightPayload represents the insight-specific payload.
type InsightPayload struct {
	Text string `json:"text"`
}

// AudiencePayload represents the audience-specific payload.
type AudiencePayload struct {
	Gender                   string `json:"gender"`
	BirthCountry             string `json:"birthCountry"`
	AgeGroup                 string `json:"ageGroup"`
	HoursOnSocialMediaPerDay int    `json:"hoursOnSocialMediaPerDay"`
	PurchasesLastMonth       int    `json:"purchasesLastMonth"`
}

// CreateAssetRequest represents an asset to be created inline.
type CreateAssetRequest struct {
	Type        AssetType        `json:"type"`
	Description string           `json:"description"`
	Chart       *ChartPayload    `json:"chart,omitempty"`
	Insight     *InsightPayload  `json:"insight,omitempty"`
	Audience    *AudiencePayload `json:"audience,omitempty"`
}

// CreateFavouriteRequest is the request body for creating a new favourite asset.
type CreateFavouriteRequest struct {
	Asset CreateAssetRequest `json:"asset"`
}

// EditFavouriteDescriptionRequest is the body for updating a favourite's description.
type EditFavouriteDescriptionRequest struct {
	Description string `json:"description"`
}

// AssetResponse is the asset representation returned by the API.
type AssetResponse struct {
	ID          string           `json:"id"`
	Type        AssetType        `json:"type"`
	Description string           `json:"description"`
	Chart       *ChartPayload    `json:"chart,omitempty"`
	Insight     *InsightPayload  `json:"insight,omitempty"`
	Audience    *AudiencePayload `json:"audience,omitempty"`
}

// FavouriteResponse wraps an asset with favourite metadata.
type FavouriteResponse struct {
	Asset     AssetResponse `json:"asset"`
	CreatedAt time.Time     `json:"createdAt"`
}

// CreateFavouriteResponse is the response body for creating a favourite.
type CreateFavouriteResponse struct {
	UserID    string    `json:"userId"`
	AssetID   string    `json:"assetId"`
	CreatedAt time.Time `json:"createdAt"`
}

// FavouriteItemResponse represents a single favourite item in a list response.
type FavouriteItemResponse struct {
	UserID    string        `json:"userId"`
	AssetID   string        `json:"assetId"`
	CreatedAt time.Time     `json:"createdAt"`
	Asset     AssetResponse `json:"asset"`
}

// ListMetaResponse contains pagination metadata for list responses.
type ListMetaResponse struct {
	Total int `json:"total"`
	Count int `json:"count"`
}

// ListFavouritesResponse is the response body for listing favourites.
type ListFavouritesResponse struct {
	Items []FavouriteItemResponse `json:"items"`
	Meta  ListMetaResponse        `json:"meta"`
}
