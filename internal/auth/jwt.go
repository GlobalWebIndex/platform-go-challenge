package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"platform-go-challenge/internal/config"
	"platform-go-challenge/internal/models"
)

// Claims stored inside the JWT
type Claims struct {
	UserID uuid.UUID   `json:"user_id"`
	Role   models.Role `json:"role"`
	jwt.RegisteredClaims
}

// JWTGenerator is an interface used by handlers
type JWTGenerator interface {
	Generate(cfg *config.Config, userID uuid.UUID, role models.Role) (string, error)
}

// Real implementation used in production
type RealJWTGenerator struct{}

// Generate creates a JWT token with user ID and role
func (g *RealJWTGenerator) Generate(cfg *config.Config, userID uuid.UUID, role models.Role) (string, error) {

	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.TokenExpiry)),
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// FIX: Ensure secret is []byte
	return token.SignedString([]byte(cfg.JWTSecret))
}
