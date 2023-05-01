package main

import (
	"net/http"
	"platform-go-challenge/integrations"
	"platform-go-challenge/server"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Create a DB connection
	integrations.CreateDBConnection()

	// Start http server
	logger := log.New()
	router := mux.NewRouter()
	httpServer := &http.Server{
		Handler: router,
	}
	server := server.NewServer(httpServer, router, logger)
	server.Start()
}
