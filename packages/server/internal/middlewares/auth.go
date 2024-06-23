package middlewares

import (
	"net/http"

	"github.com/alexandros-georgantas/platform-go-challenge/internal/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) {
	// var user models.User

	tokenString := c.GetHeader("Authorization")

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is missing"})
		c.Abort()
		return
	}

	tokenString = tokenString[len("Bearer "):]

	userId, vErr := utils.VerifyToken(tokenString)

	if vErr != nil || userId == -1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		c.Abort()
		return
	}

	c.Set("userId", userId)
	c.Next()
}
