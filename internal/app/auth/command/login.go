package command

import (
	"fmt"
	"github.com/akazantzidis/gwi-ass/internal/domain/token"
	"github.com/akazantzidis/gwi-ass/internal/domain/user"
	"github.com/akazantzidis/gwi-ass/internal/pkg/helper"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest encapsulates username/password
type LoginRequest struct {
	Username string
	Password string
}

// LoginHandler interface
type LoginHandler interface {
	Handle(req LoginRequest) (accessToken string, refreshToken string, err error)
}

type loginHandler struct {
	userRepo    user.Repository // you can define a UserRepo interface
	refreshRepo token.RefreshRepository
}

func NewLoginHandler(userRepo user.Repository, refreshRepo token.RefreshRepository) LoginHandler {
	return &loginHandler{userRepo: userRepo, refreshRepo: refreshRepo}
}

func (h *loginHandler) Handle(req LoginRequest) (string, string, error) {
	u, err := h.userRepo.GetByUsername(req.Username)
	if err != nil {
		return "", "", fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return "", "", fmt.Errorf("invalid credentials")
	}

	access, err := helper.GenerateAccessToken(u.ID.String(), u.Roles)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refresh, exp, err := helper.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	h.refreshRepo.Save(refresh, token.RefreshRecord{
		UserID: u.ID,
		Expiry: exp,
		Roles:  u.Roles,
	})

	return access, refresh, nil
}
