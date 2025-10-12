package api_errors

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrFavoriteNotFound   = errors.New("favorite not found")
	ErrInvalidAssetType   = errors.New("invalid asset type")
	ErrInvalidRequest     = errors.New("invalid request")
	ErrMissingAssetID     = errors.New("missing asset id")
	ErrMissingAssetType   = errors.New("missing asset type")
	ErrInvalidJSON        = errors.New("invalid json")
)