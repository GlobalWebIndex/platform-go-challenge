package models

type User struct {
	ID       int    `gorm:"primary_key;"`
	Email    string `gorm:"not null;"`
	Password string `gorm:"not null;"`
}
