package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/GlobalWebIndex/platform-go-challenge/internal/api/v1/dto"
	"github.com/GlobalWebIndex/platform-go-challenge/internal/domain/favourites"
	"github.com/GlobalWebIndex/platform-go-challenge/internal/shared/httpx"
)

// FavouritesHandler exposes HTTP endpoints for managing user favourites.
type FavouritesHandler struct {
	svc favourites.Service
}

// NewFavouritesHandler constructs a new FavouritesHandler.
func NewFavouritesHandler(svc favourites.Service) *FavouritesHandler {
	if svc == nil {
		panic("favourites handler: nil service")
	}
	return &FavouritesHandler{svc: svc}
}

// RegisterRoutes registers favourites endpoints under the given router.
// Example base path: /v1/users/{userId}/favourites
func (h *FavouritesHandler) RegisterRoutes(r chi.Router) {
	r.Get("/users/{userId}/favourites", h.ListUserFavourites)
	r.Post("/users/{userId}/favourites", h.AddFavourite)
	r.Delete("/users/{userId}/favourites/{assetId}", h.RemoveFavourite)
	r.Patch("/users/{userId}/favourites/{assetId}", h.EditFavouriteDescription)
}

// ListUserFavourites handles GET /v1/users/{userId}/favourites?type=chart
func (h *FavouritesHandler) ListUserFavourites(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := chi.URLParam(r, "userId")
	if userID == "" {
		httpx.WriteProblem(w, http.StatusBadRequest, "invalid_user", "userId is required", "", r.URL.Path)
		return
	}

	filter := favourites.ListFilter{}
	if t := r.URL.Query().Get("type"); t != "" {
		at, err := parseAssetType(t)
		if err != nil {
			httpx.WriteProblem(
				w,
				http.StatusBadRequest,
				"invalid_type",
				"invalid asset type",
				"type must be chart, insight, or audience",
				r.URL.Path,
			)
			return
		}
		filter.Type = &at
	}

	entries, err := h.svc.ListUserFavourites(ctx, userID, filter)
	if err != nil {
		status, typ, title, detail := mapServiceError(err)
		httpx.WriteProblem(w, status, typ, title, detail, r.URL.Path)
		return
	}

	items := make([]dto.FavouriteItemResponse, 0, len(entries))
	for _, e := range entries {
		items = append(items, favouriteEntryToListItemDTO(e))
	}

	resp := dto.ListFavouritesResponse{
		Items: items,
		Meta: dto.ListMetaResponse{
			Total: len(entries),
			Count: len(entries),
		},
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}

// AddFavourite handles POST /v1/users/{userId}/favourites
func (h *FavouritesHandler) AddFavourite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := chi.URLParam(r, "userId")
	if userID == "" {
		httpx.WriteProblem(w, http.StatusBadRequest, "invalid_user", "userId is required", "", r.URL.Path)
		return
	}

	var req dto.CreateFavouriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "invalid_json", "invalid JSON payload", err.Error(), r.URL.Path)
		return
	}

	domainAsset, err := dtoToDomainAsset(req)
	if err != nil {
		httpx.WriteProblem(w, http.StatusUnprocessableEntity, "invalid_asset", "invalid asset payload", err.Error(), r.URL.Path)
		return
	}

	entry, err := h.svc.AddFavourite(ctx, userID, domainAsset)
	if err != nil {
		status, typ, title, detail := mapServiceError(err)
		httpx.WriteProblem(w, status, typ, title, detail, r.URL.Path)
		return
	}

	resp := favouriteEntryToCreateResponse(entry)
	httpx.WriteJSON(w, http.StatusCreated, resp)
}

// RemoveFavourite handles DELETE /v1/users/{userId}/favourites/{assetId}
func (h *FavouritesHandler) RemoveFavourite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := chi.URLParam(r, "userId")
	assetID := chi.URLParam(r, "assetId")

	if userID == "" || assetID == "" {
		httpx.WriteProblem(w, http.StatusBadRequest, "invalid_path", "userId and assetId are required", "", r.URL.Path)
		return
	}

	if err := h.svc.RemoveFavourite(ctx, userID, assetID); err != nil {
		status, typ, title, detail := mapServiceError(err)
		httpx.WriteProblem(w, status, typ, title, detail, r.URL.Path)
		return
	}

	// Idempotent delete: always respond 204 even if favourite wasn't present.
	w.WriteHeader(http.StatusNoContent)
}

// EditFavouriteDescription handles PATCH /v1/users/{userId}/favourites/{assetId}
func (h *FavouritesHandler) EditFavouriteDescription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := chi.URLParam(r, "userId")
	assetID := chi.URLParam(r, "assetId")

	if userID == "" || assetID == "" {
		httpx.WriteProblem(w, http.StatusBadRequest, "invalid_path", "userId and assetId are required", "", r.URL.Path)
		return
	}

	var req dto.EditFavouriteDescriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "invalid_json", "invalid JSON payload", err.Error(), r.URL.Path)
		return
	}
	if req.Description == "" {
		httpx.WriteProblem(w, http.StatusUnprocessableEntity, "invalid_description", "description is required", "", r.URL.Path)
		return
	}

	entry, err := h.svc.EditFavouriteDescription(ctx, userID, assetID, req.Description)
	if err != nil {
		status, typ, title, detail := mapServiceError(err)
		httpx.WriteProblem(w, status, typ, title, detail, r.URL.Path)
		return
	}

	resp := favouriteEntryToListItemDTO(entry)
	httpx.WriteJSON(w, http.StatusOK, resp)
}
