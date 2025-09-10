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

func TestNewJWT_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	secretKey := "test-secret"
	router.POST("/jwt", NewJWT(secretKey))

	// Act
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

	// Assert the token exists and is valid.
	assert.NotEmpty(t, response["token"])

	// Further validation of the signed token.
	token, _ := jwt.Parse(response["token"], func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	assert.True(t, token.Valid)
}

func TestNewJWT_InvalidJSON(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	secretKey := "test-secret"
	router.Use(middleware.ErrorHandler())
	router.POST("/jwt", NewJWT(secretKey))

	// Act
	requestBody := strings.NewReader(`{"userID": 123}`) // Invalid JSON type
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
