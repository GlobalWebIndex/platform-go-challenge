package main

import (
	"log"
	"net/http"
	"os"

	"github.com/artemis13/platform-go-challenge/routes"
	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Print loaded environment variables (for debugging)
	log.Println("GOPROXY:", os.Getenv("GOPROXY"))

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//security-related headers
	e.Use(middleware.Secure())

	//CORS settings
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	//limit the number of requests
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	// Use custom middleware for HTTPS redirection -- a valid certificate should be present
	//e.Use(myMiddleware.HTTPSRedirect)

	//also caching midlleware should me applyed for better performance

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "GWI Go Platform Challenge running!\n")
	})

	routes.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":1323"))
}
