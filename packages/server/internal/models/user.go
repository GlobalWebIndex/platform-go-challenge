package models

import (
	"errors"

	"github.com/alexandros-georgantas/platform-go-challenge/internal/utils"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string `gorm:"unique;not null"`
	GivenName string
	Surname   string
	Password  string
	Favorites []Favorite
}

type CurrentUserResponse struct {
	ID        uint
	Email     string
	GivenName string
	Surname   string
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {

	hash, err := utils.HashPassword(u.Password)

	if err != nil {
		return errors.New("problem when tried to hash password")
	}

	u.Password = hash

	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {

	if u.Password != "" {

		hash, err := utils.HashPassword(u.Password)

		if err != nil {
			return errors.New("problem when tried to hash password")
		}

		u.Password = hash
	}
	return
}
