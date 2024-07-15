package main

import (
	"fmt"
	"net/http"

	"github.com/ceciivanov/platform-go-challenge/internal/handlers"
	"github.com/ceciivanov/platform-go-challenge/internal/repository"
	"github.com/ceciivanov/platform-go-challenge/internal/service"
	"github.com/gorilla/mux"
)

func main() {

	// Try to change the number of users and assets to see how the application behaves with large data sets
	NumberOfUsers := 2
	NumberOfAssets := 3

	// Create and initialize a new UsersRepository instance
	repo := repository.NewInMemoryUserRepository()
	repo.GenerateSampleUsers(NumberOfUsers, NumberOfAssets)

	// Create UserService and Handler for it
	userService := service.NewUserService(repo)
	userHandler := handlers.NewUserHandler(userService)

	// Create a new router from the Gorilla Mux package and register the respective routes for the userHandler
	r := mux.NewRouter()
	userHandler.RegisterRoutes(r)

	// Start the server
	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", r)
}
