package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/GlobalWebIndex/platform-go-challenge/internal/domain/favourites"
)

func TestMapServiceError(t *testing.T) {
	t.Parallel()

	unknownErr := errors.New("boom")

	tests := []struct {
		name               string
		err                error
		wantStatus         int
		wantType           string
		wantTitle          string
		wantDetailNonEmpty bool
	}{
		{
			name:               "userID_required",
			err:                favourites.ErrUserIDRequired,
			wantStatus:         http.StatusBadRequest,
			wantType:           "invalid_request",
			wantTitle:          "invalid request",
			wantDetailNonEmpty: true,
		},
		{
			name:               "assetID_required",
			err:                favourites.ErrAssetIDRequired,
			wantStatus:         http.StatusBadRequest,
			wantType:           "invalid_request",
			wantTitle:          "invalid request",
			wantDetailNonEmpty: true,
		},
		{
			name:               "description_required",
			err:                favourites.ErrDescriptionRequired,
			wantStatus:         http.StatusBadRequest,
			wantType:           "invalid_request",
			wantTitle:          "invalid request",
			wantDetailNonEmpty: true,
		},
		{
			name:               "invalid_asset_payload",
			err:                favourites.ErrInvalidAssetPayload,
			wantStatus:         http.StatusUnprocessableEntity,
			wantType:           "invalid_asset",
			wantTitle:          "invalid asset payload",
			wantDetailNonEmpty: true,
		},
		{
			name:               "invalid_asset_type",
			err:                favourites.ErrInvalidAssetType,
			wantStatus:         http.StatusUnprocessableEntity,
			wantType:           "invalid_asset",
			wantTitle:          "invalid asset payload",
			wantDetailNonEmpty: true,
		},
		{
			name:               "favourite_not_found",
			err:                favourites.ErrFavouriteNotFound,
			wantStatus:         http.StatusNotFound,
			wantType:           "not_found",
			wantTitle:          "favourite not found",
			wantDetailNonEmpty: true,
		},
		{
			name:               "unknown_error_defaults_to_500",
			err:                unknownErr,
			wantStatus:         http.StatusInternalServerError,
			wantType:           "internal_error",
			wantTitle:          "internal server error",
			wantDetailNonEmpty: false, // we return empty detail for unknown errors
		},
	}

	for _, tt := range tests {
		tt := tt // capture
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			status, typ, title, detail := mapServiceError(tt.err)

			if status != tt.wantStatus {
				t.Fatalf("status = %d, want %d", status, tt.wantStatus)
			}
			if typ != tt.wantType {
				t.Fatalf("type = %q, want %q", typ, tt.wantType)
			}
			if title != tt.wantTitle {
				t.Fatalf("title = %q, want %q", title, tt.wantTitle)
			}
			if tt.wantDetailNonEmpty && detail == "" {
				t.Fatalf("detail is empty, expected non-empty")
			}
			if !tt.wantDetailNonEmpty && detail != "" {
				t.Fatalf("detail = %q, expected empty", detail)
			}
		})
	}
}
