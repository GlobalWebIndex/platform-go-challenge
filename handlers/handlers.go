package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"platform-go-challenge/models"
	"platform-go-challenge/storage"
	"sync"
)

// GetFavoritesHandler retrieves a user's favorites concurrently
func GetFavoritesHandler(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}
	// Check if user exists in storage
	if !storage.UserExists(userID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var allFavorites []models.AssetInterface
	getAndAppend := func(assetType string) {
		defer wg.Done()
		assets := storage.GetFavoritesByType(userID, assetType)
		mu.Lock()
		allFavorites = append(allFavorites, assets...)
		mu.Unlock()
	}
	wg.Add(3)
	go getAndAppend("chart")
	go getAndAppend("insight")
	go getAndAppend("audience")
	wg.Wait()
	c.JSON(http.StatusOK, gin.H{"favorites": allFavorites})
}

// Add a favorite
func AddFavoriteHandler(c *gin.Context) {
	var req struct {
		UserID string          `json:"user_id"`
		Asset  json.RawMessage `json:"asset"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	var base models.Asset
	if err := json.Unmarshal(req.Asset, &base); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset format"})
		return
	}
	// Check if user exists in storage
	if !storage.UserExists(req.UserID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	// Check if asset with same ID already exists
	if storage.AssetExists(req.UserID, base.ID) {
		c.JSON(http.StatusConflict, gin.H{"error": "Asset with this ID already exists"})
		return
	}
	switch base.Type {
	case models.ChartType:
		var chart models.Chart
		if err := json.Unmarshal(req.Asset, &chart); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chart format"})
			return
		}
		storage.AddFavorite(req.UserID, chart)
	case models.AudienceType:
		var audience models.Audience
		if err := json.Unmarshal(req.Asset, &audience); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audience format"})
			return
		}
		storage.AddFavorite(req.UserID, audience)
	case models.InsightType:
		var insight models.Insight
		if err := json.Unmarshal(req.Asset, &insight); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid insight format"})
			return
		}
		storage.AddFavorite(req.UserID, insight)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset type"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Asset added to favorites"})
}

// Remove a favorite
func RemoveFavoriteHandler(c *gin.Context) {
	userID := c.Param("user_id")
	assetID := c.Param("asset_id")
	if userID == "" || assetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID and Asset ID are required"})
		return
	}
	// Check if user exists in storage
	if !storage.UserExists(userID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	// Check if asset exists before removing
	if !storage.AssetExists(userID, assetID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}
	storage.RemoveFavorite(userID, assetID)
	c.JSON(http.StatusOK, gin.H{"message": "Asset removed"})
}

// Edit favorite description
func EditFavoriteHandler(c *gin.Context) {
	userID := c.Param("user_id")
	assetID := c.Param("asset_id")
	if userID == "" || assetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID and Asset ID are required"})
		return
	}
	var req struct {
		NewDescription string `json:"new_description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	if req.NewDescription == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New description cannot be empty"})
		return
	}
	// Check if user exists in storage
	if !storage.UserExists(userID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	// Check if asset exists
	if !storage.AssetExists(userID, assetID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}
	storage.EditFavorite(userID, assetID, req.NewDescription)
	c.JSON(http.StatusOK, gin.H{"message": "Asset updated"})
}
