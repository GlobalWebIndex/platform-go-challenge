package handlers

import (
	"net/http"
	"platform-go-challenge/internal/domain/dto"
	"platform-go-challenge/internal/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// NewJWT
// @Summary Create a new JWT
// @Description Creates a new JSON Web Token for authentication.
// @Accept json
// @Produce json
// @Param payload body dto.JWTPayload true "User ID for token generation"
// @Success 200 {object} map[string]string "Successful token creation"
// @Failure 400 {object} middleware.APIError "Invalid JSON body"
// @Failure 500 {object} map[string]string "Failed to create token"
// @Router /jwt [post]
func NewJWT(secretKey string) gin.HandlerFunc {
	//TODO consider adding context in jwt handling e.g. claims, roles
	return func(c *gin.Context) {

		// Define a struct to bind the request body to.
		var reqBody dto.JWTPayload

		// Bind the JSON payload to the struct.
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.Error(middleware.NewAPIError(http.StatusBadRequest, "invalid json body", "User Id missing"))
			return
		}

		var role string
		if reqBody.GetUserID() == "admin_user_id" {
			role = "admin"
		} else {
			role = "user"
		}

		// Define the JWT claims.
		claims := jwt.MapClaims{
			"sub":  reqBody.GetUserID(),
			"role": role,
			"exp":  time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours.
			"iat":  time.Now().Unix(),
		}

		// Create the token.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Sign the token with the secret key.
		signedToken, err := token.SignedString([]byte(secretKey))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
			return
		}

		// Return the token in the response.
		c.JSON(http.StatusOK, gin.H{"token": signedToken})
	}
}
