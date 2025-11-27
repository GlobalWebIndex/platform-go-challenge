package http

import (
	"fmt"
	"github.com/akazantzidis/gwi-ass/internal/app"
	"github.com/akazantzidis/gwi-ass/internal/infra/http/auth"
	"github.com/akazantzidis/gwi-ass/internal/infra/http/favourite"
	"github.com/akazantzidis/gwi-ass/internal/pkg/middleware"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// Server Represents the http server running for this service
type Server struct {
	appServicesF app.Services
	router       *mux.Router
}

// NewServer HTTP Server constructor
func NewServer(appServicesF app.Services) *Server {
	httpServer := &Server{appServicesF: appServicesF}
	httpServer.router = mux.NewRouter()

	// use services to initialize handlers
	authHandler := auth.NewAuthHandler(appServicesF.AuthServices)

	// Public routes
	public := httpServer.router.PathPrefix("/").Subrouter()
	public.HandleFunc("/login", authHandler.Login).Methods("POST")
	public.HandleFunc("/refresh", authHandler.Refresh).Methods("POST")
	public.HandleFunc("/logout", authHandler.Logout).Methods("POST")
	public.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}).Methods("GET")

	// Private routes - apply JWT middleware
	private := httpServer.router.PathPrefix("/").Subrouter()
	private.Use(middleware.JWTMiddleware)

	h := favourite.NewHandler(httpServer.appServicesF.FavoriteServices, httpServer.appServicesF.UserServices)
	base := "/users/{userID}/favorites"

	private.HandleFunc(base, h.GetAll).Methods("GET")
	private.HandleFunc(base+"/{favoriteId}", h.GetByID).Methods("GET")
	private.HandleFunc(base, h.Create).Methods("POST")
	private.HandleFunc(base+"/{favoriteId}", h.Patch).Methods("PATCH")
	private.HandleFunc(base+"/{favoriteId}", h.Update).Methods("PUT")
	private.HandleFunc(base+"/{favoriteId}", h.Delete).Methods("DELETE")

	// admin-only route
	adminOnly := private.PathPrefix("/admin").Subrouter()
	adminOnly.Use(middleware.RequireRole("admin"))
	adminOnly.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("secret admin stats"))
	}).Methods("GET")

	http.Handle("/", httpServer.router)
	return httpServer
}

func (httpServer *Server) ListenAndServe(port string) {
	fmt.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
