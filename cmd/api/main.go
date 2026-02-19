package main

import (
	"log"
	"time"

	"go.uber.org/zap"

	"github.com/gwi/favourites-api/internal/cache"
	"github.com/gwi/favourites-api/internal/env"
	"github.com/gwi/favourites-api/internal/service"
	"github.com/gwi/favourites-api/internal/store"
)

const version = "1.0.0"

// @title           Favourites API
// @version         1.0.0
// @description     API for managing user favourite assets
//
//	@termsOfService	http://swagger.io/terms/
//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//
// @BasePath        /api/v1
func main() {
	// config
	cfg := config{
		addr:     env.GetString("ADDR", ":8080"),
		env:      env.GetString("ENV", "development"),
		apiURL:   env.GetString("API_URL", "localhost:8080"),
		cacheTTL: env.GetDuration("CACHE_TTL", 5*time.Minute),
	}

	// logger
	var logger *zap.SugaredLogger
	if cfg.env == "production" {
		logger = zap.Must(zap.NewProduction()).Sugar()
	} else {
		logger = zap.Must(zap.NewDevelopment()).Sugar()
	}
	defer logger.Sync()

	storage := store.NewStorage()
	logger.Info("Initialized storage layer")

	cache := cache.NewInMemoryCache()
	logger.Info("Initialized in-memory cache")

	favouritesService := service.NewFavouritesService(storage, cache, cfg.cacheTTL)
	logger.Info("Initialized favourites service with caching")

	// app
	app := &application{
		config:            cfg,
		storage:           storage,
		cache:             cache,
		favouritesService: favouritesService,
		logger:            logger,
	}

	// mount routes and start server
	mux := app.mount()

	if err := app.run(mux); err != nil {
		log.Fatal(err)
	}
}
