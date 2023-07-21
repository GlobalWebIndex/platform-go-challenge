package service

import (
	"fmt"
	"strconv"

	"gwi_api/internal/dto"
	"gwi_api/internal/repository"

	//"gwi_api/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignUp(user dto.UserDto) (*int64, *string, error)
	SignIn(email, password string) (*string, error)
	GetUserID(token string) (*int64, error)
	Logout(userID int64) error
}

type authService struct {
	db           repository.DBHandler
	tokenManager TokenManager
}

func NewAuthService(db repository.DBHandler, tokenManager TokenManager) AuthService {
	return &authService{db: db, tokenManager: tokenManager}
}

// GetUserID implements AuthService.
func (a *authService) GetUserID(token string) (*int64, error) {
	return a.tokenManager.Parse(token)
}

func (a *authService) SignUp(user dto.UserDto) (*int64, *string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return nil, nil, err
	}
	user.Password = string(hashedPassword)
	id, err := a.db.NewUserQuery().CreateUser(user)
	if err != nil {
		return nil, nil, err
	}
	userId := int64(*id)
	token, err := a.tokenManager.NewJWT(strconv.Itoa(int(userId)))

	if err != nil {
		return nil, nil, err
	}
	return &userId, &token, nil
}

func (a *authService) SignIn(email, reqPassword string) (*string, error) {
	password, err := a.db.NewUserQuery().GetUserPasswordByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(*password), []byte(reqPassword))
	if err != nil {
		return nil, fmt.Errorf("passwords don't match %v", err)
	} else {
		userID, err := a.db.NewUserQuery().GetUserIdByEmail(email)
		if err != nil {
			return nil, err
		}

		jwt, err := a.tokenManager.NewJWT(strconv.Itoa(int(*userID)))
		if err != nil {
			return nil, err
		}

		return &jwt, nil
	}
}

func (a *authService) Logout(userID int64) error {
	_, err := a.tokenManager.NewJWT(strconv.Itoa(int(userID)))
	if err != nil {
		return err
	}
	return nil
}
