// Package model holds small domain types used across the service: Asset and
// Favourite. Keep these lean and serialisable to JSON for the HTTP API.
package model

import (
	"time"

	"github.com/google/uuid"
)

// AssetType enumerates the supported asset kinds used in favourites.
type AssetType string

const (
	// ChartType represents a chart asset.
	ChartType AssetType = "chart"
	// InsightType represents a textual insight asset.
	InsightType AssetType = "insight"
	// AudienceType represents an audience definition asset.
	AudienceType AssetType = "audience"
)

// ChartPoint represents one data point in a chart (e.g., "2025-01" → 10000).
type ChartPoint struct {
	X string  `json:"label"` // e.g., "2025-01"
	Y float64 `json:"value"` // e.g., 10000.0
}

// ChartData represents the data for a chart asset.
type ChartData struct {
	Title  string       `json:"title"`
	XAxis  string       `json:"x_axis"`
	YAxis  string       `json:"y_axis"`
	Values []ChartPoint `json:"values"`
}

// InsightData represents the data for an insight asset.
type InsightData struct {
	Text string `json:"text"`
}

// AudienceData represents the data for an audience asset.
type AudienceData struct {
	Gender             string  `json:"gender"` // "Male" or "Female"
	BirthCountry       string  `json:"birth_country"`
	AgeGroup           string  `json:"age_group"`          // e.g. "24-35"
	HoursSocialMedia   float64 `json:"hours_social_media"` // hours per day
	PurchasesLastMonth int     `json:"purchases_last_month"`
}

// Asset is a minimal representation of an asset that can be favourited.
// Fields are intentionally simple for the exercise: ID, Type, Description
// and CreatedAt timestamp.
type Asset struct {
	ID          uuid.UUID `json:"id"`
	Type        AssetType `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Data        any       `json:"data,omitempty"` // ChartData, InsightData, or AudienceData
}

// NewAsset creates a new Asset with a generated UUID.
func NewAsset(assetType AssetType, description string, data any) Asset {
	timeNow := time.Now().UTC()
	return Asset{
		Type:        assetType,
		Description: description,
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
		Data:        data,
	}
}
