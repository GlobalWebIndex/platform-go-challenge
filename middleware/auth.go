package middleware

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

const testToken = "gwi-token-12345" //for future implementation this token should be encrypted and handled properly

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
		}

		log.Printf("Received token: %s", token) // Debug statement

		// Validate the token
		if token != testToken {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}

		// Token is valid, proceed with the next handler
		return next(c)
	}
}
