package server

import (
	"context"
	"gwi-platform/auth"
	"gwi-platform/database"
	"gwi-platform/handlers"
	"gwi-platform/utils"
	"net/http"
	"time"
)

// App struct to hold the HTTP server and any other app-wide dependencies
type App struct {
	server *http.Server
}

// NewApp creates and returns a new App instance
func NewApp() *App {

	// Initialize database connection
	var err error
	retryDuration := time.Duration(2) * time.Second
	for {
		err = database.InitDB()
		if err == nil {
			utils.ErrorLogger.Println("Database connection established successfully")
			break
		}
		utils.WarnLogger.Printf("Failed to initialize database: %v. Retrying!", err)
		time.Sleep(retryDuration)
	}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Set up routes
	mux.HandleFunc("/", utils.LogHandler(handlers.HomeHandler))
	mux.HandleFunc("/docs", utils.LogHandler(handlers.DocsHandler))
	mux.HandleFunc("/ping", utils.LogHandler(handlers.PingHandler))
	mux.HandleFunc("/assets/add", utils.LogHandler(handlers.AddHandler))
	mux.HandleFunc("/assets/delete", utils.LogHandler(handlers.DeleteHandler))
	mux.HandleFunc("/assets/modify", utils.LogHandler(handlers.ModifyHandler))
	mux.HandleFunc("/assets/get", utils.LogHandler(handlers.GetHandler))
	mux.HandleFunc("/user/add", utils.LogHandler(handlers.AddUserHandler))
	mux.HandleFunc("/user/favorite/add", utils.LogHandler(handlers.AddFavoriteHandler))
	mux.HandleFunc("/user/favorites", utils.LogHandler(auth.UserStatusAuth(handlers.GetUserFavoritesDetailedHandler)))

	// Create a new http.Server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	utils.InfoLogger.Println("*********Server Init**********")
	return &App{server: server}
}

// Start begins listening for requests
func (a *App) Start() error {
	return a.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
