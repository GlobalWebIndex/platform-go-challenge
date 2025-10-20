package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("super-secret-key") // store in env variable in production

func GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,                                // user ID
		"exp": time.Now().Add(24 * time.Hour).Unix(), // expires in 24h
		"iat": time.Now().Unix(),                     // issued at
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
