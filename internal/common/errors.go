package common

const (
	ErrInvalidAssetType  = "invalid asset type"
	ErrEmptyUsername     = "username cannot be empty"
	ErrEmptyPassword     = "password cannot be empty"
	ErrInvalidRole       = "invalid role"
	ErrHashPassword      = "cannot hash password"
	ErrAdminAccess       = "admin access required"
	ErrUnauthorized      = "unauthorized"
	ErrForbiddenAccess   = "access to other users' favorites is forbidden"
	ErrInvalidUserID     = "invalid user id"
	ErrInvalidJSON       = "invalid json"
	ErrMissingAssetID    = "asset id is required"
	ErrAddFavorite       = "could not add favorite"
	ErrRemoveFavorite    = "failed to remove favorite"
	ErrUpdateFavorite    = "failed to edit favorite"
	ErrFetchingFavorites = "error fetching favorites"
	ErrAssetTypeNotFound = "asset type not found"
	ErrAssetIDRequired   = "asset ID is required for editing"
	ErrUserNotFound      = "user not found"
	ErrAssetNotFound     = "asset not found"
	ErrUsernameExists    = "username already exists"
)

// ErrorResponse is used for structured error responses in the API.
type ErrorResponse struct {
	Error string `json:"error"`
}
