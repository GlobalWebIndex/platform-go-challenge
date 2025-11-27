package command

import (
	"github.com/akazantzidis/gwi-ass/internal/domain/token"
)

type LogoutHandler interface {
	Handle(refreshToken string) error
}

type logoutHandler struct {
	refreshRepo token.RefreshRepository
}

func NewLogoutHandler(rr token.RefreshRepository) LogoutHandler {
	return &logoutHandler{refreshRepo: rr}
}

func (h *logoutHandler) Handle(refreshToken string) error {
	h.refreshRepo.Delete(refreshToken)
	return nil
}
