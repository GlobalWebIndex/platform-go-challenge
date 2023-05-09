package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/loukaspe/platform-go-challenge/internal/core/domain"
	"time"
)

type claimskey int

var claimsKey claimskey

type AuthMechanism struct {
	secret        []byte
	signingMethod string
}

func NewAuthMechanism(
	secret string,
	signingMethod string,
) *AuthMechanism {
	return &AuthMechanism{
		secret:        []byte(secret),
		signingMethod: signingMethod,
	}
}
func (j *AuthMechanism) CreateToken(sub string, userInfo interface{}) (string, error) {
	token := jwt.New(jwt.GetSigningMethod(j.signingMethod))
	expiration := time.Now().Add(time.Hour)
	token.Claims = &domain.JwtClaims{
		&jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			Subject:   sub,
		}, userInfo,
	}
	val, err := token.SignedString(j.secret)
	if err != nil {

		return "", err
	}
	return val, nil
}
func (j *AuthMechanism) GetClaimsFromToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
func (j *AuthMechanism) SetJWTClaimsContext(ctx context.Context, claims jwt.MapClaims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}
