package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	secretKey = []byte(os.Getenv("SERVER_SECRET"))
)

func CreateToken(uid uint) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": uid,
		"iss": "github.com/alexandros-georgantas/platform-go-challenge",
		"exp": time.Now().Add(time.Hour * 720).Unix(), // Expiration time in one month
		"iat": time.Now().Unix(),
	})

	tokenString, err := token.SignedString(secretKey)

	if err != nil || tokenString == "" {
		return nil, errors.New("error while creating JWT")
	}

	return &tokenString, nil
}

func VerifyToken(tokenString string) (int, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	claims, _ := token.Claims.(jwt.MapClaims)
	userId, cErr := strconv.Atoi(fmt.Sprint(claims["sub"]))

	if err != nil || cErr != nil {
		return -1, errors.New("error while parsing token")
	}

	if !token.Valid {
		return -1, errors.New("invalid token")
	}

	return userId, nil
}
