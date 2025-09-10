package handlers

import (
	"net/http"
	"platform-go-challenge/internal/domain/dto"
	"platform-go-challenge/internal/middleware"
	"platform-go-challenge/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FavouriteHandler provides methods for handling favourite-related requests.
type FavouriteHandler struct {
	FavouriteService service.FavouriteServiceIntf
}

// NewFavouriteHandler creates and returns a new FavouriteHandler.
func NewFavouriteHandler(s *service.FavouriteService) *FavouriteHandler {
	return &FavouriteHandler{
		FavouriteService: s,
	}
}

// isAuthorized checks if the user from the JWT has permission to access the resource.
// A user is authorized if their ID matches the resource's user ID or if they are an admin.
func isAuthorized(c *gin.Context, userIDInPath string) bool {
	// Get the user ID and role from the context.
	userIDFromToken, _ := c.Get(middleware.JWT_USER_ID_KEY)
	userRole, _ := c.Get(middleware.JWT_ROLE_KEY)

	// Check if the user is an admin.
	if userRole == "admin" {
		return true
	}

	// Check if the user ID from the token matches the user ID in the path.
	if userIDFromToken == userIDInPath {
		return true
	}

	return false
}

// @Summary Get a user's favourites
// @Description Retrieves a paginated list of all favourite assets for a given user.
// @Tags favourites
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param userID path string true "User ID"
// @Param page query int false "Page number (default: 1)"
// @Param size query int false "Page size (default: 10)"
// @Success 200 {array} entities.Asset
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{userID}/favourites [get]
func (h *FavouriteHandler) GetFavourites(c *gin.Context) {
	userID := c.Param("userID")

	// Authorization check.
	if !isAuthorized(c, userID) {
		c.AbortWithError(http.StatusForbidden, middleware.NewAPIError("access denied", ""))
		return
	}

	page, size, err := ExtractPageAndSize(c)
	if err != nil {
		return
	}

	if userID == "" {
		c.AbortWithError(http.StatusBadRequest, middleware.NewAPIError("userID is required", "userID"))
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
// @Param Authorization header string true "Bearer token"
// @Param userID path string true "User ID"
// @Param request body dto.AddFavouriteRequest true "Favourite request payload"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{userID}/favourites [post]
func (h *FavouriteHandler) AddFavourite(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.AbortWithError(http.StatusBadRequest, middleware.NewAPIError("userID is required", "userID"))
		return
	}

	// Authorization check.
	if !isAuthorized(c, userID) {
		c.AbortWithError(http.StatusForbidden, middleware.NewAPIError("access denied", ""))
		return
	}

	var reqBody dto.AddFavouriteRequest

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithError(http.StatusBadRequest, middleware.NewAPIError("invalid json body", ""))
		return
	}

	ctx := c.Request.Context()
	// Call the service layer to perform the business logic.
	err := h.FavouriteService.AddFavourite(ctx, userID, reqBody.AssetID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, middleware.NewAPIError(err.Error(), ""))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "favourite added successfully"})
}

// @Summary Update a favourite's description
// @Description Updates the description of a specific favourite asset.
// @Tags favourites
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param userID path string true "User ID"
// @Param assetID path string true "Asset ID"
// @Param request body dto.UpdateFavouriteDescriptionRequest true "Update description payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{userID}/favourites/{assetID} [patch]
func (h *FavouriteHandler) UpdateFavouriteDescription(c *gin.Context) {
	userID := c.Param("userID")
	assetID := c.Param("assetID")

	// Authorization check.
	if !isAuthorized(c, userID) {
		c.AbortWithError(http.StatusForbidden, middleware.NewAPIError("access denied", ""))
		return
	}

	var reqBody dto.UpdateFavouriteDescriptionRequest

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithError(http.StatusBadRequest, middleware.NewAPIError("invalid json body", ""))
		return
	}

	ctx := c.Request.Context()
	err := h.FavouriteService.UpdateFavouriteDescription(ctx, userID, assetID, reqBody.GetDescription())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, middleware.NewAPIError(err.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "favourite description updated successfully"})
}

// @Summary Remove a favourite asset
// @Description Removes a specific asset from a user's favourites list.
// @Tags favourites
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param userID path string true "User ID"
// @Param assetID path string true "Asset ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{userID}/favourites/{assetID} [delete]
func (h *FavouriteHandler) DeleteFavourite(c *gin.Context) {
	userID := c.Param("userID")
	assetID := c.Param("assetID")

	// Authorization check.
	if !isAuthorized(c, userID) {
		c.AbortWithError(http.StatusForbidden, middleware.NewAPIError("access denied", ""))
		return
	}

	if userID == "" {
		c.AbortWithError(http.StatusBadRequest, middleware.NewAPIError("userID is required", "userID"))
		return
	}

	ctx := c.Request.Context()
	err := h.FavouriteService.RemoveFavourite(ctx, userID, assetID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, middleware.NewAPIError(err.Error(), ""))
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
