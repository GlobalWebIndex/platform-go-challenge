package handlers

import (
	"errors"
	"net/http"
	"platform-go-challenge/internal/domain/dto"
	"platform-go-challenge/internal/middleware"
	"platform-go-challenge/internal/repository"
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
		c.AbortWithStatusJSON(http.StatusForbidden, middleware.NewAPIError(http.StatusForbidden, "access denied", ""))
		return
	}

	page, size, err := ExtractPageAndSize(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, middleware.NewAPIError(http.StatusForbidden, "invalid paging parameters", ""))
		return
	}

	if userID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, middleware.NewAPIError(http.StatusBadRequest, "userID is required", "userID"))
		return
	}

	favourites, err := h.FavouriteService.GetFavourites(userID, page, size)
	if err != nil {
		HandleServiceErrors(c, err)
		return
	}

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
		c.AbortWithStatusJSON(http.StatusBadRequest, middleware.NewAPIError(http.StatusBadRequest, "userID is required", "userID"))
		return
	}

	// Authorization check.
	if !isAuthorized(c, userID) {
		c.AbortWithStatusJSON(http.StatusForbidden, middleware.NewAPIError(http.StatusForbidden, "access denied", ""))
		return
	}

	var reqBody dto.AddFavouriteRequest

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, middleware.NewAPIError(http.StatusBadRequest, "invalid json body", ""))
		return
	}

	// Call the service layer to perform the business logic.
	err := h.FavouriteService.AddFavourite(userID, reqBody.AssetID)
	if err != nil {
		HandleServiceErrors(c, err)
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
		c.AbortWithStatusJSON(http.StatusForbidden, middleware.NewAPIError(http.StatusForbidden, "access denied", ""))
		return
	}

	var reqBody dto.UpdateFavouriteDescriptionRequest

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, middleware.NewAPIError(http.StatusBadRequest, "invalid json body", ""))
		return
	}

	err := h.FavouriteService.UpdateFavouriteDescription(userID, assetID, reqBody.GetDescription())
	if err != nil {
		HandleServiceErrors(c, err)
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
		c.AbortWithStatusJSON(http.StatusForbidden, middleware.NewAPIError(http.StatusForbidden, "access denied", ""))
		return
	}

	if userID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, middleware.NewAPIError(http.StatusBadRequest, "userID is required", "userID"))
		return
	}

	err := h.FavouriteService.RemoveFavourite(userID, assetID)
	if err != nil {
		HandleServiceErrors(c, err)
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

func HandleServiceErrors(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.UserNotFoundError):
		c.AbortWithStatusJSON(http.StatusNotFound, middleware.NewAPIError(http.StatusNotFound, err.Error(), "userID"))
	case errors.Is(err, repository.AssetNotFoundError):
		c.AbortWithStatusJSON(http.StatusNotFound, middleware.NewAPIError(http.StatusNotFound, err.Error(), "assetID"))
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, middleware.NewAPIError(http.StatusInternalServerError, err.Error(), ""))
	}
	c.Abort()
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
