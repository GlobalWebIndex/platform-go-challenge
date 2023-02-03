package app

import (
	"ownify_api/internal/service"
	desc "ownify_api/pkg"
)

type MicroserviceServer struct {
	desc.UnimplementedMicroserviceServer
	userService      service.UserService
	businessService  service.BusinessService
	ownershipService service.OwnershipService
	authService      service.AuthService
	tokenManager     service.TokenManager
	productService   service.ProductService
	walletService    service.WalletService
}

func NewMicroservice(
	userService service.UserService,
	businessService service.BusinessService,
	ownershipService service.OwnershipService,
	authService service.AuthService,
	tokenManager service.TokenManager,
	productService service.ProductService,
	walletService service.WalletService,
) *MicroserviceServer {
	return &MicroserviceServer{
		userService:      userService,
		businessService:  businessService,
		ownershipService: ownershipService,
		authService:      authService,
		tokenManager:     tokenManager,
		productService:   productService,
		walletService:    walletService,
	}
}
