package services

import (
	"errors"

	"github.com/alexandros-georgantas/platform-go-challenge/internal/models"
	"github.com/alexandros-georgantas/platform-go-challenge/internal/serializers"
	"github.com/alexandros-georgantas/platform-go-challenge/internal/utils"
	"gorm.io/gorm"
)

type UserService interface {
	SignUp(su *serializers.SignUpUser) (*uint, error)
	Login(uc *serializers.UserCredentials) (*string, error)
	GetCurrentUser(uId uint) (*serializers.CurrentUserResponse, error)
}

type userService struct {
	db gorm.DB
}

func NewUserService(db gorm.DB) (UserService, error) {
	return &userService{db: db}, nil
}

func (us *userService) SignUp(su *serializers.SignUpUser) (*uint, error) {

	user := models.User{GivenName: su.GivenName, Surname: su.Surname, Email: su.Email, Password: su.Password}
	dbErr := us.db.Create(&user).Error

	if dbErr != nil {
		return nil, errors.New("could not create user")
	}

	return &user.ID, nil
}

func (us *userService) Login(uc *serializers.UserCredentials) (*string, error) {

	user := models.User{}

	dbErr := us.db.Where("email = ?", uc.Email).First(&user).Error

	if errors.Is(dbErr, gorm.ErrRecordNotFound) {
		return nil, errors.New("wrong user credentials")
	}

	if !utils.VerifyPassword(uc.Password, user.Password) {
		return nil, errors.New("wrong user credentials")
	}

	token, tErr := utils.CreateToken(user.ID)

	if tErr != nil {
		return nil, tErr
	}

	return token, nil
}

func (us *userService) GetCurrentUser(uId uint) (*serializers.CurrentUserResponse, error) {
	u := models.User{}
	cu := serializers.CurrentUserResponse{}

	dbErr := us.db.First(&u, uId).Error

	if errors.Is(dbErr, gorm.ErrRecordNotFound) {
		return nil, errors.New("invalid user id")
	}

	cu.Email = u.Email
	cu.ID = u.ID
	cu.GivenName = u.GivenName
	cu.Surname = u.Surname

	return &cu, nil
}
