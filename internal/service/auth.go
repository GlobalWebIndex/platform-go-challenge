package service

import (
	"strconv"

	"ownify_api/internal/repository"
	//"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	//SignUp(user T) (*int64, error)
	SignIn(email, password string) (*string, error)
	SignInWithPhone(firebaseToken string, userId int64) (*string, error)
	Logout(userID int64) error
}

type authService struct {
	dbHandler    repository.DBHandler
	tokenManager TokenManager
}

func NewAuthService(dbHandler repository.DBHandler, tokenManager TokenManager) AuthService {
	return &authService{dbHandler: dbHandler, tokenManager: tokenManager}
}

// func (a *authService[T]) SignUp(user T) (*int64, error) {
// 	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// user.Password = string(hashedPassword)
// 	briefUser := dto.BriefUser{0, "", ""}
// 	id, err := a.dao.NewUserQuery().CreateUser(briefUser)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return id, nil
// }

func (a *authService) SignIn(email, reqPassword string) (*string, error) {
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

func (a *authService) SignInWithPhone(firebaseToken string, userId int64) (*string, error) {
	token, err := a.tokenManager.NewFirebaseToken(firebaseToken, userId)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (a *authService) Logout(userID int64) error {
	_, err := a.tokenManager.NewJWT(strconv.Itoa(int(userID)))
	if err != nil {
		return err
	}
	return nil
}
