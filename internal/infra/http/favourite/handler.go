package favourite

import (
	"encoding/json"
	"fmt"
	queries2 "github.com/akazantzidis/gwi-ass/internal/app/user/queries"
	"github.com/akazantzidis/gwi-ass/internal/pkg/helper"
	"github.com/akazantzidis/gwi-ass/internal/pkg/middleware"
	"net/http"

	"github.com/akazantzidis/gwi-ass/internal/app"
	"github.com/akazantzidis/gwi-ass/internal/app/favourite/commands"
	"github.com/akazantzidis/gwi-ass/internal/app/favourite/queries"
	"github.com/akazantzidis/gwi-ass/internal/domain/favourite"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Handler Favorite http request Handler
type Handler struct {
	favoriteServices app.FavoriteServices
	userServices     app.UserServices
}

func NewHandler(app app.FavoriteServices, userApp app.UserServices) *Handler {
	return &Handler{favoriteServices: app, userServices: userApp}
}

// URL param constants
const (
	UserIDURLParam           = "userID"
	GetFavoriteIDURLParam    = "favoriteId"
	UpdateFavoriteIDURLParam = "favoriteId"
	DeleteFavoriteIDURLParam = "favoriteId"
)

// GetAll returns all favorites for a given user
func (c Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := extractUserID(w, r)
	if !ok {
		return
	}

	favorites, err := c.favoriteServices.Queries.GetAllFavoritesHandler.Handle(
		queries.GetAllFavoritesRequest{UserID: userID},
	)
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err, nil)
		return
	}

	if favorites == nil {
		favorites = []queries.GetAllFavoritesResult{}
	}

	json.NewEncoder(w).Encode(favorites)
}

// GetByID returns a favorite for a given user
func (c Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := extractUserID(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	favoriteID, err := uuid.Parse(vars[GetFavoriteIDURLParam])
	if err != nil {
		helper.WriteJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid favorite ID"), nil)
		return
	}

	fav, err := c.favoriteServices.Queries.GetFavoriteHandler.Handle(
		queries.GetFavoriteRequest{
			UserID:     userID,
			FavoriteID: favoriteID,
		},
	)

	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err, nil)
		return
	}
	if fav == nil {
		helper.WriteJSONError(w, http.StatusNotFound, fmt.Errorf("favorite not found"), nil)
		return
	}

	json.NewEncoder(w).Encode(fav)
}

// CreateFavoriteRequestModel represents the request model expected for Add request
type CreateFavoriteRequestModel struct {
	Type        favourite.AssetType `json:"type"`
	Description string              `json:"description"`
	Data        json.RawMessage     `json:"data"`
}

// Create adds a new favorite for a user
func (c Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := extractUserID(w, r)
	if !ok {
		return
	}

	var req CreateFavoriteRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid request body"), nil)
		return
	}

	if !req.Type.IsValid() {
		helper.WriteJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid asset type"), nil)
		return
	}

	// Ensure user exists
	if _, err := c.userServices.Queries.GetUserHandler.Handle(
		queries2.GetUserRequest{ID: userID},
	); err != nil {
		helper.WriteJSONError(w, http.StatusNotFound, err, nil)
		return
	}

	newFavID := uuid.New()

	err := c.favoriteServices.Commands.CreateFavoriteHandler.Handle(
		commands.AddFavoriteRequest{
			ID:          newFavID,
			UserID:      userID,
			Type:        req.Type,
			Description: req.Description,
			Data:        req.Data,
		},
	)
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err, nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": newFavID.String()})
}

// UpdateFavoriteRequestModel represents the request model of Update
type UpdateFavoriteRequestModel struct {
	ID          uuid.UUID           `json:"id"`
	Type        favourite.AssetType `json:"type"`
	Description string              `json:"description"`
	Data        json.RawMessage     `json:"data"`
}

// PatchFavoriteRequestModel represents the fields that can be partially updated
type PatchFavoriteRequestModel struct {
	Type        *favourite.AssetType `json:"type,omitempty"`
	Description *string              `json:"description,omitempty"`
	Data        *json.RawMessage     `json:"data,omitempty"`
}

// Patch handles partial updates for favorites
func (c Handler) Patch(w http.ResponseWriter, r *http.Request) {
	userID, ok := extractUserID(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	favID, err := uuid.Parse(vars[UpdateFavoriteIDURLParam])
	if err != nil {
		helper.WriteJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid favorite ID"), nil)
		return
	}

	var req PatchFavoriteRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid request body"), nil)
		return
	}

	// Validate if type provided
	if req.Type != nil && !req.Type.IsValid() {
		helper.WriteJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid asset type"), nil)
		return
	}

	result, err := c.favoriteServices.Commands.UpdatePartialFavoriteHandler.HandlePartial(
		userID, favID,
		commands.PatchFavoriteRequest{
			Type:        req.Type,
			Description: req.Description,
			Data:        req.Data,
		},
	)
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err, nil)
		return
	}

	json.NewEncoder(w).Encode(result)
}

// Update handles full updates
func (c Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := extractUserID(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	favoriteID, err := uuid.Parse(vars[UpdateFavoriteIDURLParam])
	if err != nil {
		helper.WriteJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid favorite ID"), nil)
		return
	}

	var req UpdateFavoriteRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid request body"), nil)
		return
	}

	// Validation
	if !req.Type.IsValid() {
		helper.WriteJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid asset type"), nil)
		return
	}

	err = c.favoriteServices.Commands.UpdateFavoriteHandler.Handle(
		commands.UpdateFavoriteRequest{
			UserID:      userID,
			ID:          favoriteID,
			Type:        req.Type,
			Description: req.Description,
			Data:        req.Data,
		},
	)
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err, nil)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Delete deletes a favorite
func (c Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := extractUserID(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	favoriteID, err := uuid.Parse(vars[DeleteFavoriteIDURLParam])
	if err != nil {
		helper.WriteJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid favorite ID"), nil)
		return
	}

	err = c.favoriteServices.Commands.DeleteFavoriteHandler.Handle(
		commands.DeleteFavoriteRequest{
			UserID:     userID,
			FavoriteID: favoriteID,
		},
	)
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err, nil)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// extractUserID checks both context and URL param, ensuring they match.
func extractUserID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	// 1. Read from JWT context
	ctxVal := r.Context().Value(middleware.ContextUserKey)
	if ctxVal == nil {
		helper.WriteJSONError(w, http.StatusUnauthorized, fmt.Errorf("missing user in token"), nil)
		return uuid.Nil, false
	}

	userIDStr, ok := ctxVal.(string) // JWT middleware should store as string
	if !ok {
		helper.WriteJSONError(w, http.StatusUnauthorized, fmt.Errorf("invalid user ID in token"), nil)
		return uuid.Nil, false
	}

	ctxUserID, err := uuid.Parse(userIDStr)
	if err != nil {
		helper.WriteJSONError(w, http.StatusUnauthorized, fmt.Errorf("invalid user ID format in token"), nil)
		return uuid.Nil, false
	}

	// 2. Read from URL param
	vars := mux.Vars(r)
	paramID, err := uuid.Parse(vars[UserIDURLParam])
	if err != nil {
		helper.WriteJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID parameter"), nil)
		return uuid.Nil, false
	}

	// 3. Compare — forbid accessing other users
	if ctxUserID != paramID {
		helper.WriteJSONError(w, http.StatusForbidden, fmt.Errorf("forbidden: cannot access another user's data"), nil)
		return uuid.Nil, false
	}

	return ctxUserID, true
}
