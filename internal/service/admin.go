package service

import (
	"ownify_api/internal/repository"
)

type AdminService interface {
	GrantBusiness(email string, isApproved bool) error
}

type adminService struct {
	dbHandler repository.DBHandler
}

// GrantBusiness implements AdminService
func (a *adminService) GrantBusiness(email string, isApproved bool) error {
	return a.dbHandler.NewAdminQuery().GrantBusiness(email, isApproved)
}

func NewAdminService(dbHandler repository.DBHandler) AdminService {
	return &adminService{dbHandler}
}
