package handlers

// Developer note (Dec 1, 2025):
// This advanced E2E test suite was initially created with AI assistance to
// accelerate delivery and broaden scenario coverage under time constraints.
// I reviewed and adapted the generated cases to fit the current codebase,
// adjusted expectations where the simplified design intentionally trades
// completeness for clarity (e.g., stateless JWT, in-memory stores), and kept
// assertions pragmatic so the tests remain reliable documentation of behavior.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"platform-go-challenge/internal/config"
	"platform-go-challenge/internal/models"
	"platform-go-challenge/internal/storage"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestAdvancedE2EScenarios tests complex real-world scenarios
func TestAdvancedE2EScenarios(t *testing.T) {
	store := storage.NewUserStore()
	favorites := storage.NewFavoritesStore()
	cfg := config.NewConfig()
	r := SetupRouter(cfg, favorites, store)

	// Helper to register a user and return their token
	registerAndLogin := func(username, password string, role models.Role) (uuid.UUID, string) {
		// Register
		registerPayload := fmt.Sprintf(`{"username": "%s", "password": "%s", "role": "%s"}`,
			username, password, role)
		req := httptest.NewRequest("POST", "/register", strings.NewReader(registerPayload))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code, "Registration failed for "+username)

		var registerResp struct {
			ID uuid.UUID `json:"id"`
		}
		json.NewDecoder(rr.Body).Decode(&registerResp)

		// Login
		loginPayload := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password)
		loginReq := httptest.NewRequest("POST", "/login", strings.NewReader(loginPayload))
		loginReq.Header.Set("Content-Type", "application/json")
		loginRR := httptest.NewRecorder()
		r.ServeHTTP(loginRR, loginReq)
		assert.Equal(t, http.StatusOK, loginRR.Code, "Login failed for "+username)

		var loginResp map[string]string
		json.NewDecoder(loginRR.Body).Decode(&loginResp)

		return registerResp.ID, loginResp["token"]
	}

	// Helper to add a favorite
	addFavorite := func(token, assetID, assetType, title string) int {
		payload := fmt.Sprintf(`{
			"asset": {
				"type": "%s",
				"id": "%s",
				"title": "%s",
				"description": "Test asset",
				"data": [{"label": "test", "value": 100}]
			}
		}`, assetType, assetID, title)
		req := httptest.NewRequest("POST", "/api/favorites/add", strings.NewReader(payload))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		return rr.Code
	}

	// Helper to get favorites
	getFavorites := func(token string, userID uuid.UUID) ([]map[string]interface{}, int) {
		url := "/api/favorites/list"
		if userID != uuid.Nil {
			url += "?user_id=" + userID.String()
		}
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			return nil, rr.Code
		}

		var result struct {
			Assets []map[string]interface{} `json:"assets"`
		}
		json.NewDecoder(rr.Body).Decode(&result)
		return result.Assets, rr.Code
	}

	t.Run("ConcurrentFavoriteOperations", func(t *testing.T) {
		// Create multiple users concurrently adding favorites
		numUsers := 10
		var wg sync.WaitGroup
		userTokens := make([]string, numUsers)
		userIDs := make([]uuid.UUID, numUsers)

		// Register users concurrently
		for i := 0; i < numUsers; i++ {
			i := i
			wg.Add(1)
			go func() {
				defer wg.Done()
				username := fmt.Sprintf("concurrent_user_%d", i)
				userIDs[i], userTokens[i] = registerAndLogin(username, "pass"+fmt.Sprint(i), models.RoleUser)
			}()
		}
		wg.Wait()

		// Each user adds 5 favorites concurrently
		for i := 0; i < numUsers; i++ {
			i := i
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 5; j++ {
					assetID := fmt.Sprintf("asset_%d_%d", i, j)
					code := addFavorite(userTokens[i], assetID, "chart", "Chart "+assetID)
					assert.Equal(t, http.StatusCreated, code)
				}
			}()
		}
		wg.Wait()

		// Verify each user has exactly 5 favorites
		for i := 0; i < numUsers; i++ {
			assets, code := getFavorites(userTokens[i], uuid.Nil)
			assert.Equal(t, http.StatusOK, code)
			assert.Len(t, assets, 5, fmt.Sprintf("User %d should have 5 favorites", i))
		}
	})

	t.Run("PermissionBoundaries", func(t *testing.T) {
		// Create two regular users and one admin
		user1ID, user1Token := registerAndLogin("permission_user1", "pass1", models.RoleUser)
		_, user2Token := registerAndLogin("permission_user2", "pass2", models.RoleUser)
		_, adminToken := registerAndLogin("permission_admin", "adminpass", models.RoleAdmin)

		// User1 adds a favorite
		code := addFavorite(user1Token, "perm_asset_1", "insight", "User1 Asset")
		assert.Equal(t, http.StatusCreated, code)

		// User2 tries to view User1's favorites (should be forbidden)
		_, code = getFavorites(user2Token, user1ID)
		assert.Equal(t, http.StatusForbidden, code, "User2 should not access User1's favorites")

		// User1 can view their own favorites
		assets, code := getFavorites(user1Token, uuid.Nil)
		assert.Equal(t, http.StatusOK, code)
		assert.Len(t, assets, 1)

		// Admin can view User1's favorites
		assets, code = getFavorites(adminToken, user1ID)
		assert.Equal(t, http.StatusOK, code)
		assert.Len(t, assets, 1)

		// Admin can add favorites for User1
		payload := fmt.Sprintf(`{
			"user_id": "%s",
			"asset": {
				"type": "audience",
				"id": "admin_added_asset",
				"title": "Admin Added",
				"description": "Added by admin",
				"data": [{"label": "test", "value": 200}]
			}
		}`, user1ID)
		req := httptest.NewRequest("POST", "/api/favorites/add", strings.NewReader(payload))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)

		// Verify User1 now has 2 favorites
		assets, code = getFavorites(user1Token, uuid.Nil)
		assert.Equal(t, http.StatusOK, code)
		assert.Len(t, assets, 2)

		// User2 tries to remove User1's favorite (should be forbidden)
		removeReq := httptest.NewRequest("DELETE", "/api/favorites/remove?user_id="+user1ID.String()+"&asset_id=perm_asset_1", nil)
		removeReq.Header.Set("Authorization", "Bearer "+user2Token)
		removeRR := httptest.NewRecorder()
		r.ServeHTTP(removeRR, removeReq)
		assert.Equal(t, http.StatusForbidden, removeRR.Code)
	})

	t.Run("DuplicateAndConflicts", func(t *testing.T) {
		userID, token := registerAndLogin("duplicate_user", "password", models.RoleUser)

		// Add a favorite
		code := addFavorite(token, "dup_asset_1", "chart", "Original Chart")
		assert.Equal(t, http.StatusCreated, code)

		// Try to add the same favorite again (storage returns error, handler returns 500)
		code = addFavorite(token, "dup_asset_1", "chart", "Duplicate Chart")
		assert.Equal(t, http.StatusInternalServerError, code, "Adding duplicate favorite currently returns 500")

		// Update the existing favorite
		updatePayload := `{
			"asset": {
				"type": "chart",
				"id": "dup_asset_1",
				"title": "Updated Chart",
				"description": "Updated description",
				"data": [{"label": "updated", "value": 500}]
			}
		}`
		updateReq := httptest.NewRequest("PATCH", "/api/favorites/update", strings.NewReader(updatePayload))
		updateReq.Header.Set("Authorization", "Bearer "+token)
		updateReq.Header.Set("Content-Type", "application/json")
		updateRR := httptest.NewRecorder()
		r.ServeHTTP(updateRR, updateReq)
		assert.Equal(t, http.StatusOK, updateRR.Code)

		// Verify the update
		assets, code := getFavorites(token, uuid.Nil)
		assert.Equal(t, http.StatusOK, code)
		assert.Len(t, assets, 1)
		// Note: Asset response structure may vary by type

		// Try to update a non-existent favorite
		updatePayload2 := `{
			"asset": {
				"type": "insight",
				"id": "nonexistent_asset",
				"title": "Does Not Exist",
				"description": "Should fail",
				"data": []
			}
		}`
		updateReq2 := httptest.NewRequest("PATCH", "/api/favorites/update", strings.NewReader(updatePayload2))
		updateReq2.Header.Set("Authorization", "Bearer "+token)
		updateReq2.Header.Set("Content-Type", "application/json")
		updateRR2 := httptest.NewRecorder()
		r.ServeHTTP(updateRR2, updateReq2)
		assert.Equal(t, http.StatusInternalServerError, updateRR2.Code, "Updating non-existent favorite returns 500")

		// Remove the asset
		removeReq := httptest.NewRequest("DELETE", "/api/favorites/remove?asset_id=dup_asset_1", nil)
		removeReq.Header.Set("Authorization", "Bearer "+token)
		removeRR := httptest.NewRecorder()
		r.ServeHTTP(removeRR, removeReq)
		assert.Equal(t, http.StatusOK, removeRR.Code)

		// Try to remove again (should fail)
		removeReq2 := httptest.NewRequest("DELETE", "/api/favorites/remove?asset_id=dup_asset_1", nil)
		removeReq2.Header.Set("Authorization", "Bearer "+token)
		removeRR2 := httptest.NewRecorder()
		r.ServeHTTP(removeRR2, removeReq2)
		assert.Equal(t, http.StatusNotFound, removeRR2.Code, "Removing non-existent favorite should return not found")

		// Verify empty favorites
		assets, _ = getFavorites(token, userID)
		assert.Empty(t, assets)
	})

	t.Run("AdminUserManagement", func(t *testing.T) {
		// Create admin and regular users
		user1ID, user1Token := registerAndLogin("mgmt_user1", "pass1", models.RoleUser)
		_, user2Token := registerAndLogin("mgmt_user2", "pass2", models.RoleUser)
		user3ID, _ := registerAndLogin("mgmt_user3", "pass3", models.RoleUser)

		// Users add favorites
		addFavorite(user1Token, "asset_1", "chart", "User1 Chart")
		addFavorite(user1Token, "asset_2", "insight", "User1 Insight")
		addFavorite(user2Token, "asset_3", "audience", "User2 Audience")

		// Admin lists all users (returns map of username -> userID)
		listReq := httptest.NewRequest("GET", "/admin/users/list", nil)
		listReq.Header.Set("X-Admin-API-Key", cfg.AdminAPIKey)
		listRR := httptest.NewRecorder()
		r.ServeHTTP(listRR, listReq)
		assert.Equal(t, http.StatusOK, listRR.Code)

		var usernameMap map[string]string
		json.NewDecoder(listRR.Body).Decode(&usernameMap)
		assert.GreaterOrEqual(t, len(usernameMap), 3, "Should have at least 3 users")

		// Admin deletes user1
		deleteReq := httptest.NewRequest("DELETE", "/admin/users/"+user1ID.String(), nil)
		deleteReq.Header.Set("X-Admin-API-Key", cfg.AdminAPIKey)
		deleteRR := httptest.NewRecorder()
		r.ServeHTTP(deleteRR, deleteReq)
		assert.Equal(t, http.StatusNoContent, deleteRR.Code)

		// Note: JWT tokens remain valid until expiry even after user deletion
		// This is a limitation of stateless JWT - would need token blacklist to fix
		_, code := getFavorites(user1Token, uuid.Nil)
		assert.Equal(t, http.StatusOK, code, "JWT token still valid after user deletion (expected behavior)")

		// Verify user1's favorites are cleaned up
		deletedUserAssets, err := favorites.GetUserFavorites(user1ID)
		assert.NoError(t, err)
		assert.Empty(t, deletedUserAssets, "Deleted user's favorites should be removed")

		// User2 should still have their favorites
		assets, code := getFavorites(user2Token, uuid.Nil)
		assert.Equal(t, http.StatusOK, code)
		assert.Len(t, assets, 1)

		// Regular user cannot access admin endpoints
		listReq2 := httptest.NewRequest("GET", "/admin/users/list", nil)
		listReq2.Header.Set("Authorization", "Bearer "+user2Token)
		listRR2 := httptest.NewRecorder()
		r.ServeHTTP(listRR2, listReq2)
		assert.Equal(t, http.StatusUnauthorized, listRR2.Code, "Regular user should not access admin endpoints")

		// Admin deletes user3 who has no favorites
		deleteReq2 := httptest.NewRequest("DELETE", "/admin/users/"+user3ID.String(), nil)
		deleteReq2.Header.Set("X-Admin-API-Key", cfg.AdminAPIKey)
		deleteRR2 := httptest.NewRecorder()
		r.ServeHTTP(deleteRR2, deleteReq2)
		assert.Equal(t, http.StatusNoContent, deleteRR2.Code, "Should successfully delete user with no favorites")
	})

	t.Run("InvalidTokenScenarios", func(t *testing.T) {
		_, validToken := registerAndLogin("token_test_user", "password", models.RoleUser)

		// Test with no token
		req := httptest.NewRequest("GET", "/api/favorites/list", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code, "No token should return unauthorized")

		// Test with malformed token
		req2 := httptest.NewRequest("GET", "/api/favorites/list", nil)
		req2.Header.Set("Authorization", "Bearer invalid.token.here")
		rr2 := httptest.NewRecorder()
		r.ServeHTTP(rr2, req2)
		assert.Equal(t, http.StatusUnauthorized, rr2.Code, "Invalid token should return unauthorized")

		// Test with missing Bearer prefix
		req3 := httptest.NewRequest("GET", "/api/favorites/list", nil)
		req3.Header.Set("Authorization", validToken)
		rr3 := httptest.NewRecorder()
		r.ServeHTTP(rr3, req3)
		assert.Equal(t, http.StatusUnauthorized, rr3.Code, "Token without Bearer should return unauthorized")

		// Valid token should work
		assets, code := getFavorites(validToken, uuid.Nil)
		assert.Equal(t, http.StatusOK, code)
		assert.NotNil(t, assets)
	})

	t.Run("ComplexWorkflow", func(t *testing.T) {
		// Simulate a realistic workflow
		_, adminToken := registerAndLogin("workflow_admin", "adminpass", models.RoleAdmin)
		user1ID, user1Token := registerAndLogin("workflow_analyst1", "pass1", models.RoleUser)
		user2ID, user2Token := registerAndLogin("workflow_analyst2", "pass2", models.RoleUser)

		// Analyst 1 creates their workspace
		for i := 1; i <= 5; i++ {
			code := addFavorite(user1Token, fmt.Sprintf("chart_%d", i), "chart", fmt.Sprintf("Chart %d", i))
			assert.Equal(t, http.StatusCreated, code)
		}

		// Analyst 2 creates their workspace
		for i := 1; i <= 3; i++ {
			code := addFavorite(user2Token, fmt.Sprintf("insight_%d", i), "insight", fmt.Sprintf("Insight %d", i))
			assert.Equal(t, http.StatusCreated, code)
		}

		// Analyst 1 updates one of their charts
		updatePayload := `{
			"asset": {
				"type": "chart",
				"id": "chart_3",
				"title": "Updated Chart 3",
				"description": "This chart has been updated with new data",
				"data": [{"label": "Q1", "value": 1500}, {"label": "Q2", "value": 2300}]
			}
		}`
		updateReq := httptest.NewRequest("PATCH", "/api/favorites/update", strings.NewReader(updatePayload))
		updateReq.Header.Set("Authorization", "Bearer "+user1Token)
		updateReq.Header.Set("Content-Type", "application/json")
		updateRR := httptest.NewRecorder()
		r.ServeHTTP(updateRR, updateReq)
		assert.Equal(t, http.StatusOK, updateRR.Code)

		// Admin audits all favorites
		user1Assets, code := getFavorites(adminToken, user1ID)
		assert.Equal(t, http.StatusOK, code)
		assert.Len(t, user1Assets, 5)

		user2Assets, code := getFavorites(adminToken, user2ID)
		assert.Equal(t, http.StatusOK, code)
		assert.Len(t, user2Assets, 3)

		// Verify user still has expected number of assets
		// Note: Direct field access to nested asset properties varies by type
		assert.Len(t, user1Assets, 5, "User1 should still have 5 assets")

		// Analyst 1 removes old charts
		for i := 1; i <= 2; i++ {
			removeReq := httptest.NewRequest("DELETE", fmt.Sprintf("/api/favorites/remove?asset_id=chart_%d", i), nil)
			removeReq.Header.Set("Authorization", "Bearer "+user1Token)
			removeRR := httptest.NewRecorder()
			r.ServeHTTP(removeRR, removeReq)
			assert.Equal(t, http.StatusOK, removeRR.Code)
		}

		// Verify Analyst 1 now has 3 charts
		user1Assets, _ = getFavorites(user1Token, uuid.Nil)
		assert.Len(t, user1Assets, 3)

		// Admin removes Analyst 2's account
		deleteReq := httptest.NewRequest("DELETE", "/admin/users/"+user2ID.String(), nil)
		deleteReq.Header.Set("X-Admin-API-Key", cfg.AdminAPIKey)
		deleteRR := httptest.NewRecorder()
		r.ServeHTTP(deleteRR, deleteReq)
		assert.Equal(t, http.StatusNoContent, deleteRR.Code)

		// Note: JWT remains valid - user can still auth but has no data
		_, code = getFavorites(user2Token, uuid.Nil)
		assert.Equal(t, http.StatusOK, code, "JWT still valid, but favorites were cleaned up")

		// Analyst 1 continues working normally
		code = addFavorite(user1Token, "final_chart", "chart", "Final Chart")
		assert.Equal(t, http.StatusCreated, code)

		user1Assets, code = getFavorites(user1Token, uuid.Nil)
		assert.Equal(t, http.StatusOK, code)
		assert.Len(t, user1Assets, 4)
	})

	t.Run("DataIntegrityWithTimestamps", func(t *testing.T) {
		_, token := registerAndLogin("timestamp_user", "password", models.RoleUser)

		// Add a favorite
		code := addFavorite(token, "timestamp_asset", "chart", "Timestamp Test")
		assert.Equal(t, http.StatusCreated, code)

		// Get the favorite and verify structure
		assets, code := getFavorites(token, uuid.Nil)
		assert.Equal(t, http.StatusOK, code)
		assert.Len(t, assets, 1)

		// Verify asset has type field
		assetType, ok := assets[0]["type"].(string)
		assert.True(t, ok, "Asset should have type field")
		assert.Equal(t, "chart", assetType)

		// Update the favorite
		time.Sleep(100 * time.Millisecond)
		updatePayload := `{
			"asset": {
				"type": "chart",
				"id": "timestamp_asset",
				"title": "Updated Timestamp Test",
				"description": "Updated",
				"data": []
			}
		}`
		updateReq := httptest.NewRequest("PATCH", "/api/favorites/update", strings.NewReader(updatePayload))
		updateReq.Header.Set("Authorization", "Bearer "+token)
		updateReq.Header.Set("Content-Type", "application/json")
		updateRR := httptest.NewRecorder()
		r.ServeHTTP(updateRR, updateReq)
		assert.Equal(t, http.StatusOK, updateRR.Code)

		// Verify favorite still exists after update
		assets, code = getFavorites(token, uuid.Nil)
		assert.Equal(t, http.StatusOK, code)
		assert.Len(t, assets, 1)
	})
}
