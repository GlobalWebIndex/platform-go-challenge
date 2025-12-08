// Command api is the application entrypoint for the Assets Favourites API service.
// It reads configuration, initializes the HTTP router, and starts the server with graceful shutdown support.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GlobalWebIndex/platform-go-challenge/internal/app"
	"github.com/GlobalWebIndex/platform-go-challenge/internal/config"
)

func main() {
	cfg := config.Load()
	router := app.NewRouter()
	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadTimeout:       10 * time.Second, // Timeout for reading the request body
		WriteTimeout:      10 * time.Second, // Timeout for writing the response body
		IdleTimeout:       60 * time.Second, // Timeout for idle connections
		ReadHeaderTimeout: 5 * time.Second,  // Timeout for reading the request headers
		MaxHeaderBytes:    1 << 20,          // 1MB max header size
	}

	// Start the HTTP server in a goroutine to allow graceful shutdown handling.
	go func() {
		log.Printf("server listening on %s", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	// Wait for interrupt signal (SIGINT or SIGTERM) to initiate graceful shutdown.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	// Gracefully shutdown the server with a 5-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
