package http

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"platform-go-challenge/src/core"
)

type Handlers struct {
	assetSvc    core.AssetService
	favoriteSvc core.FavoriteService
	now         func() time.Time
}

func NewHandlers(assetSvc core.AssetService, favSvc core.FavoriteService, now func() time.Time) *Handlers {
	return &Handlers{assetSvc: assetSvc, favoriteSvc: favSvc, now: now}
}

// GET /v1/users/:user_id/favorites?limit=&after_t=&after_id=&type=
func (h *Handlers) getFavorites(c *gin.Context) {
	userID := core.UserID(c.Param("user_id"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var typ *core.AssetType
	if s := c.Query("type"); s != "" {
		at := core.AssetType(s)
		typ = &at
	}

	var after *core.FavoriteEdge
	tStr := c.Query("after_t")
	idStr := c.Query("after_id")
	if tStr != "" || idStr != "" {
		if tStr == "" || strings.TrimSpace(idStr) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "both after_t and after_id are required"})
			return
		}
		ns, err := strconv.ParseInt(tStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid after_t"})
			return
		}
		after = &core.FavoriteEdge{AssetID: core.AssetID(idStr), CreatedAt: time.Unix(0, ns).UTC()}
	}

	items, next, err := h.favoriteSvc.ListFavorites(userID, limit, after, typ)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list favorites"})
		return
	}

	var nextAfterT, nextAfterID *string
	if next != nil {
		ns := strconv.FormatInt(next.CreatedAt.UTC().UnixNano(), 10)
		id := string(next.AssetID)
		nextAfterT, nextAfterID = &ns, &id
	}

	c.JSON(http.StatusOK, gin.H{
		"items":         items,        // concrete entities
		"next_after_t":  nextAfterT,   // string | null
		"next_after_id": nextAfterID,  // string | null
	})
}

// POST /v1/users/:user_id/favorites   body: { "assetId": "..." }
func (h *Handlers) addFavorite(c *gin.Context) {
	userID := core.UserID(c.Param("user_id"))
	var body map[string]any
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	assetID, _ := body["assetId"].(string)
	if strings.TrimSpace(assetID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "assetId is required"})
		return
	}
	if err := h.favoriteSvc.AddToFavorites(userID, core.AssetID(assetID), h.now()); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// DELETE /v1/users/:user_id/favorites/:asset_id
func (h *Handlers) removeFavorite(c *gin.Context) {
	userID := core.UserID(c.Param("user_id"))
	assetID := core.AssetID(c.Param("asset_id"))
	if err := h.favoriteSvc.RemoveFromFavorites(userID, assetID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// PATCH /v1/users/:user_id/favorites/:asset_id   body: { "description": "...", "title": "..." }
func (h *Handlers) patchFavoriteOverride(c *gin.Context) {
	userID := core.UserID(c.Param("user_id"))
	assetID := core.AssetID(c.Param("asset_id"))

	var body map[string]any
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	// whitelist allowed shared fields
	override := map[string]any{}
	if v, ok := body["description"].(string); ok && strings.TrimSpace(v) != "" {
		override["description"] = v
	}
	if v, ok := body["title"].(string); ok && strings.TrimSpace(v) != "" {
		override["title"] = v
	}

	if len(override) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no supported fields to update"})
		return
	}

	if err := h.favoriteSvc.SetFavoriteOverride(userID, assetID, override, h.now()); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func SetupRouter(h *Handlers) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	v1 := r.Group("/v1")
	{
		v1.GET("/users/:user_id/favorites", h.getFavorites)
		v1.POST("/users/:user_id/favorites", h.addFavorite)
		v1.DELETE("/users/:user_id/favorites/:asset_id", h.removeFavorite)
		v1.PATCH("/users/:user_id/favorites/:asset_id", h.patchFavoriteOverride)
	}

	return r
}
