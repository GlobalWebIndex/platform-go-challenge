package main

import (
	"fmt"
	"net/http"
    "github.com/dintzeler/platform-go-challenge/routes"
	"github.com/joho/godotenv"
	"log"
)

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}

func main() {
	routes.SetupFavoritesRoutes()
    routes.SetupAuthRoutes()
    http.HandleFunc("/docs/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./openapi/redoc-static.html")
    })
	fmt.Println("Server is running on port 8090")
	http.ListenAndServe(":8090", nil)

}