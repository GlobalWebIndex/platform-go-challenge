package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/gwi/favourites-api/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/gwi/favourites-api/internal/cache"
	"github.com/gwi/favourites-api/internal/service"
	"github.com/gwi/favourites-api/internal/store"
	"go.uber.org/zap"
)

type application struct {
	config            config
	storage           *store.Storage
	cache             cache.Cache
	favouritesService *service.FavouritesService
	logger            *zap.SugaredLogger
}

type config struct {
	addr     string
	env      string
	apiURL   string
	cacheTTL time.Duration
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set timeout on request context
	r.Use(middleware.Timeout(60 * time.Second))

	// Swagger docs
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Health check
		r.Get("/health", app.healthCheckHandler)

		// User routes
		r.Route("/users", func(r chi.Router) {

			r.Route("/{userID}", func(r chi.Router) {

				// Favourites routes
				r.Route("/favourites", func(r chi.Router) {
					r.Get("/", app.getFavouritesHandler)
					r.Post("/", app.addFavouriteHandler)

					r.Route("/{assetID}", func(r chi.Router) {
						r.Get("/", app.getFavouriteHandler)
						r.Delete("/", app.removeFavouriteHandler)
						r.Put("/description", app.updateDescriptionHandler)
					})
				})
			})
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	// Channel to listen for errors from the server
	shutdown := make(chan error)

	// Goroutine to listen for OS signals for graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Infow("server has started", "addr", app.config.addr, "env", app.config.env, "version", version)

	// Start the server
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Wait for shutdown signal
	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Infow("server has stopped", "addr", app.config.addr)

	return nil
}
