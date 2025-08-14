package domain

import "errors"

var (
	// Asset errors
	ErrAssetNotFound      = errors.New("asset not found")
	ErrInvalidAssetType   = errors.New("invalid asset type")
	ErrAssetAlreadyExists = errors.New("asset already exists")

	// User errors
	ErrUserNotFound  = errors.New("user not found")
	ErrInvalidUserID = errors.New("invalid user ID")

	// Favorite errors
	ErrFavoriteNotFound      = errors.New("favorite not found")
	ErrFavoriteAlreadyExists = errors.New("favorite already exists")
	ErrMaxFavoritesReached   = errors.New("maximum favorites limit reached")

	// Validation errors
	ErrInvalidInput         = errors.New("invalid input")
	ErrMissingRequiredField = errors.New("missing required field")

	// Auth errors
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)
