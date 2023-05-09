package handlers

import (
	"github.com/loukaspe/platform-go-challenge/internal/core/domain"
	"net/http"
	"strings"
)

type AuthenticationMw struct {
	claimsDomain domain.JwtClaimsInterface
}
type AuthenticationMechanismInterface interface {
	AuthenticationMW(next http.Handler) http.Handler
}

func NewAuthenticationMw(claims domain.JwtClaimsInterface) *AuthenticationMw {
	return &AuthenticationMw{claimsDomain: claims}
}
func (a *AuthenticationMw) AuthenticationMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer") {
			http.Error(w, "Not Authorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := a.claimsDomain.GetClaimsFromToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		r = r.WithContext(a.claimsDomain.SetJWTClaimsContext(r.Context(), claims))
		next.ServeHTTP(w, r)
	})
}
