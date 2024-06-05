package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/meletios/gwi-engineering-challenge/handlers"
	"github.com/meletios/gwi-engineering-challenge/middlewares"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/login", handlers.Login).Methods("POST")
	router.HandleFunc("/favorites/{userID}", handlers.GetFavorites).Methods("GET")
	router.HandleFunc("/favorites/{userID}", handlers.AddFavorite).Methods("POST")
	router.HandleFunc("/favorites/{userID}/{assetID}", handlers.RemoveFavorite).Methods("DELETE")
	router.HandleFunc("/favorites/{userID}/{assetID}", handlers.EditFavorite).Methods("PUT")

	router.Use(middlewares.JwtVerify)
	router.Use(middlewares.RateLimit)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", router)
}
