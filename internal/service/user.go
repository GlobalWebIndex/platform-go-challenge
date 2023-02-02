package service

import (
	"ownify_api/internal/dto"
	"ownify_api/internal/repository"
)

type UserService interface {
	CreateUser(
		user dto.BriefUser,
	) error

	GetUser(pubKey string) (*dto.BriefUser, error)
	DeleteUser(pubKey string) error
	GetLastUserId(walletType string) (*int64, error)
}

type userService struct {
	dbHandler repository.DBHandler
}

func NewUserService(dbHandler repository.DBHandler) UserService {
	return &userService{dbHandler}
}

func (u *userService) CreateUser(
	user dto.BriefUser) error {
	return u.dbHandler.NewUserQuery().CreateUser(user)
}

// GetUser implements UserService
func (u *userService) GetUser(pubKey string) (*dto.BriefUser, error) {
	return u.dbHandler.NewUserQuery().GetUser(pubKey, "")
}

func (u *userService) GetLastUserId(walletType string) (*int64, error) {
	return u.dbHandler.NewUserQuery().GetLastUserId(walletType)
}

func (u *userService) DeleteUser(pubKey string) error {
	return u.dbHandler.NewUserQuery().DeleteUser(pubKey)
}
