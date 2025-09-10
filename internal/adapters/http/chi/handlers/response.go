package handlers

import (
	"github.com/SokratisChaimanas/platform-go-challenge/internal/domain"
	"github.com/google/uuid"
)

// ErrorResponse is used for error payloads (consistent across endpoints).
type ErrorResponse struct {
	Error string `json:"error" example:"error description"`
}

// --- Users ---

// UserResponse is returned by GET /api/users/{user_id}.
type UserResponse struct {
	ID        uuid.UUID `json:"id" example:"11111111-1111-1111-1111-111111111111"`
	CreatedAt string    `json:"created_at" example:"2025-09-08T12:34:56Z"`
}

// --- Assets ---

// AssetEditRequest is the body for PATCH /api/assets/{asset_id}/description.
type AssetEditRequest struct {
	Description string `json:"description" example:"New description from Swagger"`
}

// AssetResponse is returned when assets are requested.
type AssetResponse struct {
	ID          uuid.UUID        `json:"id" example:"aaaaaaa1-0000-0000-0000-000000000001"`
	Type        domain.AssetType `json:"type" example:"chart"`
	Description string           `json:"description" example:"Daily active users - last 7 days"`
	Payload     map[string]any   `json:"payload" swaggertype:"object"`
}

// --- Favourites ---

// FavouriteAddRequest is the body for POST /api/users/{user_id}/favourites.
type FavouriteAddRequest struct {
	AssetID string `json:"asset_id" example:"aaaaaaa1-0000-0000-0000-000000000001"`
}

// FavouriteResponse is returned when creating favourites.
type FavouriteResponse struct {
	UserID    uuid.UUID `json:"user_id"  example:"11111111-1111-1111-1111-111111111111"`
	AssetID   uuid.UUID `json:"asset_id" example:"aaaaaaa1-0000-0000-0000-000000000001"`
	CreatedAt string    `json:"created_at" example:"2025-09-08T12:34:56Z"`
}

// AssetsListResponse wraps a list of assets plus the next page cursor.
type AssetsListResponse struct {
	Items     []AssetResponse `json:"items"`
	NextAfter *string         `json:"next_after,omitempty"`
}

// HealthResponse is used by GET /api/healthz.
type HealthResponse struct {
	OK bool `json:"ok" example:"true"`
}
