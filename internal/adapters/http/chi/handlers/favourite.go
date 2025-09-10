package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/SokratisChaimanas/platform-go-challenge/internal/app"
	"github.com/SokratisChaimanas/platform-go-challenge/internal/domain"
	"github.com/google/uuid"
)

// FavouritesHandler serves favourites endpoints.
type FavouritesHandler struct {
	favService *app.FavouritesService
}

func NewFavouritesHandler(favService *app.FavouritesService) *FavouritesHandler {
	return &FavouritesHandler{favService: favService}
}

// ListByUser godoc
// @Summary      List favourites for a user
// @Description  Returns assets the user has favourited using keyset pagination.
// @Tags         favourites
// @Accept       json
// @Produce      json
// @Param        user_id  path   string  true   "User ID (UUID)"
// @Param        limit    query  int     false  "Max items to return (default 20, max 50)"
// @Param        after    query  string  false  "Opaque cursor from next_after"
// @Success      200      {object}  handlers.AssetsListResponse
// @Failure      400      {object}  handlers.ErrorResponse
// @Failure      404      {object}  handlers.ErrorResponse
// @Failure      500      {object}  handlers.ErrorResponse
// @Router       /users/{user_id}/favourites [get]
func (handler *FavouritesHandler) ListByUser(writer http.ResponseWriter, req *http.Request) {
	// Returns assets + a next_after cursor. Uses keyset pagination, not offset.
	writer.Header().Set("Content-Type", "application/json")

	userID, ok := parseUUIDParam(writer, req, "user_id")
	if !ok {
		return
	}

	// Reuse your parsePagination to keep limit defaults/caps consistent; ignore offset.
	limit, _, ok := parsePagination(writer, req)
	if !ok {
		return
	}

	after := req.URL.Query().Get("after")

	items, nextAfter, err := handler.favService.ListByUserKeyset(req.Context(), userID, limit, after)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			WriteJsonError(writer, "user not found", http.StatusNotFound)
			return
		case errors.Is(err, domain.ErrBadCursor):
			WriteJsonError(writer, "bad cursor", http.StatusBadRequest)
			return
		default:
			WriteJsonError(writer, "internal error", http.StatusInternalServerError)
			return
		}
	}

	out := make([]AssetResponse, 0, len(items))
	for _, a := range items {
		out = append(out, AssetResponse{
			ID:          a.ID,
			Type:        a.Type,
			Description: a.Description,
			Payload:     a.Payload,
		})
	}

	_ = json.NewEncoder(writer).Encode(AssetsListResponse{
		Items:     out,
		NextAfter: nextAfter,
	})
}

// Add godoc
// @Summary      Add favourite
// @Description  Adds an asset to the user's favourites.
// @Tags         favourites
// @Accept       json
// @Produce      json
// @Param        user_id  path   string                         true  "User ID (UUID)"
// @Param        payload  body   handlers.FavouriteAddRequest    true  "Favourite payload"
// @Success      201      {object} handlers.FavouriteResponse
// @Failure      400      {object} handlers.ErrorResponse
// @Failure      404      {object} handlers.ErrorResponse
// @Failure      409      {object} handlers.ErrorResponse
// @Failure      500      {object} handlers.ErrorResponse
// @Router       /users/{user_id}/favourites [post]
func (handler *FavouritesHandler) Add(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	userID, ok := parseUUIDParam(writer, req, "user_id")
	if !ok {
		return
	}

	var body struct {
		AssetID string `json:"asset_id"`
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		WriteJsonError(writer, "invalid json", http.StatusBadRequest)
		return
	}

	assetID, err := uuid.Parse(body.AssetID)
	if err != nil {
		WriteJsonError(writer, "invalid asset_id", http.StatusBadRequest)
		return
	}

	f, err := handler.favService.Add(req.Context(), userID, assetID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			WriteJsonError(writer, "user not found", http.StatusNotFound)
			return
		case errors.Is(err, domain.ErrAssetNotFound):
			WriteJsonError(writer, "asset not found", http.StatusNotFound)
			return
		case errors.Is(err, domain.ErrFavouriteAlreadyExists):
			WriteJsonError(writer, "favourite already exists", http.StatusConflict)
			return
		default:
			WriteJsonError(writer, "internal error", http.StatusInternalServerError)
			return
		}
	}

	type resp struct {
		UserID    uuid.UUID `json:"user_id"`
		AssetID   uuid.UUID `json:"asset_id"`
		CreatedAt string    `json:"created_at"`
	}

	writer.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(writer).Encode(resp{
		UserID:    f.UserID,
		AssetID:   f.AssetID,
		CreatedAt: f.CreatedAt.UTC().Format(time.RFC3339),
	})
}

// Remove godoc
// @Summary      Remove favourite
// @Description  Removes an asset from the user's favourites.
// @Tags         favourites
// @Accept       json
// @Produce      json
// @Param        user_id   path   string  true  "User ID (UUID)"
// @Param        asset_id  path   string  true  "Asset ID (UUID)"
// @Success      204
// @Failure      400       {object} handlers.ErrorResponse
// @Failure      404       {object} handlers.ErrorResponse
// @Failure      500       {object} handlers.ErrorResponse
// @Router       /users/{user_id}/favourites/{asset_id} [delete]
func (handler *FavouritesHandler) Remove(writer http.ResponseWriter, req *http.Request) {
	userID, ok := parseUUIDParam(writer, req, "user_id")
	if !ok {
		return
	}
	assetID, ok := parseUUIDParam(writer, req, "asset_id")
	if !ok {
		return
	}

	if err := handler.favService.Remove(req.Context(), userID, assetID); err != nil {
		if errors.Is(err, domain.ErrFavouriteNotFound) {
			WriteJsonError(writer, "favourite not found", http.StatusNotFound)
			return
		}

		WriteJsonError(writer, "internal error", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
