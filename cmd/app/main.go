// Command server runs a tiny HTTP API that implements the GlobalWebIndex
// favourites challenge. It wires the in-memory store to the HTTP handlers and
// listens on :3002.
package main

import (
	"log"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/mr-polymath/platform-go-challenge/internal/api"
	"github.com/mr-polymath/platform-go-challenge/internal/db"
)

func main() {
	mStore := db.NewMemoryStore()
	db.SeedMemoryStore(mStore)
	handler := api.New(mStore)

	r := gin.Default()                     // includes Logger and Recovery by default
	r.Use(gzip.Gzip(gzip.BestCompression)) // Enable gzip compression

	handler.RegisterRoutes(r)

	port := ":3002"
	log.Printf("starting server on %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
