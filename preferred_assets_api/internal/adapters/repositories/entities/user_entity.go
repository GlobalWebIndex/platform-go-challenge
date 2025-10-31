package entities

import "time"

type UserEntity struct {
	Id        string
	Name      string
	Email     string
	Password  string // hashed
	CreatedAt time.Time
	UpdatedAt time.Time
}
