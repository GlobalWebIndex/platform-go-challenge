package repository

import (
	"errors"
	"strings"

	"platform-go-challenge/models"
	"platform-go-challenge/utils"

	"gorm.io/gorm"
)

type UserRepository struct{}

func (u UserRepository) AddUser(db *gorm.DB, user models.User) (int, error) {
	err := db.Create(&user).Error
	if err != nil {
		if strings.Contains(err.Error(), utils.UniqueConstrainViolationString) {
			return 0, utils.NewUniqueConstrainViolation("Email already exists")
		}

		return 0, err
	}
	return user.ID, nil
}

func (u UserRepository) LoginUser(db *gorm.DB, user models.User, email string) (models.User, error) {
	err := utils.DB.Where("email = ?", email).Find(&user).Error
	errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil {
		return user, err
	}
	return user, nil
}
