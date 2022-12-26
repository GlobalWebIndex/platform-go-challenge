package service

import (
	"fmt"
	"ownify_api/internal/dto"
	"ownify_api/internal/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService interface {
	CreateUser(
		user dto.BriefUser,
	) (*int64, error)
	GetUser(userID int64, walletType string) (*interface{}, error)
	DeleteUser(userID int64, walletType string) error
}

type userService struct {
	dbHandler repository.DBHandler
}

func NewUserService(dbHandler repository.DBHandler) UserService {
	return &userService{dbHandler}
}

func (u *userService) CreateUser(
	user dto.BriefUser) (*int64, error) {

	id, err := u.dbHandler.NewUserQuery().GetUserByBriefInfo(user)
	if id != nil {
		return nil, fmt.Errorf("[ERR] this user already exist: id%s", id)
	}
	id, err = u.dbHandler.NewUserQuery().CreateUser(
		dto.BriefUser{
			ChainId: 0, Wallet: "", WalletType: "",
		},
	)

	if err != nil {
		return nil, err
	}

	return id, nil
}

func (u *userService) GetUser(userID int64, walletType string) (*interface{}, error) {

	return nil, nil
}

func (u *userService) DeleteUser(userID int64, walletType string) error {
	_, err := u.dbHandler.NewUserQuery().GetUser(userID, walletType)
	if err != nil {
		return err
	}

	// if user.Role == domain.ADMIN || id == user.ID {
	// 	err = u.dao.NewUserQuery().DeleteUser(id)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return nil
	// }
	return status.Errorf(codes.PermissionDenied, "you have no access")
}

// func (u *userService) UpdateUser(user T) (*T, error) {
// 	// email checking
// 	// phone number checking
// 	//_, err := u.dao.NewUserQuery().GetUser(person.ID)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// if user.Role == domain.ADMIN || user.ID == person.ID {
// 	// 	updatedUser, err := u.dao.NewUserQuery().UpdateUser(person)
// 	// 	if err != nil {
// 	// 		return nil, err
// 	// 	}
// 	// 	return updatedUser, nil
// 	// }
// 	return nil, status.Errorf(codes.PermissionDenied, "you don't have access")
// }
