package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func main() {
	// Add a default user for demo/testing
	defaultID := uuid.New()
	log.Printf("Default user_id: %s\n", defaultID)
	store.AddUser(&User{ID: defaultID})
	setupRoutes()
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
