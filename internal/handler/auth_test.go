package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"platform-go-challenge/internal/auth"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	config := auth.Config{
		SecretKey:  "test-secret-key",
		ExpiryTime: time.Hour * 24,
	}
	jwtManager := auth.NewJWTManager(config)

	t.Run("successfully generates token", func(t *testing.T) {
		payload := map[string]string{"user_id": "user-123"}
		data, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(data))
		w := httptest.NewRecorder()

		Login(jwtManager)(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotEmpty(t, response["token"])
		assert.Equal(t, "user-123", response["user_id"])
	})

	t.Run("rejects invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte("invalid")))
		w := httptest.NewRecorder()

		Login(jwtManager)(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "Invalid JSON")
	})

	t.Run("rejects empty user_id", func(t *testing.T) {
		payload := map[string]string{"user_id": ""}
		data, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(data))
		w := httptest.NewRecorder()

		Login(jwtManager)(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "user_id is required")
	})

	t.Run("handles token generation error", func(t *testing.T) {
		// Create JWT manager with empty secret to trigger error
		badConfig := auth.Config{
			SecretKey:  "",
			ExpiryTime: time.Hour,
		}
		badJWT := auth.NewJWTManager(badConfig)

		payload := map[string]string{"user_id": "user-123"}
		data, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(data))
		w := httptest.NewRecorder()

		Login(badJWT)(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "Failed to generate token")
	})
}