package domain

import "time"

type User struct {
	Id        string
	Name      string
	Email     string
	Password  string // hashed
	CreatedAt time.Time
	UpdatedAt time.Time
}
