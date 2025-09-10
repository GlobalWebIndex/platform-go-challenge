package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"platform-go-challenge/internal/middleware"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// TestNewJWT_Success tests the happy path for a standard user.
func TestNewJWT_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	secretKey := "test-secret"
	router.POST("/jwt", NewJWT(secretKey))

	// Act
	// The request body now only includes the userID field.
	requestBody := []byte(`{"userID": "user123"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/jwt", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["token"])

	// Validate the token and claims.
	token, err := jwt.Parse(response["token"], func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	assert.NoError(t, err)
	assert.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, "user123", claims["sub"])
	// The role should be "user" as determined by the server.
	assert.Equal(t, "user", claims["role"])
}

// TestNewJWT_AdminUserGeneratesTokenWithRole tests that a token for an admin user contains the correct role.
func TestNewJWT_AdminUserGeneratesTokenWithRole(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	secretKey := "test-secret"
	adminUserID := "admin_user_id"
	router.POST("/jwt", NewJWT(secretKey))

	// Act
	// The request body now only includes the admin userID field.
	requestBody := []byte(`{"userID": "` + adminUserID + `"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/jwt", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["token"])

	// Validate the token and claims.
	token, err := jwt.Parse(response["token"], func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	assert.NoError(t, err)
	assert.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, adminUserID, claims["sub"])
	// The role should be "admin" as determined by the server.
	assert.Equal(t, "admin", claims["role"])
}

func TestNewJWT_InvalidJSON(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	secretKey := "test-secret"
	router.Use(middleware.ErrorHandler())
	router.POST("/jwt", NewJWT(secretKey))

	// Act
	// The request body has an invalid JSON type for the userID field.
	requestBody := strings.NewReader(`{"userID": 123}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/jwt", requestBody)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response middleware.APIError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid json body", response.Message)
	assert.Equal(t, "User Id missing", response.Field)
}

func TestNewJWT_MissingUserID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	secretKey := "test-secret"
	router.Use(middleware.ErrorHandler())
	router.POST("/jwt", NewJWT(secretKey))

	// Act
	// The request body is missing the required userID field.
	requestBody := strings.NewReader(`{}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/jwt", requestBody)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response middleware.APIError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid json body", response.Message)
	assert.Equal(t, "User Id missing", response.Field)
}
