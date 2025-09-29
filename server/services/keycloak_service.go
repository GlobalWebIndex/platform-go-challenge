package services

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/golang-jwt/jwt/v5"
	"favourite_assets/server/errors"
)

type KeycloakService struct{}

func NewKeycloakService() *KeycloakService {
	return &KeycloakService{}
}
func (k *KeycloakService) VerifyToken(ctx context.Context, token string) (*gocloak.UserInfo, []string, error) {
	parsed, _, err := jwt.NewParser().ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return nil, nil, errors.ErrBadRequest
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, errors.ErrUnauthorized
	}

	// Build user info from claims
	userInfo := &gocloak.UserInfo{
		Sub:               getClaimString(claims, "sub"),
		PreferredUsername: getClaimString(claims, "preferred_username"),
		Email:             getClaimString(claims, "email"),
		Name:              getClaimString(claims, "name"),
	}

	roles := []string{}

	// Realm roles
	if realmAccess, ok := claims["realm_access"].(map[string]interface{}); ok {
		if r, ok := realmAccess["roles"].([]interface{}); ok {
			for _, role := range r {
				if roleStr, ok := role.(string); ok {
					roles = append(roles, roleStr)
				}
			}
		}
	}

	return userInfo, roles, nil
}

func getClaimString(claims jwt.MapClaims, key string) *string {
	if v, ok := claims[key].(string); ok {
		return &v
	}
	return nil
}
