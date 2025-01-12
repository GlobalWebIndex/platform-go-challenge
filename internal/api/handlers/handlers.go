package handlers

// @title           Favorites Service API
// @version         1.0
// @description     API for managing user favorite assets including charts, insights, and audiences

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"platform-go-challenge/internal/api/middleware"
	"platform-go-challenge/internal/api/types"
	"platform-go-challenge/internal/domain"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/swaggo/swag"
)

type Handlers struct {
	assetService domain.AssetService
}

func New(fs domain.AssetService) *Handlers {
	return &Handlers{
		assetService: fs,
	}
}

func (h *Handlers) RegisterRoutes(r *mux.Router) {
	// Favorites endpoints
	r.HandleFunc("/favorites", h.ListFavorites).Methods(http.MethodGet)
	r.HandleFunc("/favorites/{assetId}", h.AddFavorite).Methods(http.MethodPost)
	r.HandleFunc("/favorites/{assetId}", h.RemoveFavorite).Methods(http.MethodDelete)
	r.HandleFunc("/favorites/{assetId}/description", h.UpdateDescription).Methods(http.MethodPatch)

	// swager
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The URL to find swagger.json
	))

	// redirect from /swagger to /swagger/
	r.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})
}

type UpdateDescriptionRequest struct {
	Description string `json:"description"`
}

// ListFavorites godoc
// @Summary      List user favorites
// @Description  Get list of user's favorite assets with pagination
// @Tags         favorites
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number"
// @Param        limit  query     int  false  "Items per page"
// @Success      200  {object}  types.APIResponse
// @Failure      401  {object}  types.APIResponse
// @Router       /favorites [get]
// @Security     Bearer
func (h *Handlers) ListFavorites(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		sendError(w, types.ErrCodeUnauthorized, "User not authenticated")
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	assets, err := h.assetService.ListFavorites(r.Context(), userID, page, limit)
	if err != nil {
		handleError(w, err)
		return
	}

	sendJSON(w, types.NewSuccessResponse(assets).WithMeta(&types.MetaData{
		RequestID: middleware.GetRequestID(r.Context()),
		Page:      page,
		PerPage:   limit,
	}))
}

// AddFavorite godoc
// @Summary      Add favorite
// @Description  Add an asset to user's favorites
// @Tags         favorites
// @Accept       json
// @Produce      json
// @Param        assetId   path     string  true  "Asset ID"
// @Success      201  {object}  types.APIResponse
// @Failure      400  {object}  types.APIResponse
// @Failure      401  {object}  types.APIResponse
// @Failure      409  {object}  types.APIResponse
// @Router       /favorites/{assetId} [post]
// @Security     Bearer
func (h *Handlers) AddFavorite(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		sendError(w, types.ErrCodeUnauthorized, "User not authenticated")
		return
	}

	assetID := mux.Vars(r)["assetId"]
	if assetID == "" {
		sendError(w, types.ErrCodeValidation, "Asset ID is required")
		return
	}

	err := h.assetService.AddFavorite(r.Context(), userID, assetID)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	sendJSON(w, types.NewSuccessResponse(nil).WithMeta(&types.MetaData{
		RequestID: middleware.GetRequestID(r.Context()),
	}))
}

// RemoveFavorite godoc
// @Summary      Remove favorite
// @Description  Remove an asset from user's favorites
// @Tags         favorites
// @Accept       json
// @Produce      json
// @Param        assetId   path     string  true  "Asset ID"
// @Success      200  {object}  types.APIResponse
// @Failure      400  {object}  types.APIResponse
// @Failure      401  {object}  types.APIResponse
// @Failure      404  {object}  types.APIResponse
// @Router       /favorites/{assetId} [delete]
// @Security     Bearer
func (h *Handlers) RemoveFavorite(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		sendError(w, types.ErrCodeUnauthorized, "User not authenticated")
		return
	}

	assetID := mux.Vars(r)["assetId"]
	if assetID == "" {
		sendError(w, types.ErrCodeValidation, "Asset ID is required")
		return
	}

	err := h.assetService.RemoveFavorite(r.Context(), userID, assetID)
	if err != nil {
		handleError(w, err)
		return
	}

	sendJSON(w, types.NewSuccessResponse(nil).WithMeta(&types.MetaData{
		RequestID: middleware.GetRequestID(r.Context())}))
}

// UpdateDescription godoc
// @Summary      Update description
// @Description  Update description of a favorite asset
// @Tags         favorites
// @Accept       json
// @Produce      json
// @Param        assetId   path     string  true  "Asset ID"
// @Param        request   body     UpdateDescriptionRequest  true  "Update description request"
// @Success      200  {object}  types.APIResponse
// @Failure      400  {object}  types.APIResponse
// @Failure      401  {object}  types.APIResponse
// @Failure      404  {object}  types.APIResponse
// @Router       /favorites/{assetId}/description [patch]
// @Security     Bearer
func (h *Handlers) UpdateDescription(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		sendError(w, types.ErrCodeUnauthorized, "User not authenticated")
		return
	}

	assetID := mux.Vars(r)["assetId"]
	if assetID == "" {
		sendError(w, types.ErrCodeValidation, "Asset ID is required")
		return
	}

	var req UpdateDescriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, types.ErrCodeValidation, "Invalid request body")
		return
	}

	err := h.assetService.UpdateDescription(r.Context(), userID, assetID, req.Description)
	if err != nil {
		handleError(w, err)
		return
	}

	sendJSON(w, types.NewSuccessResponse(nil))
}

// Helper functions for consistent error handling and responses
func handleError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrNotFound:
		sendError(w, types.ErrCodeNotFound, "Resource not found")
	case domain.ErrInvalidAssetType:
		sendError(w, types.ErrCodeInvalidAssetType, "Invalid asset type")
	case domain.ErrAlreadyExists:
		sendError(w, types.ErrCodeAlreadyFavorited, "Asset already favorited")
	case context.DeadlineExceeded:
		sendError(w, types.ErrCodeAlreadyFavorited, "Asset already favorited")
	default:
		sendError(w, types.ErrCodeInternal, "Internal server error")
	}
}

func sendError(w http.ResponseWriter, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(getHTTPStatus(code))
	json.NewEncoder(w).Encode(types.NewErrorResponse(code, message))
}

func sendJSON(w http.ResponseWriter, response types.APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getHTTPStatus(code string) int {
	switch code {
	case types.ErrCodeValidation:
		return http.StatusBadRequest
	case types.ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case types.ErrCodeNotFound:
		return http.StatusNotFound
	case types.ErrTimedOut:
		return http.StatusGatewayTimeout
	default:
		return http.StatusInternalServerError
	}
}
