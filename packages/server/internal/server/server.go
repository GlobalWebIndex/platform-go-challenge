package server

import (
	"net/http"
	"time"

	"github.com/alexandros-georgantas/platform-go-challenge/internal/database"
	"gorm.io/gorm"
)

type Server struct {
	db *gorm.DB
}

func Create() *http.Server {
	NewServer := &Server{
		db: database.GetDBConnection(),
	}

	// Basic server config
	server := &http.Server{
		Addr:         ":3000",
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
