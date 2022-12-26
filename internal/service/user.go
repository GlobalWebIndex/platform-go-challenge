package service

import (
	"fmt"
	"ownify_api/internal/domain"
	"ownify_api/internal/dto"
	"ownify_api/internal/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService[T domain.Userable] interface {
	CreateUser(
		user dto.BriefUser,
	) (*int64, error)
	GetUser(requestedUserID int64, userID int64) (*T, error)
	DeleteUser(id int64, userID int64) error
	UpdateUser(user T) (*T, error)
}

type userService[T domain.Userable] struct {
	dbHandler repository.DBHandler[T]
}

func NewUserService[T domain.Userable](dbHandler repository.DBHandler[T]) UserService[T] {
	return &userService[T]{dbHandler}
}

func (u *userService[T]) CreateUser(
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

func (u *userService[T]) GetUser(requestedUserID int64, userID int64) (*T, error) {

	return nil, nil
}

func (u *userService[T]) DeleteUser(id int64, userID int64) error {
	_, err := u.dbHandler.NewUserQuery().GetUser(userID)
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

func (u *userService[T]) UpdateUser(user T) (*T, error) {
	// email checking
	// phone number checking
	//_, err := u.dao.NewUserQuery().GetUser(person.ID)
	// if err != nil {
	// 	return nil, err
	// }

	// if user.Role == domain.ADMIN || user.ID == person.ID {
	// 	updatedUser, err := u.dao.NewUserQuery().UpdateUser(person)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return updatedUser, nil
	// }
	return nil, status.Errorf(codes.PermissionDenied, "you don't have access")
}
