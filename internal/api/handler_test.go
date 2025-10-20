package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mr-polymath/platform-go-challenge/internal/api"
	"github.com/mr-polymath/platform-go-challenge/internal/auth"
	"github.com/mr-polymath/platform-go-challenge/internal/db"
	"github.com/mr-polymath/platform-go-challenge/internal/model"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	s := db.NewMemoryStore()
	h := &api.Handler{FavouriteRepository: s}

	// Protected routes
	apiRoutes := r.Group("/users")
	apiRoutes.Use(api.AuthMiddleware([]byte("super-secret-key")))
	{
		apiRoutes.POST("/:userId/favourites", h.AddFavourite)
		apiRoutes.GET("/:userId/favourites", h.GetAllFavourites)
		apiRoutes.PATCH("/:userId/favourites/:assetId", h.PatchFavourite)
		apiRoutes.DELETE("/:userId/favourites/:assetId", h.DeleteFavourite)
	}

	return r
}

func TestHandler(t *testing.T) {
	router := setupRouter()

	// Create a new favourite asset
	newAsset := model.NewAsset(model.InsightType, "Test desc", "Data")
	body, _ := json.Marshal(newAsset)

	userId := uuid.New().String()
	token, err := auth.GenerateToken(userId)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/users/"+userId+"/favourites", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Unmarshal response into wrapper struct
	var postResp struct {
		Message   string      `json:"message"`
		Favourite model.Asset `json:"favourite"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &postResp)
	assert.NoError(t, err)
	assert.Equal(t, "favourite added", postResp.Message)
	assert.Equal(t, "Test desc", postResp.Favourite.Description)

	// Now GET the favourites
	getReq := httptest.NewRequest("GET", "/users/"+userId+"/favourites", nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	getResAfterPost := httptest.NewRecorder()
	router.ServeHTTP(getResAfterPost, getReq)

	assert.Equal(t, http.StatusOK, getResAfterPost.Code)

	// GET returns a JSON array of assets

	var resp struct {
		Favourites []model.Asset `json:"favourites"`
	}

	err = json.Unmarshal(getResAfterPost.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp.Favourites, 1)
	assert.Equal(t, "Test desc", resp.Favourites[0].Description)

	// Now PATCH to update the description
	assetID := resp.Favourites[0].ID.String()

	patchBody, _ := json.Marshal(map[string]string{"description": "Updated desc"})
	patchReq := httptest.NewRequest("PATCH", "/users/"+userId+"/favourites/"+assetID, bytes.NewBuffer(patchBody))
	patchReq.Header.Set("Content-Type", "application/json")
	patchReq.Header.Set("Authorization", "Bearer "+token)
	patchRes := httptest.NewRecorder()
	router.ServeHTTP(patchRes, patchReq)

	// handler currently returns 204 No Content (with a JSON message)
	assert.Equal(t, http.StatusNoContent, patchRes.Code)

	// GET again to ensure the description was updated
	getResAfterPatch := httptest.NewRecorder()
	router.ServeHTTP(getResAfterPatch, getReq)

	assert.Equal(t, http.StatusOK, getResAfterPatch.Code)

	err = json.Unmarshal(getResAfterPatch.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp.Favourites, 1)
	assert.Equal(t, "Updated desc", resp.Favourites[0].Description)

	// Now DELETE the favourite
	delReq := httptest.NewRequest("DELETE", "/users/"+userId+"/favourites/"+assetID, nil)
	delReq.Header.Set("Authorization", "Bearer "+token)
	delRes := httptest.NewRecorder()
	router.ServeHTTP(delRes, delReq)

	assert.Equal(t, http.StatusOK, delRes.Code)

	// Final GET to ensure it's removed
	getResAfterDelete := httptest.NewRecorder()
	router.ServeHTTP(getResAfterDelete, getReq)

	assert.Equal(t, http.StatusOK, getResAfterDelete.Code)

	err = json.Unmarshal(getResAfterDelete.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp.Favourites, 0)
}
