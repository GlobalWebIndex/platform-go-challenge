package service

import (
	"ownify_api/internal/dto"
	"ownify_api/internal/repository"
)

type UserService interface {
	CreateUser(
		user dto.BriefUser,
	) error

	GetUser(userId string, pubKey string) (*dto.BriefUser, error)

	DeleteUser(pubKey string) error
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

func (a *userService) GetUser(userId string, pubKey string) (*dto.BriefUser, error) {
	return a.dbHandler.NewUserQuery().VerifyUser(userId, pubKey)
}

func (u *userService) DeleteUser(pubKey string) error {
	return u.dbHandler.NewUserQuery().DeleteUser(pubKey)
}
