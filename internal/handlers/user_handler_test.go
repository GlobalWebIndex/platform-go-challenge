package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"platform-go-challenge/internal/config"
	"platform-go-challenge/internal/models"
	"platform-go-challenge/internal/storage"
)

// Mock GenerateJWT for testing
type mockJWT struct{}

func (m *mockJWT) Generate(cfg *config.Config, userID uuid.UUID, role models.Role) (string, error) {
	return "mocktoken", nil
}

func TestUserHandler_Register(t *testing.T) {
	t.Parallel()

	store := storage.NewUserStore()
	mockJWT := &mockJWT{}
	handler := NewUserHandler(&config.Config{}, store, mockJWT)

	tests := []struct {
		name           string
		payload        map[string]string
		setupExisting  bool
		expectedStatus int
		expectJSON     bool
	}{
		{
			name: "successful registration - user",
			payload: map[string]string{
				"username": "user1",
				"password": "pass123",
				"role":     "user",
			},
			expectedStatus: http.StatusCreated,
			expectJSON:     true,
		},
		{
			name: "successful registration - admin",
			payload: map[string]string{
				"username": "admin1",
				"password": "pass123",
				"role":     "admin",
			},
			expectedStatus: http.StatusCreated,
			expectJSON:     true,
		},
		{
			name: "duplicate username",
			payload: map[string]string{
				"username": "duplicate",
				"password": "pass",
				"role":     "user",
			},
			setupExisting:  true,
			expectedStatus: http.StatusBadRequest,
			expectJSON:     false,
		},
		{
			name: "empty username",
			payload: map[string]string{
				"username": "",
				"password": "pass",
				"role":     "user",
			},
			expectedStatus: http.StatusBadRequest,
			expectJSON:     false,
		},
		{
			name: "empty password",
			payload: map[string]string{
				"username": "user2",
				"password": "",
				"role":     "user",
			},
			expectedStatus: http.StatusBadRequest,
			expectJSON:     false,
		},
		{
			name: "empty role",
			payload: map[string]string{
				"username": "emptyrole",
				"password": "pass123",
				"role":     "",
			},
			expectedStatus: http.StatusBadRequest,
			expectJSON:     false,
		},
		{
			name: "invalid role",
			payload: map[string]string{
				"username": "invalidrole",
				"password": "pass123",
				"role":     "superuser",
			},
			expectedStatus: http.StatusBadRequest,
			expectJSON:     false,
		},
		{
			name:           "invalid JSON",
			payload:        nil,
			expectedStatus: http.StatusBadRequest,
			expectJSON:     false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var body []byte
			if tt.payload != nil {
				body, _ = json.Marshal(tt.payload)
			} else {
				body = []byte("{badjson")
			}

			if tt.setupExisting {
				user, _ := models.CreateUser(tt.payload["username"], tt.payload["password"], models.Role(tt.payload["role"]))
				store.AddUser(&user)
			}

			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
			w := httptest.NewRecorder()

			handler.register(w, req)

			res := w.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectJSON {
				assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

				var resp map[string]interface{}
				json.NewDecoder(res.Body).Decode(&resp)
				assert.Equal(t, tt.payload["username"], resp["username"])
				assert.Equal(t, tt.payload["role"], resp["role"])
				assert.NotEmpty(t, resp["id"])
			}
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	t.Parallel()

	store := storage.NewUserStore()
	mockJWT := &mockJWT{}
	handler := NewUserHandler(&config.Config{}, store, mockJWT)

	// Add a valid user
	user, _ := models.CreateUser("loginuser", "mypassword", models.RoleUser)
	store.AddUser(&user)

	tests := []struct {
		name           string
		payload        map[string]string
		expectedStatus int
		expectJSON     bool
	}{
		{
			name: "successful login",
			payload: map[string]string{
				"username": "loginuser",
				"password": "mypassword",
			},
			expectedStatus: http.StatusOK,
			expectJSON:     true,
		},
		{
			name: "wrong password",
			payload: map[string]string{
				"username": "loginuser",
				"password": "wrongpass",
			},
			expectedStatus: http.StatusUnauthorized,
			expectJSON:     false,
		},
		{
			name: "nonexistent user",
			payload: map[string]string{
				"username": "nouser",
				"password": "pass",
			},
			expectedStatus: http.StatusUnauthorized,
			expectJSON:     false,
		},
		{
			name: "empty username",
			payload: map[string]string{
				"username": "",
				"password": "mypassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectJSON:     false,
		},
		{
			name: "empty password",
			payload: map[string]string{
				"username": "loginuser",
				"password": "",
			},
			expectedStatus: http.StatusUnauthorized,
			expectJSON:     false,
		},
		{
			name: "both empty",
			payload: map[string]string{
				"username": "",
				"password": "",
			},
			expectedStatus: http.StatusUnauthorized,
			expectJSON:     false,
		},
		{
			name:           "invalid JSON",
			payload:        nil,
			expectedStatus: http.StatusBadRequest,
			expectJSON:     false,
		},
		{
			name: "long username/password",
			payload: map[string]string{
				"username": string(make([]byte, 1000)),
				"password": string(make([]byte, 1000)),
			},
			expectedStatus: http.StatusUnauthorized,
			expectJSON:     false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var body []byte
			if tt.payload != nil {
				body, _ = json.Marshal(tt.payload)
			} else {
				body = []byte("{badjson")
			}

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			w := httptest.NewRecorder()

			handler.login(w, req)

			res := w.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectJSON {
				assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

				var resp map[string]string
				json.NewDecoder(res.Body).Decode(&resp)
				assert.Equal(t, "mocktoken", resp["token"])
			}
		})
	}
}
