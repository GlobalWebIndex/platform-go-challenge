package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string    `json:"id" binding:"required"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewUser() (*User, error) {
	newID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &User{
		ID:        newID.String(),
		CreatedAt: time.Now(),
	}, nil
}

func (user *User) GetID() string {
	return user.ID
}

func (user *User) GetCreatedAt() time.Time {
	return user.CreatedAt
}
