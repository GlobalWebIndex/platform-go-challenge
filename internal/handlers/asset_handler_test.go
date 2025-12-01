package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"platform-go-challenge/internal/common"
	"platform-go-challenge/internal/config"
	"platform-go-challenge/internal/models"
	"platform-go-challenge/internal/storage"
)

func setupAssetTestEnv() (*config.Config, *storage.FavoritesStore, *AssetHandler) {
	// Setup the test environment: config, storage, and handler
	cfg := config.NewConfig() // Use NewConfig to initialize the Config
	favorites := storage.NewFavoritesStore()

	handler := NewAssetHandler(cfg, favorites)
	return cfg, favorites, handler
}

// Test for AssetHandler's list function
func TestAssetHandler_List(t *testing.T) {
	cfg, favorites, handler := setupAssetTestEnv()

	// Create assets
	userID := uuid.New()
	chartAsset := &models.Chart{
		BaseAsset: models.BaseAsset{
			ID:          "1",
			Description: "Chart asset",
		},
		Title: "Revenue over time",
		Data: []models.ChartDataPoint{
			{Label: "Q1", Value: 1000},
			{Label: "Q2", Value: 2000},
		},
	}
	insightAsset := &models.Insight{
		BaseAsset: models.BaseAsset{
			ID:          "2",
			Description: "Insight asset",
		},
		InsightText: "Q2 growth is significant.",
	}

	// Add mock favorites
	favorites.AddFavorite(userID, chartAsset)
	favorites.AddFavorite(userID, insightAsset)

	// Create another user (for the forbidden test case)
	anotherUserID := uuid.New()

	// Setup test cases
	tests := []struct {
		name           string
		role           models.Role
		userID         uuid.UUID
		queryUserID    string
		expectedStatus int
		expectedAssets []models.Asset
	}{
		{
			name:           "Unauthorized", // This test checks for missing user context
			role:           models.RoleUser,
			expectedStatus: http.StatusUnauthorized, // Expect 401 when no valid context is passed
		},
		{
			name:           "Forbidden Access for Non-Admin", // Non-admin trying to access another user's assets
			role:           models.RoleUser,
			userID:         userID,
			queryUserID:    anotherUserID.String(), // Simulate a different user
			expectedStatus: http.StatusForbidden,   // Expect 403 if role is not Admin
		},
		{
			name:           "Success for Admin", // Admin should have access to all assets
			role:           models.RoleAdmin,
			userID:         userID,
			expectedStatus: http.StatusOK,
			expectedAssets: []models.Asset{
				chartAsset,
				insightAsset,
			},
		},
		{
			name:           "Success for Admin (No user_id)", // Admin without a user_id should get their own favorites
			role:           models.RoleAdmin,
			userID:         userID,
			expectedStatus: http.StatusOK,
			expectedAssets: []models.Asset{
				chartAsset,
				insightAsset,
			},
		},
		{
			name:           "Success for Normal User", // Normal user should get their own favorites
			role:           models.RoleUser,
			userID:         userID,
			expectedStatus: http.StatusOK,
			expectedAssets: []models.Asset{
				chartAsset,
				insightAsset,
			},
		},
		{
			name:           "Forbidden for Normal User with Different user_id", // Normal user trying to view another user's favorites
			role:           models.RoleUser,
			userID:         userID,
			queryUserID:    anotherUserID.String(),
			expectedStatus: http.StatusForbidden,
		},
	}

	// Running the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate request with the role and user ID
			req := httptest.NewRequest("GET", "/assets", nil)
			req = attachRole(req, cfg, tt.role)
			req = attachUserID(req, cfg, tt.userID)

			// Set query params if any
			if tt.queryUserID != "" {
				query := req.URL.Query()             // Get a reference to the query parameters
				query.Set("user_id", tt.queryUserID) // Modify the query
				req.URL.RawQuery = query.Encode()    // Set the modified query back
			}

			// Set up the recorder and handle the request
			rr := httptest.NewRecorder()
			handler.list(rr, req)

			// Assert the correct status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// If the status is OK, check the returned assets
			if tt.expectedStatus == http.StatusOK {
				var response struct {
					Assets []struct {
						Type  string          `json:"type"`
						Asset json.RawMessage `json:"asset"`
					} `json:"assets"`
				}

				// Unmarshal the response into the wrapper struct
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)

				var assets []models.Asset
				for _, assetResponse := range response.Assets {
					// Use the "type" field to determine the asset type and unmarshal the asset data
					var asset models.Asset
					switch assetResponse.Type {
					case "chart":
						var chart models.Chart
						err := json.Unmarshal(assetResponse.Asset, &chart)
						assert.NoError(t, err)
						asset = &chart
					case "insight":
						var insight models.Insight
						err := json.Unmarshal(assetResponse.Asset, &insight)
						assert.NoError(t, err)
						asset = &insight
					// Add other cases as necessary...
					default:
						t.Errorf("unknown asset type: %s", assetResponse.Type)
					}

					// Append the unmarshaled asset to the assets slice
					assets = append(assets, asset)
				}

				// Compare the unmarshaled assets with the expected assets
				assert.ElementsMatch(t, tt.expectedAssets, assets)
			}
		})
	}
}

