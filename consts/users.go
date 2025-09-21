package consts

import (
	"platform-go-challenge/domain"
	"time"
)

var Users = []domain.User{
	domain.User{
		FirstName: "William",
		LastName:  "Smith",
		Username:  "wsmith",
		Email:     "wsmith@outlook.com",
		CreatedAt: time.Now(),
	},
	domain.User{
		FirstName: "John",
		LastName:  "Doe",
		Username:  "jdoe",
		Email:     "jdoe@outlook.com",
		CreatedAt: time.Now(),
	},
	domain.User{
		FirstName: "Alice",
		LastName:  "Johnson",
		Username:  "ajohnson",
		Email:     "ajohnson@outlook.com",
		CreatedAt: time.Now(),
	},
}
