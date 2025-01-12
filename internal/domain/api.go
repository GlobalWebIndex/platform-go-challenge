package domain

import "context"

// Api manages the lifecycle of an API server, handling startup and graceful shutdown
type Api interface {
	// Start begins the API server and blocks until context cancellation or error.
	// Returns error on startup failure, nil on context cancellation.
	Start(context.Context) error

	// Stop gracefully shuts down the API server, waiting for in-flight requests
	// to complete or context timeout. Returns error if shutdown fails.
	Stop(context.Context) error
}
