package models

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// User model
//
// for both admins and plain users
//
// swagger:model User
type User struct {
	// ID of the User
	// in: string
	// example: 669c34226029d2ef83fc38f8
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	// Username of the User
	Username string `json:"username"`
	// Hashed Password of the User
	Password string `json:"password,omitempty"`
	// CreatedAt is the date-time of user's creation
	CreatedAt time.Time `json:"created_at" bson:"created_at" `
	// ModifiedAt is the date-time of user's last modification
	ModifiedAt time.Time `json:"modified_at" bson:"modified_at" `
	// Role of the User
	Role string `json:"role,omitempty"`
}

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}
