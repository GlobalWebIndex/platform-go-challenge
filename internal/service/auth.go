package service

import (
	"strconv"

	"ownify_api/internal/domain"
	"ownify_api/internal/repository"

	//"golang.org/x/crypto/bcrypt"
)

type AuthService[T domain.Userable] interface {
	SignUp(user T) (*int64, error)
	SignIn(email, password string) (*string, error)
	Logout(userID int64) error
}

type authService[T domain.Userable] struct {
	dao          repository.DBHandler[T]
	tokenManager TokenManager
}

func NewAuthService[T domain.Userable](dao repository.DBHandler[T], tokenManager TokenManager) AuthService[T] {
	return &authService[T]{dao: dao, tokenManager: tokenManager}
}

func (a *authService[T]) SignUp(user T) (*int64, error) {
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	// if err != nil {
	// 	return nil, err
	// }
	// user.Password = string(hashedPassword)
	id, err := a.dao.NewUserQuery().CreateUser(user)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (a *authService[T]) SignIn(email, reqPassword string) (*string, error) {
	// password, err := a.dao.NewUserQuery().GetUserPasswordByEmail(email)
	// if err != nil {
	// 	return nil, err
	// }

	// err = bcrypt.CompareHashAndPassword([]byte(*password), []byte(reqPassword))
	// if err != nil {
	// 	return nil, fmt.Errorf("passwords don't match %v", err)
	// } else {
	// 	userID, err := a.dao.NewUserQuery().GetUserIdByEmail(email)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	jwt, err := a.tokenManager.NewJWT(strconv.Itoa(int(*userID)))
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	return &jwt, nil
	// }
	return nil, nil
}

func (a *authService[T]) Logout(userID int64) error {
	_, err := a.tokenManager.NewJWT(strconv.Itoa(int(userID)))
	if err != nil {
		return err
	}
	return nil
}
