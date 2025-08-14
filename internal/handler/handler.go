package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"gwi-favorites-service/internal/domain"
	"gwi-favorites-service/internal/service"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	favoritesService *service.FavoritesService
	logger           *logrus.Logger
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type UpdateDescriptionRequest struct {
	Description string `json:"description"`
}

func NewHandler(favoritesService *service.FavoritesService, logger *logrus.Logger) *Handler {
	return &Handler{
		favoritesService: favoritesService,
		logger:           logger,
	}
}

func (h *Handler) SetupRoutes() http.Handler {
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Apply middleware
	api.Use(h.LoggingMiddleware)
	api.Use(h.CORSMiddleware)

	// User favorites routes
	userRoutes := api.PathPrefix("/users/{userID}/favorites").Subrouter()
	userRoutes.HandleFunc("", h.GetUserFavorites).Methods("GET")
	userRoutes.HandleFunc("", h.AddFavorite).Methods("POST")
	userRoutes.HandleFunc("/{assetID}", h.RemoveFavorite).Methods("DELETE")
	userRoutes.HandleFunc("/{assetID}", h.UpdateFavoriteDescription).Methods("PUT")
	userRoutes.HandleFunc("/{assetID}/check", h.CheckIsFavorite).Methods("GET")

	// Health check
	r.HandleFunc("/health", h.HealthCheck).Methods("GET")

	return r
}

// GetUserFavorites handles GET /api/users/{userID}/favorites
func (h *Handler) GetUserFavorites(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	// Parse pagination parameters
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50 // Default limit
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	favorites, err := h.favoritesService.GetUserFavorites(r.Context(), userID, limit, offset)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.sendResponse(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    favorites,
	})
}

// AddFavorite handles POST /api/users/{userID}/favorites
func (h *Handler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	var rawAsset json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&rawAsset); err != nil {
		h.handleError(w, domain.ErrInvalidInput)
		return
	}

	asset, err := domain.AssetFromJSON(rawAsset)
	if err != nil {
		h.handleError(w, err)
		return
	}

	if err := h.favoritesService.AddFavorite(r.Context(), userID, asset); err != nil {
		h.handleError(w, err)
		return
	}

	h.sendResponse(w, http.StatusCreated, APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Asset added to favorites"},
	})
}

// RemoveFavorite handles DELETE /api/users/{userID}/favorites/{assetID}
func (h *Handler) RemoveFavorite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	assetID := vars["assetID"]

	if err := h.favoritesService.RemoveFavorite(r.Context(), userID, assetID); err != nil {
		h.handleError(w, err)
		return
	}

	h.sendResponse(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Asset removed from favorites"},
	})
}

// UpdateFavoriteDescription handles PUT /api/users/{userID}/favorites/{assetID}
func (h *Handler) UpdateFavoriteDescription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	assetID := vars["assetID"]

	var req UpdateDescriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, domain.ErrInvalidInput)
		return
	}

	if err := h.favoritesService.UpdateFavoriteDescription(r.Context(), userID, assetID, req.Description); err != nil {
		h.handleError(w, err)
		return
	}

	h.sendResponse(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Asset description updated"},
	})
}

// CheckIsFavorite handles GET /api/users/{userID}/favorites/{assetID}/check
func (h *Handler) CheckIsFavorite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	assetID := vars["assetID"]

	isFavorite, err := h.favoritesService.IsFavorite(r.Context(), userID, assetID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.sendResponse(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]bool{"is_favorite": isFavorite},
	})
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.sendResponse(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]string{
			"status":  "healthy",
			"service": "gwi-favorites-service",
		},
	})
}

// Helper methods
func (h *Handler) sendResponse(w http.ResponseWriter, statusCode int, response APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	var statusCode int
	var message string

	switch err {
	case domain.ErrUserNotFound:
		statusCode = http.StatusNotFound
		message = "User not found"
	case domain.ErrAssetNotFound:
		statusCode = http.StatusNotFound
		message = "Asset not found"
	case domain.ErrFavoriteNotFound:
		statusCode = http.StatusNotFound
		message = "Favorite not found"
	case domain.ErrFavoriteAlreadyExists:
		statusCode = http.StatusConflict
		message = "Asset is already in favorites"
	case domain.ErrInvalidInput, domain.ErrMissingRequiredField:
		statusCode = http.StatusBadRequest
		message = "Invalid input"
	case domain.ErrInvalidUserID:
		statusCode = http.StatusBadRequest
		message = "Invalid user ID"
	case domain.ErrInvalidAssetType:
		statusCode = http.StatusBadRequest
		message = "Invalid asset type"
	default:
		statusCode = http.StatusInternalServerError
		message = "Internal server error"
		h.logger.WithError(err).Error("Unexpected error occurred")
	}

	h.sendResponse(w, statusCode, APIResponse{
		Success: false,
		Error:   message,
	})
}

// Middleware
func (h *Handler) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom response writer to capture status code
		wrapped := &responseWriterWrapper{ResponseWriter: w}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		h.logger.WithFields(logrus.Fields{
			"method":     r.Method,
			"path":       r.URL.Path,
			"status":     wrapped.statusCode,
			"duration":   duration,
			"user_agent": r.UserAgent(),
		}).Info("HTTP request completed")
	})
}

func (h *Handler) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
