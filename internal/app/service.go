package app

import (
	"ownify_api/internal/service"
	desc "ownify_api/pkg"
)

type MicroserviceServer struct {
	desc.UnimplementedMicroserviceServer
	userService    service.UserService
	authService    service.AuthService
	tokenManager   service.TokenManager
	productService service.ProductService
}

func NewMicroservice(
	userService service.UserService,
	authService service.AuthService,
	tokenManager service.TokenManager,
	productService service.ProductService,
) *MicroserviceServer {
	return &MicroserviceServer{
		userService:    userService,
		authService:    authService,
		tokenManager:   tokenManager,
		productService: productService,
	}
}
