package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/arhsxro/platform-go-challenge/api"
	"github.com/arhsxro/platform-go-challenge/config"
	"github.com/arhsxro/platform-go-challenge/storage"
)

func main() {

	cfg := config.LoadConfig()

	// Initialize database
	maxAttempts := 5
	var err error
	var dbInstance *storage.PostgresStore
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		dbInstance, err = storage.NewPostgresStore(cfg)
		if err == nil {
			break
		}
		log.Println("Attempt "+strconv.Itoa(attempt)+": Failed to connect to database: ", err)
		if attempt < maxAttempts {
			time.Sleep(time.Duration(attempt) * time.Second)
		}
	}
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbInstance.Close()

	// Initialize API with the database instance
	apiInstance := api.InitApi(dbInstance)
	router := apiInstance.InitRoutes()

	http.ListenAndServe(":8080", router)
}
