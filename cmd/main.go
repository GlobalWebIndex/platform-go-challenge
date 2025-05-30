package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"pgc/internal/config"
	"pgc/internal/server"
)

func main() {
	// Create a context that listens for SIGINT or SIGTERM signals for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop() // Ensure resources associated with the signal context are released

	cfg := config.Load()

	if err := server.Run(ctx, cfg); err != nil {
		slog.Error("failed to run the application", "error", err.Error())
		cleanExit(1)
	}
}

// cleanExit is used to exit the application while ensuring deferred functions are executed
func cleanExit(code int) {
	// Allow deferred functions to run before exiting
	defer os.Exit(code)
}
