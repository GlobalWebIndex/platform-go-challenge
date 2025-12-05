package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/users/", usersHandler)
	log.Printf("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}