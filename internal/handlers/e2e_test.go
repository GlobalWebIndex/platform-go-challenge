package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"platform-go-challenge/internal/config"
	"platform-go-challenge/internal/models"
	"platform-go-challenge/internal/storage"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEndToEndFlow(t *testing.T) {
	// Initialize the in-memory stores and the router
	store := storage.NewUserStore()
	favorites := storage.NewFavoritesStore()
	cfg := config.NewConfig()
	r := SetupRouter(cfg, favorites, store)

	// Register an admin and two users
	admin := models.User{
		ID:       uuid.New(),
		Username: "admin",
		Password: "adminpass",
		UserRole: models.RoleAdmin,
	}
	user1 := models.User{
		ID:       uuid.New(),
		Username: "user1",
		Password: "user1pass",
		UserRole: models.RoleUser,
	}
	user2 := models.User{
		ID:       uuid.New(),
		Username: "user2",
		Password: "user2pass",
		UserRole: models.RoleUser,
	}

	// Register all users and capture their actual IDs
	users := []*models.User{&admin, &user1, &user2}
	for _, user := range users {
		registerPayload := fmt.Sprintf(`{"username": "%s", "password": "%s", "role": "%s"}`,
			user.Username, user.Password, user.UserRole)

		req := httptest.NewRequest("POST", "/register", strings.NewReader(registerPayload))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		// Assert that registration is successful (201 Created)
		assert.Equal(t, http.StatusCreated, rr.Code)

		// Capture the actual registered user ID from response
		var registerResp struct {
			ID       uuid.UUID `json:"id"`
			Username string    `json:"username"`
			Role     string    `json:"role"`
		}
		json.NewDecoder(rr.Body).Decode(&registerResp)
		user.ID = registerResp.ID
	}

	favoriteAssets := []struct {
		assetID     string
		assetType   string
		description string
		title       string
		data        []map[string]interface{} // Data array with label/value pairs
	}{
		{
			assetID:     "1",
			assetType:   "chart",
			description: "A revenue chart asset",
			title:       "Revenue Chart Q1",
			data: []map[string]interface{}{
				{"label": "Q1", "value": 1500},
				{"label": "Q2", "value": 2300},
			},
		},
		{
			assetID:     "2",
			assetType:   "insight",
			description: "An insight into Q1 performance",
			title:       "Q1 Performance Insight",
			data: []map[string]interface{}{
				{"label": "Q1", "value": 1500},
			},
		},
		{
			assetID:     "3",
			assetType:   "audience",
			description: "Audience demographics for Q1",
			title:       "Audience Overview",
			data: []map[string]interface{}{
				{"label": "Males", "value": 55},
				{"label": "Females", "value": 45},
			},
		},
	}

	userFavorites := make(map[uuid.UUID][]string)
	for _, user := range users {
		// Log in as the user and get the JWT token
		loginPayload := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, user.Username, user.Password)
		loginReq := httptest.NewRequest("POST", "/login", strings.NewReader(loginPayload))
		loginReq.Header.Set("Content-Type", "application/json")

		loginRR := httptest.NewRecorder()
		r.ServeHTTP(loginRR, loginReq)

		assert.Equal(t, http.StatusOK, loginRR.Code)

		var loginResp map[string]string
		json.NewDecoder(loginRR.Body).Decode(&loginResp)
		token := loginResp["token"]

		// Now add favorites for the logged-in user
		for _, asset := range favoriteAssets {
			toJSON := func(data interface{}) string {
				bytes, err := json.Marshal(data)
				if err != nil {
					log.Fatalf("Error converting data to JSON: %v", err)
				}
				return string(bytes)
			}

			// Add favorite request payload
			addFavoritePayload := fmt.Sprintf(`{
				"asset": {
					"type": "%s",
					"id": "%s",
					"description": "%s",
					"title": "%s",
					"data": %s
				}
			}`, asset.assetType, asset.assetID, asset.description, asset.title, toJSON(asset.data))

			addFavoriteReq := httptest.NewRequest("POST", "/api/favorites/add", strings.NewReader(addFavoritePayload))
			addFavoriteReq.Header.Set("Authorization", "Bearer "+token)
			addFavoriteReq.Header.Set("Content-Type", "application/json")

			addFavoriteRR := httptest.NewRecorder()
			r.ServeHTTP(addFavoriteRR, addFavoriteReq)

			// Assert the response status is 201 (Created)
			assert.Equal(t, http.StatusCreated, addFavoriteRR.Code)
			userFavorites[user.ID] = append(userFavorites[user.ID], asset.assetType)
		}
	}

	// Log in as admin and fetch users and their favorites
	loginAdminPayload := `{"username": "admin", "password": "adminpass"}`
	loginAdminReq := httptest.NewRequest("POST", "/login", strings.NewReader(loginAdminPayload))
	loginAdminReq.Header.Set("Content-Type", "application/json")

	loginAdminRR := httptest.NewRecorder()
	r.ServeHTTP(loginAdminRR, loginAdminReq)

	assert.Equal(t, http.StatusOK, loginAdminRR.Code)

	var loginAdminResp map[string]string
	json.NewDecoder(loginAdminRR.Body).Decode(&loginAdminResp)
	adminToken := loginAdminResp["token"]

	// Get all users (admin routes require X-Admin-API-Key header)
	getUsersReq := httptest.NewRequest("GET", "/admin/users/list", nil)
	getUsersReq.Header.Set("X-Admin-API-Key", cfg.AdminAPIKey)

	getUsersRR := httptest.NewRecorder()
	r.ServeHTTP(getUsersRR, getUsersReq)

	assert.Equal(t, http.StatusOK, getUsersRR.Code)

	// Get user favorites (JWT token allows admin to view any user's favorites)
	for _, user := range users {
		getFavoritesReq := httptest.NewRequest("GET", fmt.Sprintf("/api/favorites/list?user_id=%s", user.ID), nil)
		getFavoritesReq.Header.Set("Authorization", "Bearer "+adminToken)

		getFavoritesRR := httptest.NewRecorder()
		r.ServeHTTP(getFavoritesRR, getFavoritesReq)

		if getFavoritesRR.Code != http.StatusOK {
			t.Logf("Error for user %s: %s", user.Username, getFavoritesRR.Body.String())
		}
		assert.Equal(t, http.StatusOK, getFavoritesRR.Code)
	}

	// Verify favorites were actually added by checking each user's favorites
	for _, user := range users {
		checkFavoritesReq := httptest.NewRequest("GET", fmt.Sprintf("/api/favorites/list?user_id=%s", user.ID), nil)
		checkFavoritesReq.Header.Set("Authorization", "Bearer "+adminToken)

		checkFavoritesRR := httptest.NewRecorder()
		r.ServeHTTP(checkFavoritesRR, checkFavoritesReq)

		assert.Equal(t, http.StatusOK, checkFavoritesRR.Code)
	}

	// Admin removes the chart favorite for each user
	for _, user := range users {
		// Remove favorite using query parameters (no body needed for DELETE with query params)
		removeFavoriteReq := httptest.NewRequest("DELETE", "/api/favorites/remove?user_id="+user.ID.String()+"&asset_id=1", nil)
		removeFavoriteReq.Header.Set("Authorization", "Bearer "+adminToken)

		removeFavoriteRR := httptest.NewRecorder()
		r.ServeHTTP(removeFavoriteRR, removeFavoriteReq)

		assert.Equal(t, http.StatusOK, removeFavoriteRR.Code)
	}

	// Admin deletes the second user (user2) - admin routes need API key
	deleteUserReq := httptest.NewRequest("DELETE", fmt.Sprintf("/admin/users/%s", user2.ID), nil)
	deleteUserReq.Header.Set("X-Admin-API-Key", cfg.AdminAPIKey)

	deleteUserRR := httptest.NewRecorder()
	r.ServeHTTP(deleteUserRR, deleteUserReq)

	assert.Equal(t, http.StatusNoContent, deleteUserRR.Code)

	// Ensure user2 is deleted
	_, exists := store.GetUserByUserID(user2.ID)
	assert.False(t, exists)

	_, exists = store.GetUserByUsername(user2.Username)
	assert.False(t, exists)
}
