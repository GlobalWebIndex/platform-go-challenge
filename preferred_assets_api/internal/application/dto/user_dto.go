package dto

// CreateUserRequest represents a request to create a new user
// swagger:model CreateUserRequest
type CreateUserRequest struct {
	// The user's full name
	// required: true
	// example: "Alice Johnson"
	Name string `json:"name" validate:"required,min=2"`

	// The user's email address
	// required: true
	// example: "alice@example.com"
	Email string `json:"email" validate:"required,email"`

	// The user's password
	// required: true
	// example: "P@ssword123"
	Password string `json:"password" validate:"required,min=8"`
}

// UpdateUserRequest represents a request to update an existing user
// swagger:model UpdateUserRequest
type UpdateUserRequest struct {
	// The user's email address
	// example: "alice@example.com"
	Email string `json:"email,omitempty" validate:"omitempty,email"`

	// The user's full name
	// example: "Alice Johnson"
	Name string `json:"name,omitempty"`

	// The user's password
	// example: "NewP@ssword123"
	Password string `json:"password,omitempty" validate:"required,min=8"`
}

// UserResponse represents a user returned by the API
// swagger:model UserResponse
type UserResponse struct {
	// The user's unique ID
	// example: "user_123"
	ID string `json:"id"`

	// The user's email
	// example: "alice@example.com"
	Email string `json:"email"`

	// The user's full name
	// example: "Alice Johnson"
	Name string `json:"name"`
}

// UserFavouritesResponse represents a user's favourite asset
// swagger:model UserFavouritesResponse
type UserFavouritesResponse struct {
	// The type of the asset (e.g., audience, chart, insight)
	// example: "chart"
	AssetType string `json:"asset_type"`
	// The title of the asset
	// example: "Monthly Sales Chart"
	Title string `json:"title"`
	// Description of the asset
	// example: "Shows monthly sales performance across regions"
	Description string `json:"description"`
}
