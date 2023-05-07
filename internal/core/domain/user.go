package domain

import "golang.org/x/crypto/bcrypt"

type User struct {
	Username        string
	Password        string
	FavouriteAssets []Asset
}

const HashCost = 10

func (user *User) CheckPassword(providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
}

func (user *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), HashCost)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}
