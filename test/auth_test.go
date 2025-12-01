package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/dintzeler/platform-go-challenge/controllers"
	"strings"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

func init() {
    // Get the directory of this test file
    _, filename, _, _ := runtime.Caller(0)
    testDir := filepath.Dir(filename)
    
    // Go up one level to the project root
    projectRoot := filepath.Dir(testDir)
    
    // Set the absolute path to test_data.json
    dataFile := filepath.Join(projectRoot, "test_data.json")
    os.Setenv("DATA_FILE", dataFile)
    
    // Verify the file exists
    if _, err := os.Stat(dataFile); os.IsNotExist(err) {
        panic("test_data.json not found at: " + dataFile)
    }
}

func TestLogin(t *testing.T) {
	t.Run("Login with empty credentials", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(""))
		controllers.LoginHandler(w, req)

		testErrorResponse(t, w, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST_BODY")
	})

	t.Run("Login with missing email", func(t *testing.T) {
		w := httptest.NewRecorder()
		loginBody := `{"password": "hashed_password_1"}`
		req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(loginBody))
		req.Header.Set("Content-Type", "application/json")
		controllers.LoginHandler(w, req)

		testErrorResponse(t, w, http.StatusBadRequest, "Email is required", "EMAIL_REQUIRED")
	})

	t.Run("Login with missing password", func(t *testing.T) {
		w := httptest.NewRecorder()
		loginBody := `{"email": "user1@example.com"}`
		req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(loginBody))
		req.Header.Set("Content-Type", "application/json")
		controllers.LoginHandler(w, req)

		testErrorResponse(t, w, http.StatusBadRequest, "Password is required", "PASSWORD_REQUIRED")
	})

	t.Run("Login with user that does not exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		loginBody := `{"email": "dummy@example.com", "password": "hashed_password_1"}`
		req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(loginBody))
		req.Header.Set("Content-Type", "application/json")
		controllers.LoginHandler(w, req)

		testErrorResponse(t, w, http.StatusNotFound, "User not found", "USER_NOT_FOUND")
	})

	t.Run("Login with incorrect password", func(t *testing.T) {
		w := httptest.NewRecorder()
		loginBody := `{"email": "user1@example.com", "password": "wrong_password"}`
		req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(loginBody))
		req.Header.Set("Content-Type", "application/json")
		controllers.LoginHandler(w, req)

		testErrorResponse(t, w, http.StatusUnauthorized, "Wrong password", "WRONG_PASSWORD")
	})

	t.Run("Login with valid credentials", func(t *testing.T) {
		w := httptest.NewRecorder()
		loginBody := `{"email": "user1@example.com", "password": "hashed_password_1"}`
		req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(loginBody))
		req.Header.Set("Content-Type", "application/json")
		controllers.LoginHandler(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response controllers.LoginResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Token == "" {
			t.Fatalf("Expected a token in the response, got empty string")
		}
	})
}