package handlers

import (
	"log"
	"net/http"
	"platform-go-challenge/internal/errors"
	"platform-go-challenge/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FavouriteHandler provides methods for handling favourite-related requests.
type FavouriteHandler struct {
	FavouriteService *service.FavouriteService
}

// NewFavouriteHandler creates and returns a new FavouriteHandler.
func NewFavouriteHandler(s *service.FavouriteService) *FavouriteHandler {
	return &FavouriteHandler{
		FavouriteService: s,
	}
}

// @Summary Get a user's favourites
// @Description Retrieves a list of all favourite assets for a given user.
// @Tags favourites
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer"
// @Success 200 {array} entities.Asset
// @Failure 500 {object} gin.H
// @Router /favourites [get]
func (h *FavouriteHandler) GetFavourites(c *gin.Context) {
	userID := c.Param("userID")

	page, size, err := ExtractPageAndSize(c)
	if err != nil {
		return
	}

	if userID == "" {
		c.Error(errors.NewCustomError(http.StatusBadRequest, "user_id", "userID is required"))
		return
	}
	ctx := c.Request.Context()
	// Call the service layer to get the data.
	favourites := h.FavouriteService.GetFavourites(ctx, userID, page, size)

	c.JSON(http.StatusOK, favourites)
}

// @Summary Add an asset to a user's favourites
// @Description Adds a new asset to a user's favourites list.
// @Tags favourites
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer"
// @Param request body requests.AddFavouriteRequest true "Add Favourite Request"
// @Success 201 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /favourites [post]
func (h *FavouriteHandler) AddFavourite(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID is required"})
		return
	}

	var reqBody struct {
		AssetID string `json:"assetID" binding:"required"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	// Call the service layer to perform the business logic.
	err := h.FavouriteService.AddFavourite(ctx, userID, reqBody.AssetID)
	if err != nil {
		log.Printf("Error adding favourite for user %s, asset %s: %v", userID, reqBody.AssetID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add favourite"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "favourite added successfully"})
}

// @Summary Update a favourite's description
// @Description Updates the description of a specific favourite asset.
// @Tags favourites
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer"
// @Param assetID path string true "Asset ID"
// @Param request body requests.UpdateFavouriteRequest true "Update Description"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /favourites/{assetID} [patch]
func (h *FavouriteHandler) UpdateFavouriteDescription(c *gin.Context) {
	userID := c.Param("userID")
	assetID := c.Param("assetID")

	var reqBody struct {
		Description string `json:"description" binding:"required"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	err := h.FavouriteService.UpdateFavouriteDescription(ctx, userID, assetID, reqBody.Description)
	if err != nil {
		log.Printf("Error updating favourite for user %s, asset %s: %v", userID, assetID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update favourite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "favourite description updated successfully"})
}

// @Summary Remove a favourite asset
// @Description Removes a specific asset from a user's favourites list.
// @Tags favourites
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer"
// @Param assetID path string true "Asset ID"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /favourites/{assetID} [delete]
func (h *FavouriteHandler) DeleteFavourite(c *gin.Context) {
	userID := c.Param("userID")
	assetID := c.Param("assetID")

	ctx := c.Request.Context()
	err := h.FavouriteService.RemoveFavourite(ctx, userID, assetID)
	if err != nil {
		log.Printf("Error deleting favourite for user %s, asset %s: %v", userID, assetID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete favourite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "favourite removed successfully"})
}

func ExtractPageAndSize(c *gin.Context) (page, size int, err error) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	page, err = strconv.Atoi(pageStr)
	if err != nil {
		return 0, 0, err
	}

	size, err = strconv.Atoi(sizeStr)
	if err != nil {
		return 0, 0, err
	}

	return page, size, nil
}
