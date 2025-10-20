// Package handler contains HTTP handlers that translate incoming JSON requests
// into calls against a `store.Store`. Handlers are intentionally small and
// focused: decode request, call store, write response.
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mr-polymath/platform-go-challenge/internal/auth"
	"github.com/mr-polymath/platform-go-challenge/internal/db"
	"github.com/mr-polymath/platform-go-challenge/internal/model"
)

// Handler groups dependencies for HTTP handlers.
type Handler struct {
	FavouriteRepository db.FavouriteRepository
}

// New constructs a Handler with the provided store implementation.
func New(r db.FavouriteRepository) *Handler {
	return &Handler{FavouriteRepository: r}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {

	apiRoutes := r.Group("/users")
	secret := []byte("super-secret-key")

	r.POST("/token", GenerateTokenHandler(secret))

	apiRoutes.Use(AuthMiddleware(secret))
	{
		apiRoutes.GET("/:userId/favourites", h.GetAllFavourites)
		apiRoutes.POST("/:userId/favourites", h.AddFavourite)
		apiRoutes.DELETE("/:userId/favourites/:assetId", h.DeleteFavourite)
		apiRoutes.PATCH("/:userId/favourites/:assetId", h.PatchFavourite)
	}
}

// TokenRequest is the JSON payload to request a token
type TokenRequest struct {
	UserId string `json:"user_id" binding:"required,uuid"`
}

// TokenResponse is returned with the generated token
type TokenResponse struct {
	Token string `json:"token"`
}

// GenerateTokenHandler creates a JWT token for the provided user ID
func GenerateTokenHandler(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req TokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token, err := auth.GenerateToken(req.UserId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
			return
		}

		c.JSON(http.StatusOK, TokenResponse{Token: token})
	}
}

// ListFavourites returns all favourites for a given user id.
func (h *Handler) GetAllFavourites(c *gin.Context) {
	userIdFromToken := c.GetString("userId")
	if userIdFromToken != c.Param("userId") {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	userId, err := uuid.Parse(userIdFromToken)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	q := db.PagedQuery{
		Offset: parseIntOrDefault(c.Query("offset"), 0),
		Limit:  parseIntOrDefault(c.Query("limit"), 0),
		SortBy: c.Query("sortBy"),
		Order:  c.Query("order"),
	}

	assets, err := h.FavouriteRepository.GetAllFavourites(userId, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"favourites": assets})
}

// AddFavourite decodes the request and adds an asset to the user's
// favourites. It returns the created asset with 201 status.
func (h *Handler) AddFavourite(c *gin.Context) {
	userIdFromToken := c.GetString("userId")
	if userIdFromToken != c.Param("userId") {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	userId, err := uuid.Parse(userIdFromToken)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var asset model.Asset
	if err := c.BindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if err := h.FavouriteRepository.AddFavourite(userId, asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "favourite added",
		"favourite": asset})
}

// DeleteFavourite removes an asset from a user's favourites.
func (h *Handler) DeleteFavourite(c *gin.Context) {
	userIdFromToken := c.GetString("userId")
	if userIdFromToken != c.Param("userId") {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	userId, err := uuid.Parse(userIdFromToken)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	assetID, err := uuid.Parse(c.Param("assetId"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if err := h.FavouriteRepository.RemoveFavourite(userId, assetID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "asset not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "favourite removed",
		"favouriteId": assetID.String()})
}

// patchReq is the JSON payload used to update an asset's description.
type patchReq struct {
	Description string `json:"description"`
}

// PatchFavourite updates the description of a favourited asset.
func (h *Handler) PatchFavourite(c *gin.Context) {
	userIdFromToken := c.GetString("userId")
	if userIdFromToken != c.Param("userId") {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	userId, err := uuid.Parse(userIdFromToken)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	assetID, err := uuid.Parse(c.Param("assetId"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "asset not found"})
		return
	}

	var req patchReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if err := h.FavouriteRepository.PatchFavourite(userId, assetID, req.Description); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "asset not found"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "description updated"})
}
