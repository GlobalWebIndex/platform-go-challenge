package main

import (
	"log"
	"net/http"
	"platform-go-challenge/api"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/users/", api.FavoritesAPIs)

	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
