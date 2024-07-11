package main

import (
	"gwi-test/repositories"
	"gwi-test/server"
	"gwi-test/services"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := repositories.NewDbRepo()
	service := services.NewAssetHandler(db)

	server := server.NewServer(service)

	server.Initialize()
}
