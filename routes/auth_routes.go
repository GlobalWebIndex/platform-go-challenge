package routes

import (
    "net/http"
    "github.com/dintzeler/platform-go-challenge/controllers"
)

func SetupAuthRoutes() {
    http.HandleFunc("/auth/login", controllers.LoginHandler)
}