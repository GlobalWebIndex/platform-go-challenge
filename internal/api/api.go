package api

import (
	"context"
	"net/http"
	"time"

	"platform-go-challenge/internal/api/handlers"
	"platform-go-challenge/internal/api/middleware"
	"platform-go-challenge/internal/domain"

	"github.com/gorilla/mux"
)

type Api struct {
	httpSrv    *http.Server
	cancelFunc context.CancelFunc
}

// @title           Asset API
// @version         1.0
// @description     API for managing assets
// @host      localhost:8080
// @BasePath  /
func New(service domain.AssetService) *Api {
	api := &Api{}
	router := mux.NewRouter()

	// Create handlers
	h := handlers.New(service)
	h.RegisterRoutes(router)

	// Api context
	ctx, cancel := context.WithCancel(context.Background())
	api.cancelFunc = cancel

	// Apply middleware chain
	handler := middleware.Chain(
		router,
		middleware.RecoveryMiddleware(),
		middleware.RateLimitMiddleware(20),
		middleware.ContextMiddleware(ctx),
		middleware.LoggingMiddleware(),
		middleware.AuthMiddleware(),
	)

	httpSrv := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	api.httpSrv = httpSrv

	return api
}

func (api *Api) Start(ctx context.Context) error {
	return api.httpSrv.ListenAndServe()
}

func (api *Api) Stop(ctx context.Context) error {
	api.cancelFunc()
	return api.httpSrv.Shutdown(ctx)
}
