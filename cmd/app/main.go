// Command server runs a tiny HTTP API that implements the GlobalWebIndex
// favourites challenge. It wires the in-memory store to the HTTP handlers and
// listens on :3002.
package main

import (
	"log"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	// swagger
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/mr-polymath/platform-go-challenge/internal/api"
	"github.com/mr-polymath/platform-go-challenge/internal/db"
	_ "github.com/mr-polymath/platform-go-challenge/internal/swagger"
)

// @title Favourites API
// @version 1.0
// @description A small API for managing user favourites.
// @contact.name API Support
// @contact.email support@example.com
// @host localhost:3002
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter 'Bearer token'
func main() {
	mStore := db.NewMemoryStore()
	db.SeedMemoryStore(mStore)
	handler := api.New(mStore)

	r := gin.Default()                     // includes Logger and Recovery by default
	r.Use(gzip.Gzip(gzip.BestCompression)) // Enable gzip compression

	// Swagger UI route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	handler.RegisterRoutes(r)

	port := ":3002"
	log.Printf("starting server on %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
