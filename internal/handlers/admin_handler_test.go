package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"platform-go-challenge/internal/config"
	"platform-go-challenge/internal/models"
	"platform-go-challenge/internal/storage"
)

func setupAdminTestEnv() (*config.Config, *storage.UserStore, *storage.FavoritesStore) {
	cfg := &config.Config{
		JWTSecret:   []byte("secret"),
		AdminAPIKey: "admin-key",
		ServerPort:  ":8080",
		TokenExpiry: time.Hour,
	}
	cfg.ContextKeys.UserID = "userID"
	cfg.ContextKeys.Role = "role"

	users := storage.NewUserStore()
	favorites := storage.NewFavoritesStore()

	return cfg, users, favorites
}

func TestListUsers(t *testing.T) {
	cfg, users, favorites := setupAdminTestEnv()
	handler := NewUserAdminHandler(cfg, users, favorites)

	// Add 2 users for listing
	u1 := &models.User{Username: "john"}
	users.AddUser(u1)
	u2 := &models.User{Username: "mary"}
	users.AddUser(u2)

	// Note: Admin handlers no longer check roles since API key middleware handles authorization
	tests := []struct {
		name           string
		expectedStatus int
		expectedUsers  []string
		clearUsers     bool
	}{
		{
			name:           "success for admin",
			expectedStatus: http.StatusOK,
			expectedUsers:  []string{"john", "mary"},
		},
		{
			name:           "no users in store",
			expectedStatus: http.StatusOK,
			expectedUsers:  []string{}, // Empty list, no users should be returned
			clearUsers:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.clearUsers {
				// Clear users for the "no users" test
				users = storage.NewUserStore()
				handler = NewUserAdminHandler(cfg, users, favorites)
			}

			req := httptest.NewRequest("GET", "/admin/users", nil)

			rr := httptest.NewRecorder()
			handler.listUsers(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				var result map[string]uuid.UUID
				err := json.Unmarshal(rr.Body.Bytes(), &result)
				assert.NoError(t, err)

				// Check if the expected users are in the response
				for _, username := range tt.expectedUsers {
					assert.Contains(t, result, username)
				}
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	cfg, users, favorites := setupAdminTestEnv()
	handler := NewUserAdminHandler(cfg, users, favorites)

	// Add a user to delete
	user := &models.User{Username: "victim"}
	users.AddUser(user)

	// Add some favorites so we can confirm they get cleaned up
	fakeAsset := &models.Chart{} // any Asset implementation will work
	favorites.AddFavorite(user.ID, fakeAsset)

	// Note: Admin handlers no longer check roles since API key middleware handles authorization
	tests := []struct {
		name           string
		userID         string
		expectedStatus int
	}{
		{
			name:           "bad request for invalid user ID format",
			userID:         "not-a-uuid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "not found for non-existing user",
			userID:         uuid.New().String(),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "successful delete",
			userID:         user.ID.String(),
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/admin/users/"+tt.userID, nil)

			req = mux.SetURLVars(req, map[string]string{"id": tt.userID})
			rr := httptest.NewRecorder()

			handler.deleteUser(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Additional checks for successful deletion
			if tt.expectedStatus == http.StatusNoContent {
				// Confirm user removed
				_, exists := users.GetUserByUserID(user.ID)
				assert.False(t, exists)

				// Confirm favorites cleaned up (should return empty list for non-existent user)
				assets, err := favorites.GetUserFavorites(user.ID)
				assert.NoError(t, err) // GetUserFavorites returns empty list, not error
				assert.Empty(t, assets)
			}
		})
	}
}
