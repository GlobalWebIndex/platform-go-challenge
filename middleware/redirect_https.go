package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HTTPSRedirect redirects HTTP requests to HTTPS.
func HTTPSRedirect(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Scheme() == "http" {
			target := "https://" + c.Request().Host + c.Request().RequestURI
			return c.Redirect(http.StatusMovedPermanently, target)
		}
		return next(c)
	}
}
