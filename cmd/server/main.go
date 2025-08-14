package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gwi-favorites-service/internal/config"
	"gwi-favorites-service/internal/domain"
	"gwi-favorites-service/internal/handler"
	"gwi-favorites-service/internal/repository/memory"
	"gwi-favorites-service/internal/service"
	"gwi-favorites-service/pkg/logger"

	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logger
	log := logger.NewLogger()

	// Load configuration
	cfg := config.Load()
	log.WithField("port", cfg.Port).Info("Starting GWI Favorites Service")

	// Initialize repository
	repo := memory.NewRepository()

	// Seed some sample data
	seedSampleData(repo, log)

	// Initialize service
	favoritesService := service.NewFavoritesService(repo, log)

	// Initialize HTTP handler
	httpHandler := handler.NewHandler(favoritesService, log)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      httpHandler.SetupRoutes(),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.WithField("addr", server.Addr).Info("HTTP server starting")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Failed to start HTTP server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("Server forced to shutdown")
	}

	log.Info("Server exited")
}

func seedSampleData(repo *memory.Repository, log *logrus.Logger) {
	log.Info("Seeding sample data...")

	// Create sample users
	users := []*domain.User{
		domain.NewUser("user1", "john@example.com", "John Doe"),
		domain.NewUser("user2", "jane@example.com", "Jane Smith"),
		domain.NewUser("user3", "bob@example.com", "Bob Johnson"),
	}

	for _, user := range users {
		if err := repo.CreateUser(user); err != nil {
			log.WithError(err).Error("Failed to create sample user")
		}
	}

	// Create sample assets
	chartData := []domain.ChartDataPoint{
		{X: "Jan", Y: 100},
		{X: "Feb", Y: 150},
		{X: "Mar", Y: 200},
	}

	assets := []domain.Asset{
		domain.NewChart("chart1", "Monthly Sales", "Month", "Sales ($)", "Sales performance chart", chartData),
		domain.NewInsight("insight1", "40% of millennials spend more than 3 hours on social media daily", "Social media usage insight", []string{"social", "millennials"}, "demographics"),
		domain.NewAudience("audience1", "Gaming enthusiasts aged 24-35"),
	}

	// Set audience properties
	if audience, ok := assets[2].(*domain.Audience); ok {
		audience.Gender = []string{"Male", "Female"}
		audience.AgeGroups = []string{"24-35"}
		audience.SocialMediaHours = "3+"
		audience.PurchasesLastMonth = 5
		audience.BirthCountries = []string{"US", "UK", "CA"}
	}

	for _, asset := range assets {
		if err := repo.CreateAsset(asset); err != nil {
			log.WithError(err).Error("Failed to create sample asset")
		}
	}

	log.Info("Sample data seeded successfully")
}
