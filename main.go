package main

import (
	"fmt"
	"gwi/middleware"
	"gwi/routes"
	"log"
	"net/http"
)

func main() {
	router := routes.RegisterRoutes()
	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", middleware.RecoverMiddleware(router)))
}
