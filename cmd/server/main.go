package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/SpBalis/platform-go-challenge/internal/config"
	apphttp "github.com/SpBalis/platform-go-challenge/internal/http"
	"github.com/SpBalis/platform-go-challenge/internal/repo"
	"github.com/SpBalis/platform-go-challenge/internal/service"
)

func main() {
	cfg := config.FromEnv()
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL must be set")
	}

	db, err := repo.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer db.Close()

	// repos
	favRepo := repo.NewFavouritesRepo(db)
	assetRepo := repo.NewAssetsRepo(db)
	usersRepo := repo.NewUsersRepo(db)

	// services
	favSvc := &service.FavouritesService{Favs: favRepo, Assets: assetRepo}
	usersSvc := &service.UsersService{Repo: usersRepo}
	handler := apphttp.NewRouter(favSvc, usersSvc)

	// server
	srv := &http.Server{
		Addr:         cfg.Port,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("listening on %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
