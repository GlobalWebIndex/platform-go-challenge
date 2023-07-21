package service

import (
	"gwi_api/internal/dto"
	"gwi_api/internal/repository"
)

type UserService interface {
	CreateUser(
		user dto.UserDto) (*uint64, error)

	GetUser(userId string, pubKey string) (*dto.UserDto, error)

	DeleteUser(pubKey string) error
}

type userService struct {
	dbHandler repository.DBHandler
}

func NewUserService(dbHandler repository.DBHandler) UserService {
	return &userService{dbHandler}
}

func (u *userService) CreateUser(
	user dto.UserDto) (*uint64, error) {
	return u.dbHandler.NewUserQuery().CreateUser(user)
}

func (a *userService) GetUser(userId string, pubKey string) (*dto.UserDto, error) {
	return nil, nil
	//return a.dbHandler.NewUserQuery().VerifyUser(userId, pubKey)
}

func (u *userService) DeleteUser(pubKey string) error {
	return nil
	//return u.dbHandler.NewUserQuery().DeleteUser(pubKey)
}
