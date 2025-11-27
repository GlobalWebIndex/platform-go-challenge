package auth

import (
	"encoding/json"
	"fmt"
	"github.com/akazantzidis/gwi-ass/internal/app"
	"github.com/akazantzidis/gwi-ass/internal/app/auth/command"
	"github.com/akazantzidis/gwi-ass/internal/pkg/helper"
	"net/http"
	"time"
)

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type AuthHandler struct {
	authServices app.AuthServices
}

type RefreshRecord struct {
	UserID string
	Expiry time.Time
	Roles  []string
}

func NewAuthHandler(app app.AuthServices) *AuthHandler {
	return &AuthHandler{authServices: app}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var cred Credential
	if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
		helper.WriteJSONError(w, http.StatusBadRequest, err, "invalid payload")
		return
	}

	access, refresh, err := h.authServices.Commands.LoginUserHandler.Handle(command.LoginRequest{
		Username: cred.Username,
		Password: cred.Password,
	})
	if err != nil {
		helper.WriteJSONError(w, http.StatusUnauthorized, err, "invalid username or password")
		return
	}

	resp := TokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresAt:    time.Now().Add(helper.AccessTokenTTL),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Refresh - exchanges valid refresh token for new access + refresh (rotate refresh token)
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var p struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		helper.WriteJSONError(w, http.StatusBadRequest, err, "invalid payload")
		return
	}

	access, newRefresh, err := h.authServices.Commands.RefreshTokenUserHandler.Handle(command.RefreshRequest{
		RefreshToken: p.RefreshToken,
	})
	if err != nil {
		helper.WriteJSONError(w, http.StatusUnauthorized, err, nil)
		return
	}

	resp := TokenResponse{
		AccessToken:  access,
		RefreshToken: newRefresh,
		ExpiresAt:    time.Now().Add(helper.AccessTokenTTL),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var p struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		helper.WriteJSONError(w, http.StatusBadRequest, err, "invalid payload")
		return
	}

	if p.RefreshToken == "" {
		helper.WriteJSONError(w, http.StatusBadRequest, fmt.Errorf("missing refresh_token"), nil)
		return
	}

	err := h.authServices.Commands.LogoutUserHandler.Handle(p.RefreshToken)
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err, nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
