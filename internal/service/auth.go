package service

import (
	"strconv"

	"ownify_api/internal/repository"
	//"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	CheckEmail(email string) bool
	ValidUser(pubKey string, idFingerprint string) bool
	ValidBusiness(uid string, email string) bool
	ValidUserWithId(userId string, pubKey string) bool
}

type authService struct {
	dbHandler    repository.DBHandler
	tokenManager TokenManager
}

func NewAuthService(dbHandler repository.DBHandler, tokenManager TokenManager) AuthService {
	return &authService{dbHandler: dbHandler, tokenManager: tokenManager}
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
func (a *authService) CheckEmail(email string) bool {
	_, err := a.dbHandler.NewBusinessQuery().GetBusiness(email)
	return err == nil
}

// ValidUser implements AuthService
func (a *authService) ValidUser(pubKey string, idFingerprint string) bool {
	_, err := a.dbHandler.NewUserQuery().GetUser(pubKey, idFingerprint)
	return err == nil
}

func (a *authService) ValidUserWithId(userId string, pubKey string) bool {
	_, err := a.dbHandler.NewUserQuery().VerifyUser(userId, pubKey)
	return err == nil
}

func (a *authService) ValidBusiness(uid string, email string) bool {
	_, err := a.dbHandler.NewBusinessQuery().VerifyBusiness(uid, email)
	return err == nil
}
