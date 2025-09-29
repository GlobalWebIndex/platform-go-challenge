package authentication

import (
	"context"
	"net/http"
	"strings"

	"favourite_assets/server/services"
    "favourite_assets/server/errors"
	"github.com/Nerzal/gocloak/v13"
)

type contextKey string

const (
	UserInfoKey contextKey = "userInfo"
	RolesKey    contextKey = "roles"
)

func KeycloakAuth(kc *services.KeycloakService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				errors.WriteError(w, errors.ErrUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			userInfo, roles, err := kc.VerifyToken(r.Context(), token)
			if err != nil {
				errors.WriteError(w, errors.ErrInvalidToken)
				return
			}

			ctx := context.WithValue(r.Context(), UserInfoKey, userInfo)
			ctx = context.WithValue(ctx, RolesKey, roles)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserInfo from context
func GetUserInfo(ctx context.Context) *gocloak.UserInfo {
	if v := ctx.Value(UserInfoKey); v != nil {
		if u, ok := v.(*gocloak.UserInfo); ok {
			return u
		}
	}
	return nil
}

// GetRoles from context
func GetRoles(ctx context.Context) []string {
	if v := ctx.Value(RolesKey); v != nil {
		if r, ok := v.([]string); ok {
			return r
		}
	}
	return nil
}

// RequireRole ensures a user has a specific role
func RequireRole(ctx context.Context, role string) error {
	roles := GetRoles(ctx)
	if roles == nil {
		return errors.ErrUnauthorized
	}
	for _, r := range roles {
		if r == role {
			return nil
		}
	}
	return errors.ErrForbidden
}

// GetUserID helper (from Keycloak's "sub" claim) for future db authentication
func GetUserID(ctx context.Context) (string, error) {
	userInfo := GetUserInfo(ctx)
	if userInfo == nil || userInfo.Sub == nil {
		return "", http.ErrNoCookie
	}
	return *userInfo.Sub, nil
}
