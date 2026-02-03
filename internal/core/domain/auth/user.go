package auth

import (
	"errors"
)

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

func (u User) Validate() error {
	if u.ID == "" {
		return errors.New("id is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	return nil
}
