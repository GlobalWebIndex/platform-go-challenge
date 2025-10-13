package handler

import (
	"encoding/json"
	"net/http"

	"platform-go-challenge/internal/auth"
	"platform-go-challenge/internal/util"
)

// Login generates a JWT token for a user
func Login(jwtManager *auth.JWTManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UserID string `json:"user_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			util.RespondError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		if req.UserID == "" {
			util.RespondError(w, http.StatusBadRequest, "user_id is required")
			return
		}

		// Generate token
		token, err := jwtManager.GenerateToken(req.UserID)
		if err != nil {
			util.RespondError(w, http.StatusInternalServerError, "Failed to generate token")
			return
		}

		util.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"token":   token,
			"user_id": req.UserID,
		})
	}
}