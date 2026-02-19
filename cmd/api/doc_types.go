package main

import "github.com/gwi/favourites-api/internal/pagination"

// --- Swagger/OpenAPI helper types (doc-only) ---
//
// These exist to help swaggo generate correct schemas.
// They mirror JSON envelopes and asset payloads.
//

// ErrorResponse represents a generic API error payload.
type ErrorResponse struct {
	Error string `json:"error" example:"invalid cursor"`
}

// EnvelopeResponse represents { "data": ... }
type EnvelopeResponse struct {
	Data any `json:"data"`
}

// PaginatedFavouritesResponse represents paginated response.
type PaginatedFavouritesResponse struct {
	Data []any           `json:"data"`
	Meta pagination.Meta `json:"meta"`
}

// EnvelopeAsset is the { "data": <asset> } response shape.
type Enveloperesponse struct {
	Data any `json:"data"`
}

// EnvelopeHealth represents /health response.
type EnvelopeHealth struct {
	Status  string `json:"status" example:"healthy"`
	Version string `json:"version" example:"1.0.0"`
	Env     string `json:"env" example:"development"`
	Cache   struct {
		HitRate string `json:"hit_rate" example:"90.00%"`
		Hits    int    `json:"hits" example:"10"`
		Misses  int    `json:"misses" example:"1"`
		Size    int    `json:"size" example:"1"`
	} `json:"cache"`
}

// CreateFavouriteRequest documents the POST body.
// It's a union of all asset types; Swagger can't express unions well,
// so we document via examples in the handler annotation.
type CreateFavouriteRequest map[string]any

// UpdateDescriptionRequest documents PUT body.
type UpdateDescriptionRequest struct {
	Description string `json:"description" example:"Updated description"`
}
