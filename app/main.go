package main

import (
	"context"
	"fmt"
	"gwi-platform/server"
	"gwi-platform/utils"
	"os"
	"os/signal"
	"time"
)

func main() {

	// set up logger
	logFile, err := utils.SetupLogger()
	if err != nil {
		utils.ErrorLogger.Println("Failed to set up logger:", err)
	}
	defer logFile.Close()

	// Create a new App instance
	app := server.NewApp()

	// Start the server in a goroutine
	go func() {
		fmt.Println("Server is starting on port 8080...")
		if err := app.Start(); err != nil {
			utils.ErrorLogger.Printf("Error starting server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	utils.InfoLogger.Println("Server is shutting down...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := app.Shutdown(ctx); err != nil {
		utils.ErrorLogger.Printf("Server forced to shutdown: %v", err)
	}

	utils.InfoLogger.Println("Server stopped")
}
