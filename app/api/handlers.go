package api

import (
	"github.com/gin-gonic/gin"

	"app/api/schema"
	"app/api/services"

	"net/http"
)

type ApiHandler struct {
	service services.XService
}

func NewApiHandler(service services.XService) *ApiHandler {
	return &ApiHandler{service: service}
}

func (h *ApiHandler) SetupRoutes(router *gin.Engine) {
	router.GET("/assets", h.GetAssets)
	router.PATCH("/assets/:assetId", h.UpdateAsset)
	router.GET("/users/:userId/favorites", h.GetUserFavorites)
	router.POST("/users/:userId/favorites", h.CreateUserFavorite)
	router.DELETE("/users/:userId/favorites/:favoriteId", h.DeleteUserFavorite)
}

func (h *ApiHandler) GetAssets(c *gin.Context) {
	assets := h.service.GetAssets(schema.AssetFilter{})
	c.JSON(http.StatusOK, assets)
}

// TODO also check if user exists, because a user will be firing this request. or make compelling arguments about that being taken care of the security framework
func (h *ApiHandler) UpdateAsset(c *gin.Context) {
	assetId := c.Param("assetId")
	var json schema.AssetUpdateRequest
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.service.UpdateAsset(assetId, json)
	if err != nil {
		httpError, ok := err.(schema.HttpError)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(httpError.Status(), gin.H{"error": httpError.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{})
	}
}

func (h *ApiHandler) GetUserFavorites(c *gin.Context) {
	userId := c.Param("userId")
	userFavorites := h.service.GetUserFavorites(userId)
	c.JSON(http.StatusOK, userFavorites)
}

func (h *ApiHandler) CreateUserFavorite(c *gin.Context) {
	userId := c.Param("userId")
	var json schema.UserFavoriteRequest
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.service.CreateUserFavorite(userId, json.AssetID)
	if err != nil {
		httpError, ok := err.(schema.HttpError)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(httpError.Status(), gin.H{"error": httpError.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"id": id}) // not good, must return schema.UserFavoriteDto
	}
}

func (h *ApiHandler) DeleteUserFavorite(c *gin.Context) {
	userId := c.Param("userId")
	userFavoriteId := c.Param("favoriteId")
	err := h.service.DeleteUserFavorite(userId, userFavoriteId)
	if err != nil {
		httpError, ok := err.(schema.HttpError)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(httpError.Status(), gin.H{"error": httpError.Error()})
		}
	}
}
