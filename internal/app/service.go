package app

import (
	"ownify_api/internal/service"
	desc "ownify_api/pkg"
)

type MicroserviceServer struct {
	desc.UnimplementedMicroserviceServer
	adminService     service.AdminService
	userService      service.UserService
	businessService  service.BusinessService
	ownershipService service.OwnershipService
	authService      service.AuthService
	tokenManager     service.TokenManager
	productService   service.ProductService
	walletService    service.WalletService
	notifyService    service.NotifyService
}

func NewMicroservice(
	adminService service.AdminService,
	userService service.UserService,
	businessService service.BusinessService,
	ownershipService service.OwnershipService,
	authService service.AuthService,
	tokenManager service.TokenManager,
	productService service.ProductService,
	walletService service.WalletService,
	notifyService service.NotifyService,
) *MicroserviceServer {
	return &MicroserviceServer{
		adminService:     adminService,
		userService:      userService,
		businessService:  businessService,
		ownershipService: ownershipService,
		authService:      authService,
		tokenManager:     tokenManager,
		productService:   productService,
		walletService:    walletService,
		notifyService:    notifyService,
	}
}
