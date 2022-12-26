package service

import (
	"ownify_api/internal/domain"
	"ownify_api/internal/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService[T domain.Userable] interface {
	GetUser(requestedUserID int64, userID int64) (*T, error)
	DeleteUser(id int64, userID int64) error
	UpdateUser(user T) (*T, error)
}

type userService[T domain.Userable] struct {
	dao repository.DBHandler[T]
}

func NewUserService[T domain.Userable](dao repository.DBHandler[T]) UserService[T] {
	return &userService[T]{dao: dao}
}

func (u *userService[T]) CreateUser(user T) error {
	_, err := u.dao.NewUserQuery().GetUser(0)
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

func (u *userService[T]) GetUser(requestedUserID int64, userID int64) (*T, error) {
	// var userBySession *domain.Person
	// var err error

	// userBySession, err = u.dao.NewUserQuery().GetPerson(userID)
	// if err != nil {
	// 	log.Printf("user isn't authorized %v", err)
	// }

	// userByRequest, err := u.dao.NewUserQuery().GetUser(requestedUserID)
	// if err != nil {
	// 	return nil, status.Errorf(codes.NotFound, "requested user doesn't exist: %v", err)
	// }

	// if userByRequest.ID == userBySession.ID || userBySession.Role == domain.ADMIN {
	// 	return userByRequest, nil
	// } else {
	// 	return &domain.Person{ID: userByRequest.ID, FirstName: userByRequest.FirstName, LastName: userByRequest.LastName}, nil
	// }
	return nil, nil
}

func (u *userService[T]) DeleteUser(id int64, userID int64) error {
	_, err := u.dao.NewUserQuery().GetUser(userID)
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
