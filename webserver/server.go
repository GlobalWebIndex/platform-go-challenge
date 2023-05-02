package webserver

import (
	"challenge/cache"
	"challenge/cfg"
	"challenge/middleware"
	"challenge/storage"

	"github.com/gorilla/mux"
)

var (
	Version    = "v1.0.0"
	BuildTime  = ""
	BuildUser  = ""
	CommitHash = ""
)

type WebServer struct {
	Server        mux.Router
	Configuration *cfg.Configuration
}

// CreateServer create a new Server instance.
// Optionally, it accepts options for mocking purposes
func CreateServer(options ...string) (*WebServer, error) {
	var err error
	server := WebServer{}
	if len(options) > 0 {
		server.Configuration, err = cfg.ReadConfig(options[0])
	} else {
		server.Configuration, err = cfg.ReadConfig("")
	}

	if err != nil {
		return nil, err
	}
	middleware.JwtSecret = server.Configuration.JwtSecret
	server.Server = *mux.NewRouter()
	if len(options) > 1 {
		err = storage.ConnectToStorage(options[1], server.Configuration.AdminUser, server.Configuration.AdminPass)
	} else {
		err = storage.ConnectToStorage(server.Configuration.StorageOption, server.Configuration.AdminUser, server.Configuration.AdminPass)
	}
	if err != nil {
		return nil, err
	}
	if len(options) > 2 {
		err = cache.ConnectToCache(options[2])
	} else {
		err = cache.ConnectToCache(server.Configuration.CacheOption)
	}
	if err != nil {
		return nil, err
	}
	// Register the routes
	RegisterRoutes(&server)
	// Attach the profiler if selected
	if server.Configuration.Profiling {
		AttachProfiler(&server)
	}
	return &server, nil
}
