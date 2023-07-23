package app

import (
	"gwi_api/internal/service"
	desc "gwi_api/pkg"
)

type MicroserviceServer struct {
	desc.UnimplementedMicroserviceServer
	authService     service.AuthService
	userService     service.UserService
	assetService    service.AssetService
	favoriteService service.FavoritesService
}

func NewMicroservice(
	authService service.AuthService,
	userService service.UserService,
	assetService service.AssetService,
	favoriteService service.FavoritesService,

) *MicroserviceServer {
	return &MicroserviceServer{
		authService:     authService,
		userService:     userService,
		assetService:    assetService,
		favoriteService: favoriteService,
	}
}
