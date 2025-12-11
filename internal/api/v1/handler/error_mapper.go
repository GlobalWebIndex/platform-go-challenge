package handler

import (
	"errors"
	"net/http"

	"github.com/GlobalWebIndex/platform-go-challenge/internal/domain/favourites"
)

// mapServiceError maps domain/service errors to HTTP Problem Details fields.
func mapServiceError(err error) (status int, typ, title, detail string) {
	switch {
	case errors.Is(err, favourites.ErrUserIDRequired),
		errors.Is(err, favourites.ErrAssetIDRequired),
		errors.Is(err, favourites.ErrDescriptionRequired):
		return http.StatusBadRequest, "invalid_request", "invalid request", err.Error()
	case errors.Is(err, favourites.ErrInvalidAssetPayload),
		errors.Is(err, favourites.ErrInvalidAssetType):
		return http.StatusUnprocessableEntity, "invalid_asset", "invalid asset payload", err.Error()
	case errors.Is(err, favourites.ErrFavouriteNotFound):
		return http.StatusNotFound, "not_found", "favourite not found", err.Error()
	default:
		return http.StatusInternalServerError, "internal_error", "internal server error", ""
	}
}
