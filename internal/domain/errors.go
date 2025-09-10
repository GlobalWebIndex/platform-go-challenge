package domain

import "errors"

// Common domain errors. These are returned by services and repositories
// when business rules or lookups fail. They should be mapped to proper
// HTTP responses in the handlers layer.

var (
	// Asset errors
	ErrAssetNotFound    = errors.New("asset not found")
	ErrInvalidAssetType = errors.New("invalid asset type")
	ErrEmptyDescription = errors.New("asset description cannot be empty")

	// User errors
	ErrUserNotFound = errors.New("user not found")

	// Favourite errors
	ErrFavouriteNotFound      = errors.New("favourite not found")
	ErrFavouriteAlreadyExists = errors.New("favourite already exists")
	ErrBadCursor              = errors.New("bad cursor")
)
