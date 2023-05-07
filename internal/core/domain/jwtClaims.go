package domain

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
)

type JwtClaims struct {
	*jwt.RegisteredClaims
	UserInfo interface{}
}
type JwtClaimsInterface interface {
	CreateToken(sub string, userInfo interface{}) (string, error)
	GetClaimsFromToken(tokenString string) (jwt.MapClaims, error)
	SetJWTClaimsContext(ctx context.Context, claims jwt.MapClaims) context.Context
}
