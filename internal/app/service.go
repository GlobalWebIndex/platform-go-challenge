package app

import (
	"ownify_api/internal/domain"
	"ownify_api/internal/service"
	desc "ownify_api/pkg"
)

type MicroserviceServer[T domain.Userable] struct {
	desc.UnimplementedMicroserviceServer
	userService         service.UserService[T]
	authService         service.AuthService[T]
	tokenManager        service.TokenManager
}

func NewMicroservice[T domain.Userable](
	userService service.UserService[T],
	authService service.AuthService[T],
	tokenManager service.TokenManager) *MicroserviceServer[T] {
	return &MicroserviceServer[T]{
		userService:         userService,
		authService:         authService,
		tokenManager:        tokenManager,
	}
}
