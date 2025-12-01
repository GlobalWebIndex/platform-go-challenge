package models

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"platform-go-challenge/internal/common"
)

// User role type
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

func isValidRole(r string) bool {
	return r == string(RoleUser) || r == string(RoleAdmin)
}

// User struct to represent each user
type User struct {
	ID       uuid.UUID
	Username string
	Password string
	UserRole Role
}

func CreateUser(username string, password string, role Role) (User, error) {
	// Validate inputs
	if username == "" {
		return User{}, errors.New(common.ErrEmptyUsername)
	}
	if password == "" {
		return User{}, errors.New(common.ErrEmptyPassword)
	}
	if !isValidRole(string(role)) {
		return User{}, errors.New(common.ErrInvalidRole)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, errors.New(common.ErrHashPassword)
	}

	// we do not create a user ID here; it will be created when adding to storage
	user := User{
		Username: username,
		Password: string(hashed),
		UserRole: role,
	}

	return user, nil
}
