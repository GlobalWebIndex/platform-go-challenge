package command

import (
	"fmt"
	"github.com/akazantzidis/gwi-ass/internal/pkg/helper"
	"time"

	"github.com/akazantzidis/gwi-ass/internal/domain/token"
)

type RefreshRequest struct {
	RefreshToken string
}

type RefreshHandler interface {
	Handle(req RefreshRequest) (string, string, error)
}

type refreshHandler struct {
	refreshRepo token.RefreshRepository
}

func NewRefreshHandler(refreshRepo token.RefreshRepository) RefreshHandler {
	return &refreshHandler{refreshRepo: refreshRepo}
}

func (h *refreshHandler) Handle(req RefreshRequest) (string, string, error) {
	if req.RefreshToken == "" {
		return "", "", fmt.Errorf("missing refresh token")
	}

	rec, ok := h.refreshRepo.Get(req.RefreshToken)
	if !ok {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	if time.Now().After(rec.Expiry) {
		h.refreshRepo.Delete(req.RefreshToken)
		return "", "", fmt.Errorf("refresh token expired")
	}

	// ROTATE TOKEN
	h.refreshRepo.Delete(req.RefreshToken)

	access, err := helper.GenerateAccessToken(rec.UserID.String(), rec.Roles)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefresh, exp, err := helper.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	h.refreshRepo.Save(newRefresh, token.RefreshRecord{
		UserID: rec.UserID,
		Expiry: exp,
		Roles:  rec.Roles,
	})

	return access, newRefresh, nil
}
