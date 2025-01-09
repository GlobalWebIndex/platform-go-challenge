package main

import (
	"log"
	"net/http"
	"platform-go-challenge/data"
	_ "platform-go-challenge/docs" // swag init generates this file
	"platform-go-challenge/handlers"
)

// @title GWI Challenge API
// @version 1.0
// @description This is a sample server for the GWI Challenge API.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email youremail@domain.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /

func main() {
	// Create an instance of InMemoryUserDataService
	userDataService := data.NewInMemoryUserDataService()

	// Create a new UserFavouritesHandler using the data service
	userFavouritesService := handlers.NewUserFavouritesService(userDataService)

	// Configure routes
	router := ConfigureRoutes(userFavouritesService)

	// Start the server
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", router))
}
