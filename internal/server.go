package internal

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/giannis84/platform-go-challenge/internal/database"
	"github.com/giannis84/platform-go-challenge/internal/logging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// RoutesRegistry is a function that registers routes on a chi.Router
type RoutesRegistry func(r chi.Router)

// Service wraps an HTTP server with its configuration and router
type Service struct {
	// Configuration fields (set before Init)
	Addr         string
	Logger       *slog.Logger
	DB           *sql.DB
	Routes       RoutesRegistry
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	// Runtime fields (populated by Init)
	HTTPServer *http.Server
	Router     *chi.Mux
}

// Init initializes the service by setting up the router and HTTP server
func (s *Service) Init() {
	s.Router = chi.NewRouter()

	// Set up database connection for the database package
	if s.DB != nil {
		database.DB = s.DB
	}

	// Initialize common middleware
	s.Router.Use(middleware.RequestID)
	s.Router.Use(logging.RequestLogger(s.Logger))
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)

	// Register routes
	if s.Routes != nil {
		s.Routes(s.Router)
	}

	// Apply default timeouts if not provided
	readTimeout := s.ReadTimeout
	if readTimeout == 0 {
		readTimeout = 15 * time.Second
	}
	writeTimeout := s.WriteTimeout
	if writeTimeout == 0 {
		writeTimeout = 15 * time.Second
	}
	idleTimeout := s.IdleTimeout
	if idleTimeout == 0 {
		idleTimeout = 60 * time.Second
	}

	s.HTTPServer = &http.Server{
		Addr:         s.Addr,
		Handler:      s.Router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}
}

// ListenAndServeWrapper starts the http service
func (s *Service) ListenAndServeWrapper(service string) error {
	s.Logger.Info("starting http service", service, slog.String("port", s.HTTPServer.Addr))
	return s.HTTPServer.ListenAndServe()
}