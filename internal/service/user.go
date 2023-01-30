package service

import (
	"log"
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
	GetLastUserId(walletType string) (*int64, error)
}

type userService struct {
	dbHandler repository.DBHandler
}

// GetUser implements UserService
func (*userService) GetUser(userID int64, walletType string) (*interface{}, error) {
	panic("unimplemented")
}

func NewUserService(dbHandler repository.DBHandler) UserService {
	return &userService{dbHandler}
}

func (u *userService) CreateUser(
	user dto.BriefUser) (*int64, error) {
	id, err := u.dbHandler.NewUserQuery().GetUserByBriefInfo(user)
	if id != nil {
		log.Println("[Waring] Already exist: userId")
		return id, nil
	}
	if err != nil {
		log.Println(err)
	}
	id, err = u.dbHandler.NewUserQuery().CreateUser(
		user,
	)

	if err != nil {
		return nil, err
	}

	return id, nil
}



func (u *userService) GetLastUserId(walletType string) (*int64, error) {
	userId, err := u.dbHandler.NewUserQuery().GetLastUserId(walletType)
	if err != nil {
		return nil, err
	}
	return userId, nil
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