// Test for AssetHandler's add function
func TestAssetHandler_Add(t *testing.T) {
	cfg, _, handler := setupAssetTestEnv()

	tests := []struct {
		name           string
		role           models.Role
		userID         uuid.UUID
		payload        string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Unauthorized",
			role:           models.RoleUser,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error": "` + common.ErrUnauthorized + `"}`,
		},
		{
			name:           "Forbidden for User Trying to Add Other User's Favorite",
			role:           models.RoleUser,
			userID:         uuid.New(),
			payload:        `{"user_id": "` + uuid.New().String() + `", "asset": {"type": "chart", "id": "3", "description": "New Chart Asset", "title": "New Revenue Chart", "data": [{"label": "Q3", "value": 3000}]}}`,
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"error": "` + common.ErrForbiddenAccess + `"}`,
		},
		{
			name:           "Success for Admin",
			role:           models.RoleAdmin,
			userID:         uuid.New(),
			payload:        `{"user_id": "` + uuid.New().String() + `", "asset": {"type": "chart", "id": "3", "description": "New Chart Asset", "title": "New Revenue Chart", "data": [{"label": "Q3", "value": 3000}]}}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"status":"favorite added"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/assets", nil)
			req = attachRole(req, cfg, tt.role)
			req = attachUserID(req, cfg, tt.userID)
			req.Body = io.NopCloser(strings.NewReader(tt.payload))

			rr := httptest.NewRecorder()
			handler.add(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestAssetHandler_Remove(t *testing.T) {
	cfg, favorites, handler := setupAssetTestEnv()

	// Create a test asset
	assetID := "1"
	chartAsset := &models.Chart{
		BaseAsset: models.BaseAsset{
			ID:          assetID,
			Description: "Chart asset",
		},
		Title: "Revenue over time",
		Data: []models.ChartDataPoint{
			{Label: "Q1", Value: 1000},
			{Label: "Q2", Value: 2000},
		},
	}

	userID := uuid.New()
	favorites.AddFavorite(userID, chartAsset)

	tests := []struct {
		name           string
		role           models.Role
		userID         uuid.UUID
		assetID        string
		queryUserID    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Unauthorized",
			role:           models.RoleUser,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error": "` + common.ErrUnauthorized + `"}`,
		},
		{
			name:           "Forbidden Access for Non-Admin",
			role:           models.RoleUser,
			userID:         uuid.New(),
			queryUserID:    uuid.New().String(), // Simulate a different user
			assetID:        "1",
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"error": "` + common.ErrForbiddenAccess + `"}`,
		},
		{
			name:           "Asset Not Found",
			role:           models.RoleAdmin,
			userID:         userID,
			assetID:        "nonexistent-id",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error": "` + common.ErrAssetNotFound + `"}`,
		},
		{
			name:           "Success for Admin",
			role:           models.RoleAdmin,
			userID:         userID,
			assetID:        "1",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status": "favorite removed"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/assets", nil)
			req = attachRole(req, cfg, tt.role)
			req = req.WithContext(context.WithValue(req.Context(), cfg.ContextKeys.UserID, tt.userID))

			// Set the query params if any
			if tt.queryUserID != "" {
				query := req.URL.Query()             // Get a reference to the query parameters
				query.Set("user_id", tt.queryUserID) // Modify the query
				req.URL.RawQuery = query.Encode()    // Set the modified query back
			}

			// Set the asset_id query parameter
			query := req.URL.Query()          // Get a reference to the query parameters
			query.Set("asset_id", tt.assetID) // Modify the query
			req.URL.RawQuery = query.Encode() // Set the modified query back

			rr := httptest.NewRecorder()
			handler.remove(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestAssetHandler_Update(t *testing.T) {
	cfg, favorites, handler := setupAssetTestEnv()

	// Create a test asset
	assetID := "1"
	insightAsset := &models.Insight{
		BaseAsset: models.BaseAsset{
			ID:          assetID,
			Description: "Insight asset",
		},
		InsightText: "something deep",
	}

	userID := uuid.New()
	favorites.AddFavorite(userID, insightAsset)

	tests := []struct {
		name           string
		role           models.Role
		userID         uuid.UUID
		payload        string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Unauthorized",
			role:           models.RoleUser,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error": "` + common.ErrUnauthorized + `"}`,
		},
		{
			name:           "Forbidden for User Trying to Update Other User's Favorite",
			role:           models.RoleUser,
			userID:         uuid.New(), // Different user trying to update the asset
			payload:        `{"user_id": "` + userID.String() + `", "asset": {"id": "1", "type": "insight", "description": "Updated Insight Asset", "insight_text": "Updated insight text"}}`,
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"error": "` + common.ErrForbiddenAccess + `"}`,
		},
		{
			name:           "Success for Admin",
			role:           models.RoleAdmin,
			userID:         uuid.New(), // Admin can update any asset
			payload:        `{"user_id": "` + userID.String() + `", "asset": {"id": "1", "type": "insight", "description": "Updated Insight Asset", "insight_text": "Updated insight text"}}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status": "favorite updated"}`,
		},
		{
			name:           "Success for Admin (No user_id in payload)", // Admin with no user_id still works
			role:           models.RoleAdmin,
			userID:         userID, // Same user, but admin can update others' favorites
			payload:        `{"asset": {"id": "1", "type": "insight", "description": "Updated Insight Asset", "insight_text": "Updated insight text"}}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status": "favorite updated"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("PUT", "/assets", nil)
			req = attachRole(req, cfg, tt.role)
			req = req.WithContext(context.WithValue(req.Context(), cfg.ContextKeys.UserID, tt.userID))
			req.Body = io.NopCloser(strings.NewReader(tt.payload))

			rr := httptest.NewRecorder()
			handler.update(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())
		})
	}
}
