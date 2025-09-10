package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	JWT_USER_ID_KEY = "userID"
	JWT_ROLE_KEY    = "role"
)

// JWT is a middleware that handles JWT authentication.
func JWT(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header.
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, NewAPIError("Authorization header is required", "Authorization"))
			return
		}

		// Check if the token is in the correct format (Bearer <token>).
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, NewAPIError("Authorization header must be in the format 'Bearer <token>'", "Authorization"))
			return
		}

		// Parse and validate the JWT.
		tokenString := tokenParts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secretKey), nil
		})

		// Handle validation errors.
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, NewAPIError("Invalid or expired token", "Authorization"))
			return
		}

		// Extract claims and set user ID in the context.
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, NewAPIError("Failed to parse token claims", "token"))
			c.Abort()
			return
		}

		// Assuming the user ID is stored in a claim named "sub".
		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, NewAPIError("User ID not found in token claims", "sub"))
			c.Abort()
			return
		}

		// Extract the role from the claims.
		role, ok := claims["role"].(string)
		if !ok || role == "" {
			// Default to "user" role if not present, to prevent nil pointer issues.
			role = "user"
		}

		// Store userID in the context for downstream handlers.
		c.Set(JWT_USER_ID_KEY, userID)
		c.Set(JWT_ROLE_KEY, role)

		// Continue to the next handler.
		c.Next()
	}
}
