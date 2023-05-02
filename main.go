package main

import (
	"challenge/middleware"
	"challenge/webserver"
	"log"
	"net/http"
)

func main() {
	// Create a new server
	webserver, err := webserver.CreateServer()
	if err != nil {
		log.Fatal(err.Error())
	}
	// Add the logging middleware
	loggedRouter := middleware.WithLogging(&webserver.Server)
	// Add the prometheus middleware if selected
	if webserver.Configuration.Metrics {
		webserver.Server.Use(middleware.PrometheusMiddleware)
	}
	// Start the web server
	if err := http.ListenAndServe(webserver.Configuration.Address+":"+webserver.Configuration.Port, loggedRouter); err != nil {
		log.Fatal("Fatal status: ", err.Error())
	}
}
