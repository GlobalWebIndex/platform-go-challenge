package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gwi/favourites-api/internal/models"
	"github.com/gwi/favourites-api/internal/pagination"
	"github.com/gwi/favourites-api/internal/store"
)

// ==================== Health Check ====================
// healthCheckHandler godoc
//
// @Summary      Health check
// @Description  Service liveness and cache stats.
// @Tags         system
// @Produce      json
// @Success      200 {object} EnvelopeHealth
// @Failure      500 {object} ErrorResponse
// @Router       /health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	cacheStats := app.favouritesService.GetCacheStats()

	data := map[string]interface{}{
		"status":  "healthy",
		"version": version,
		"env":     app.config.env,
		"cache": map[string]interface{}{
			"hit_rate": fmt.Sprintf("%.2f%%", cacheStats.HitRate()*100),
			"hits":     cacheStats.Hits,
			"misses":   cacheStats.Misses,
			"size":     cacheStats.Size,
		},
	}

	if err := writeJSON(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}

// ==================== Favourites Handlers ====================
// getFavouritesHandler godoc
//
// @Summary      List user favourites
// @Description  Returns paginated favourites for a user using cursor-based pagination (keyset).
// @Tags         favourites
// @Produce      json
// @Param        userID  path     string true  "User ID (external identifier)" example(user-123)
// @Param        limit   query    int    false "Items per page (1-100)" minimum(1) maximum(100) default(20)
// @Param        cursor  query    string false "Pagination cursor (favourited_at|asset_id). URL-encode when using curl." example(2026-02-18T15:02:16.236972Z|asset-42)
// @Success      200 {object} PaginatedFavouritesResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /users/{userID}/favourites [get]
func (app *application) getFavouritesHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	// Parse pagination parameters
	cursor := r.URL.Query().Get("cursor")
	if cursor != "" {
		if _, err := pagination.ParseCursor(cursor); err != nil {
			app.badRequestResponse(w, r, fmt.Errorf("invalid cursor: %v", err))
			return
		}
	}
	limit := r.URL.Query().Get("limit")
	params := pagination.ParsePaginationParams(cursor, limit)

	// Get paginated favourites
	favourites, nextCursor, hasMore, err := app.favouritesService.GetFavouritesPaginated(r.Context(), userID, params)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Build paginated response
	response := pagination.NewPaginatedResponse(favourites, nextCursor, hasMore, params.Limit)

	if err := writeJSON(w, http.StatusOK, response); err != nil {
		app.internalServerError(w, r, err)
	}
}

// addFavouriteHandler godoc
//
// @Summary      Add favourite
// @Description  Adds an asset to a user's favourites. asset_id is an external identifier from upstream systems.
//
// @Tags         favourites
// @Accept       json
// @Produce      json
// @Param        userID  path     string true  "User ID (external identifier)" example(user-123)
// @Param        asset   body     CreateFavouriteRequest true "Asset payload"
//
// @Success      201 {object} EnvelopeResponse
// @Failure      400 {object} ErrorResponse
// @Failure      409 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /users/{userID}/favourites [post]
//
// @Example {json} Insight
//
//	{
//	  "asset_id": "insight-1",
//	  "type": "insight",
//	  "description": "Optional description",
//	  "text": "40% of millennials spend more than 3 hours on social media daily"
//	}
//
// @Example {json} Chart
//
//	{
//	  "asset_id": "chart-1",
//	  "type": "chart",
//	  "description": "Optional description",
//	  "title": "Q4 Sales",
//	  "x_axis_title": "Month",
//	  "y_axis_title": "Revenue",
//	  "data": [{"month":"Jan","revenue":1000}]
//	}
//
// @Example {json} Audience
//
//	{
//	  "asset_id": "audience-1",
//	  "type": "audience",
//	  "description": "Optional description",
//	  "name": "Young Male Social Media Users",
//	  "gender": "male",
//	  "birth_country": "US",
//	  "age_group": {"min":24,"max":35},
//	  "daily_social_media_hours": 3,
//	  "purchases_last_month": 5
//	}
func (app *application) addFavouriteHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	var rawAsset json.RawMessage
	if err := readJSON(w, r, &rawAsset); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	asset, err := models.UnmarshalAsset(rawAsset)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("invalid asset: %w", err))
		return
	}

	ctx := r.Context()

	if err := app.favouritesService.AddFavourite(ctx, userID, asset); err != nil {
		switch err {
		case store.ErrConflict:
			app.conflictResponse(w, r, err)
		default:
			app.badRequestResponse(w, r, err)
		}
		return
	}

	if err := app.writeEnvelopeJSON(w, http.StatusCreated, asset); err != nil {
		app.internalServerError(w, r, err)
	}
}

// getFavouriteHandler godoc
//
// @Summary      Get favourite
// @Description  Retrieves a specific favourited asset for a user.
// @Tags         favourites
// @Produce      json
// @Param        userID   path     string true "User ID" example(user-123)
// @Param        assetID  path     string true "Asset ID (external identifier)" example(asset-42)
// @Success      200 {object} EnvelopeResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /users/{userID}/favourites/{assetID} [get]
func (app *application) getFavouriteHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	assetID := chi.URLParam(r, "assetID")

	favourite, err := app.favouritesService.GetFavourite(r.Context(), userID, assetID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.writeEnvelopeJSON(w, http.StatusOK, favourite); err != nil {
		app.internalServerError(w, r, err)
	}
}

// removeFavouriteHandler godoc
//
// @Summary      Remove favourite
// @Description  Removes an asset from a user's favourites.
// @Tags         favourites
// @Param        userID   path     string true "User ID" example(user-123)
// @Param        assetID  path     string true "Asset ID (external identifier)" example(asset-42)
// @Success      204 "No Content"
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /users/{userID}/favourites/{assetID} [delete]
func (app *application) removeFavouriteHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	assetID := chi.URLParam(r, "assetID")

	if err := app.favouritesService.RemoveFavourite(r.Context(), userID, assetID); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdateDescriptionPayload struct {
	Description string `json:"description" validate:"required,max=1000"`
}

// updateDescriptionHandler godoc
//
// @Summary      Update favourite description
// @Description  Updates the description of a favourited asset.
// @Tags         favourites
// @Accept       json
// @Produce      json
// @Param        userID   path     string true "User ID" example(user-123)
// @Param        assetID  path     string true "Asset ID" example(asset-42)
// @Param        payload  body     UpdateDescriptionRequest true "New description payload"
// @Success      204 "No Content"
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /users/{userID}/favourites/{assetID}/description [put]
func (app *application) updateDescriptionHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	assetID := chi.URLParam(r, "assetID")

	var payload UpdateDescriptionPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.favouritesService.UpdateDescription(ctx, userID, assetID, payload.Description); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
