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
	}

	return r
}

func TestAddAndGetFavourite(t *testing.T) {
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
	getRes := httptest.NewRecorder()
	getReq.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(getRes, getReq)

	assert.Equal(t, http.StatusOK, getRes.Code)

	// GET returns a JSON array of assets

	var resp struct {
		Favourites []model.Asset `json:"favourites"`
	}

	err = json.Unmarshal(getRes.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp.Favourites, 1)
	assert.Equal(t, "Test desc", resp.Favourites[0].Description)
}
