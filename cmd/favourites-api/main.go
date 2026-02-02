package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	assetrepo "github.com/manos/favourites/assets/adapters/outbound/repo_inmemory"
	assetapplication "github.com/manos/favourites/assets/application"
	"github.com/manos/favourites/favourites/adapters/outbound/assets_client"
	"github.com/manos/favourites/http/auth"

	favrepo "github.com/manos/favourites/favourites/adapters/outbound/repo_inmemory"
	favouritesapplication "github.com/manos/favourites/favourites/application"
	"github.com/manos/favourites/config"
	httpapi "github.com/manos/favourites/http"
)

func main() {
	env := flag.String("env", "", "Config environment (local/dev/test/prod). Overrides APP_ENV.")
	flag.Parse()

	cfg, err := config.LoadWithEnv(*env)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// assets hexagon
	assetsRepo := assetrepo.NewInMemoryAssetRepository()
	assetsService := assetapplication.NewAssetService(assetsRepo)

	// favourites hexagon

	assetClient := assets_client.NewAssetClient(assetsService)
	favRepo := favrepo.NewInMemoryFavouriteRepository()
	favService := favouritesapplication.NewFavouriteService(favRepo, assetClient)

	// http handlers
	favHandler := httpapi.NewFavouritesHandler(
		favService,
		cfg.Pagination.DefaultPage,
		cfg.Pagination.DefaultLimit,
		cfg.Pagination.MaxLimit,
	)
	assetsHandler := httpapi.NewAssetsHandler(assetsService)

	pubKeyPath := cfg.Auth.PublicKeyPath

	pubKey, err := auth.LoadRSAPublicKey(pubKeyPath)
	if err != nil {
		log.Fatalf("failed to load public key: %v", err)
	}

	jwtCfg := auth.JWTConfig{
		PublicKey: pubKey,
		Issuer:    cfg.Auth.Issuer,
		Audience:  cfg.Auth.Audience,
	}

	router := httpapi.NewRouter(jwtCfg, favHandler, assetsHandler, cfg.EnableDocs)

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("listening on %s", cfg.HTTPAddr)

	// Start server in the background to allow graceful shutdown.
	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-stop:
		log.Printf("received signal %s, shutting down", sig)
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown: %v", err)
	}
	log.Println("server stopped")
}
