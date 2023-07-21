package app

import (
	"gwi_api/internal/service"
	desc "gwi_api/pkg"
)

type MicroserviceServer struct {
	desc.UnimplementedMicroserviceServer
	authService service.AuthService
	userService service.UserService
}

func NewMicroservice(
	authService service.AuthService,
	userService service.UserService,

) *MicroserviceServer {
	return &MicroserviceServer{
		authService: authService,
		userService: userService,
	}
}
