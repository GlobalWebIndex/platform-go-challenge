package main

import (
	"log"
	"os"

	app "github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/cmd/api/server"
)

// @title Preferred Assets API
// @version 1.0
// @description API for managing assets, users, and favourites

// @contact.name API Support
// @contact.url http://localhost:8081
// @contact.email support@yourapp.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8081
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter "Bearer" followed by a space and your JWT token. Example: "Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiw..."
func main() {
	// Check if we're running in Docker
	if os.Getenv("DOCKER_ENV") == "true" {
		log.Println("Running in Docker container")
	}

	a := app.New()
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
