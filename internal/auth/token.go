package auth

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func GetUserIDFromToken(c *gin.Context) (uuid.UUID, error) {
	header := c.GetHeader("Authorization")
	if header == "" {
		return uuid.Nil, errors.New("no authorization header")
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return uuid.Nil, errors.New("no authorization header")
	}

	bearerToken := parts[1]

	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(bearerToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("incorrect signing method used for jwt token")
		}
		// Key should come from some safe storage
		return []byte("12345678123456781234567812345678"), nil
	})

	if err != nil {
		// log details about error
		return uuid.Nil, err
	}

	if claimValue, ok := claims["userid"]; !ok {
		return uuid.Nil, errors.New("userid not present in jwt token claims")
	} else {
		s, ok := claimValue.(string)
		if !ok {
			return uuid.Nil, errors.New("cannot parse userid claim from jwt token claims as string")
		}
		userId, err := uuid.Parse(s)
		if err != nil {
			return uuid.Nil, errors.New("cannot parse userid from jwt token claims as uuid")
		}
		return userId, nil
	}
}
