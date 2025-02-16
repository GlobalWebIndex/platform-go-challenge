package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"platform-go-challenge/storage"
	"strings"
)

// Middleware to verify JWT token
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}
		// Expect "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}
		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		storedToken, exists := storage.GetToken(claims.UserID)
		if !exists || storedToken != tokenString {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalidated"})
			c.Abort()
			return
		}
		// Attach user ID to the context
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
