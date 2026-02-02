package auth

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const ContextKeySubject contextKey = "auth_subject"

type JWTConfig struct {
	PublicKey *rsa.PublicKey
	Issuer    string
	Audience  string
}

type Claims struct {
	jwt.RegisteredClaims
}

func SubjectFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(ContextKeySubject)
	s, ok := v.(string)
	return s, ok
}

func Middleware(cfg JWTConfig) func(http.Handler) http.Handler {
	if cfg.PublicKey == nil {
		panic("RSA public key is required")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr, err := readBearerToken(r.Header.Get("Authorization"))
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
				if t.Method != jwt.SigningMethodRS256 {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return cfg.PublicKey, nil
			}, jwt.WithLeeway(30*time.Second))

			if err != nil || !token.Valid {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if cfg.Issuer != "" && claims.Issuer != cfg.Issuer {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			if cfg.Audience != "" && !audienceContains(claims.Audience, cfg.Audience) {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			sub := strings.TrimSpace(claims.Subject)
			if sub == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeySubject, sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func LoadRSAPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("invalid PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}
	return rsaPub, nil
}

func readBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("missing Authorization")
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid Authorization")
	}
	t := strings.TrimSpace(parts[1])
	if t == "" {
		return "", errors.New("empty token")
	}
	return t, nil
}

func audienceContains(aud []string, expected string) bool {
	for _, a := range aud {
		if a == expected {
			return true
		}
	}
	return false
}
