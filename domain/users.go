package domain

import (
	"time"
)

type User struct {
	Key       string    `json:"_key,omitempty"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRepository interface {
	CreateUsers() ([]User, error)
	GetUser(userID string) (string, error)
}
