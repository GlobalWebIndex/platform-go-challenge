package dto

import "time"

// FavouriteRequest represents a request to add a favourite
// swagger:model FavouriteRequest
type FavouriteRequest struct {
	// The ID of the user
	// required: true
	// example: "user_123"
	UserId string `json:"_id" validate:"required"`
	// The ID of the asset
	// required: true
	// example: "asset_456"
	AssetId string `json:"asset_id" validate:"required"`
}

// FavouriteResponse represents a favourite returned by the API
// swagger:model FavouriteResponse
type FavouriteResponse struct {

	// The ID of the user
	// example: "user_123"
	UserID string `json:"user_id"`

	// The ID of the asset
	// example: "asset_456"
	AssetID string `json:"asset_id"`

	// The type of the asset (e.g., audience, chart, insight)
	// example: "chart"
	AssetType string `json:"asset_type,omitempty"`

	// Timestamp when the favourite was created
	// example: "2025-10-30T15:04:05Z"
	CreatedAt time.Time `json:"created_at"`

	// The full asset object (can be any type)
	// example: {"id":"asset_456","title":"Sales Chart","description":"Monthly sales chart","type":"chart"}
	Asset any `json:"asset"`
}
