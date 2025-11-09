package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KostasDasios/platform-go-challenge/internal/config"
	"github.com/KostasDasios/platform-go-challenge/internal/server"
)

func main() {
	// Load configuration from .env
	cfg := config.LoadConfig()

	// Build server using internal layers
	s := server.NewServer(cfg)

	// Configure HTTP server with proper timeouts
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      s.Handler(),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Start server in background goroutine
	go func() {
		log.Printf("[INFO] Server listening on :%s (env: %s)\n", cfg.Port, cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[ERROR] Server failed: %v", err)
		}
	}()

	// Listen for OS signals (Ctrl+C / docker stop)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	sig := <-quit
	log.Printf("[INFO] Caught signal: %v, initiating graceful shutdown...", sig)

	// Gracefully shut down server with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[ERROR] Graceful shutdown failed: %v", err)
	} else {
		log.Println("[INFO] Server shut down cleanly.")
	}

	log.Println("[INFO] Server exiting")
}
