package errors

import (
	"encoding/json"
	"net/http"
)

type HTTPError struct {
	Status  int
	Message string
}

func (e *HTTPError) Error() string {
	return e.Message
}

var (
	ErrUnauthorized      = &HTTPError{Status: http.StatusUnauthorized, Message: "Unauthorized"}
	ErrForbidden         = &HTTPError{Status: http.StatusForbidden, Message: "Forbidden"}
	ErrBadRequest        = &HTTPError{Status: http.StatusBadRequest, Message: "Bad request"}
	ErrNotFound          = &HTTPError{Status: http.StatusNotFound, Message: "Not found"}
	ErrInternal          = &HTTPError{Status: http.StatusInternalServerError, Message: "Internal server error"}
	ErrAssetExists       = &HTTPError{Status: http.StatusBadRequest, Message: "Asset already exists"}
	ErrUnknownAssetType  = &HTTPError{Status: http.StatusBadRequest, Message: "Unknown asset type"}
	ErrFavouriteExists   = &HTTPError{Status: http.StatusBadRequest, Message: "Favourite already exists"}
	ErrFavouriteNotFound = &HTTPError{Status: http.StatusNotFound, Message: "Favourite not found"}
	ErrUserNotFound      = &HTTPError{Status: http.StatusNotFound, Message: "User not found"}
	ErrInvalidID         = &HTTPError{Status: http.StatusBadRequest, Message: "Invalid ID"}
	ErrInvalidBody       = &HTTPError{Status: http.StatusBadRequest, Message: "Invalid request body"}
	ErrAssetNotFound     = &HTTPError{Status: http.StatusNotFound, Message: "Asset not found"}
	ErrUserExists        = &HTTPError{Status: http.StatusBadRequest, Message: "User already exists"}
	ErrConflict          = &HTTPError{Status: http.StatusConflict, Message: "Already exists"}
	ErrInvalidToken       = &HTTPError{Status: http.StatusBadRequest, Message: "Invalid or expired token"}
)

func WriteError(w http.ResponseWriter, err error) {
	if httpErr, ok := err.(*HTTPError); ok {
		http.Error(w, httpErr.Message, httpErr.Status)
		return
	}
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
