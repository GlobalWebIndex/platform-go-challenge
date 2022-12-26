package app

import (
	"ownify_api/internal/service"
	desc "ownify_api/pkg"
)

type MicroserviceServer struct {
	desc.UnimplementedMicroserviceServer
	userService         service.UserService
	authService         service.AuthService
	emailService        service.EmailService
	tokenManager        service.TokenManager
}

func NewMicroservice(
	userService service.UserService,
	authService service.AuthService,
	emailService service.EmailService,
	tokenManager service.TokenManager) *MicroserviceServer {
	return &MicroserviceServer{
		userService:         userService,
		authService:         authService,
		emailService:        emailService,
		tokenManager:        tokenManager,
	}
}
