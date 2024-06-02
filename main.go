package main

import (
	"log"
	"net/http"
	"os"
	"user-favorites-api/router"
	"user-favorites-api/store"
)

func main() {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        log.Fatal("DATABASE_URL environment variable is required")
    }

    store, err := store.NewStore(dbURL)
    if err != nil {
        log.Fatalf("failed to connect to database: %v", err)
    }
    defer store.DB.Close()

    r := router.NewRouter(store)
    http.ListenAndServe(":8080", r)
}
