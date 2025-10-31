package dto

import "time"

// AssetRequest represents a request to create or update an asset.
// swagger:model AssetRequest
type AssetRequest struct {
	// Unique identifier of the asset
	// example: 123e4567-e89b-12d3-a456-426614174000
	ID string `json:"id" validate:"required"`

	// Type of the asset (audience, chart, or insight)
	// required: true
	// enum: audience,chart,insight
	// example: chart
	Type string `json:"type" validate:"required,oneof=audience chart insight"`

	// Title of the asset
	// required: true
	// example: Customer Demographics Overview
	Title string `json:"title" validate:"required"`

	// Description of the asset
	// example: This chart shows demographic breakdown by age and gender.
	Description string `json:"description"`

	// Timestamp when the asset was created
	// example: 2025-01-01T12:00:00Z
	CreatedAt time.Time `json:"created_at"`

	// Timestamp when the asset was last updated
	// example: 2025-01-02T15:30:00Z
	UpdatedAt time.Time `json:"updated_at"`

	// Text associated with the insight (optional)
	/// example: This insight highlights key trends in customer behavior.
	Text *string `json:"text,omitempty"`

	// Gender associated with the audience (optional)
	// example: female
	Gender *string `json:"gender,omitempty"`

	// Birth country of the audience (optional)
	// example: Canada
	BirthCountry *string `json:"birth_country,omitempty"`

	// Age group of the audience (optional)
	// example: 25-34
	AgeGroup *string `json:"age_group,omitempty"`

	// Average hours spent on social media per day (optional)
	// example: 4
	HoursSocial *float64 `json:"hours_social,omitempty"`

	// Number of purchases made in the last month (optional)
	// example: 12
	PurchasesLastMo *int `json:"purchases_last_month,omitempty"`

	// Titles for the chart axes (for chart-type assets)
	// example: ["Age Group", "Average Purchases"]
	AxesTitles []string `json:"axes_titles,omitempty"`

	// Data points for the chart (for chart-type assets)
	// example: [[18, 2.3], [25, 3.5], [34, 4.1]]
	Data [][]float64 `json:"data,omitempty"`
}

// internal/application/dto/asset_dto.go

// AssetBaseResponse represents the common base fields for all asset responses
// swagger:model AssetBaseResponse
type AssetBaseResponse struct {
	// ID is the unique identifier for the asset
	// example: 550e8400-e29b-41d4-a716-446655440000
	ID string `json:"id"`

	// Type of the asset (audience, chart, insight)
	// example: audience
	Type string `json:"type"`

	// Title of the asset
	// example: "Young Social Media Users"
	Title string `json:"title"`

	// Description provides more details about the asset
	// example: "Audience segment of young adults active on social media"
	Description string `json:"description"`

	// CreatedAt timestamp when the asset was created
	// example: 2023-10-05T14:30:00Z
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt timestamp when the asset was last updated
	// example: 2023-10-05T14:30:00Z
	UpdatedAt time.Time `json:"updated_at"`
}

// AssetCreationResponse represents the response after successfully creating an asset
// swagger:model AssetCreationResponse
type AssetCreationResponse struct {
	AssetBaseResponse

	// Gender of the audience segment (only for audience assets)
	// example: female
	Gender *string `json:"gender,omitempty"`

	// BirthCountry of the audience segment (only for audience assets)
	// example: US
	BirthCountry *string `json:"birth_country,omitempty"`

	// AgeGroup of the audience segment (only for audience assets)
	// example: 18-24
	AgeGroup *string `json:"age_group,omitempty"`

	// HoursSocial represents hours spent on social media per week (only for audience assets)
	// example: 15
	HoursSocial *float64 `json:"hours_social,omitempty"`

	// PurchasesLastMo represents number of purchases in the last month (only for audience assets)
	// example: 3
	PurchasesLastMo *int `json:"purchases_last_month,omitempty"`

	// AxesTitles contains the titles for chart axes (only for chart assets)
	// example: {"x": "Time", "y": "Revenue"}
	AxesTitles []string `json:"axes_titles,omitempty"`

	// Data contains the chart or insight data (only for chart and insight assets)
	// example: [{"x": "2023-01", "y": 1000}, {"x": "2023-02", "y": 1500}]
	Data interface{} `json:"data,omitempty"`
}

// AssetCreationSuccessResponse represents a successful asset creation response
// swagger:response AssetCreationSuccessResponse
type AssetCreationSuccessResponse struct {
	// in: body
	Body AssetCreationResponse
}

// AssetErrorResponse represents an error response for asset operations
// swagger:response AssetErrorResponse
type AssetErrorResponse struct {
	// in: body
	Body struct {
		// Error message
		// example: "invalid asset type: unknown_type"
		Error string `json:"error"`

		// Error code
		// example: "VALIDATION_ERROR"
		Code string `json:"code"`

		// Detailed error message
		// example: "The provided asset type is not supported"
		Details string `json:"details,omitempty"`
	}
}

// AssetValidationErrorResponse represents a validation error response
// swagger:response AssetValidationErrorResponse
type AssetValidationErrorResponse struct {
	// in: body
	Body struct {
		// Error message
		// example: "request validation failed"
		Error string `json:"error"`

		// Error code
		// example: "VALIDATION_ERROR"
		Code string `json:"code"`

		// Field-level validation errors
		// example: {"type": "asset type is required", "title": "title must be at least 3 characters"}
		Fields map[string]string `json:"fields,omitempty"`
	}
}
