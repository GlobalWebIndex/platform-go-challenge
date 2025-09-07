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

// GetFavourites handles the GET request to retrieve a user's favourites.
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

// AddFavourite handles the POST request to add an asset to a user's favourites.
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

// UpdateFavouriteDescription handles the PUT request to edit a favourite's description.
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

// DeleteFavourite handles the DELETE request to remove a favourite.
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
