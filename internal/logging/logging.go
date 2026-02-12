package logging

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

const ErrorKey string = "error"
const loggerKey string = "logger"
const logAttributesNumber = 8 // Preallocate for common attributes. Go will reallocate if more is added to the slice.

// NewLogger creates a new JSON logger configured for production use.
// It sets the logger as the default slog logger and returns it.
func NewLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
	return logger
}

// RequestLogger is a middleware that creates a request-scoped logger with the request ID
// and stores it in the context for use by all downstream handlers and layers.
func RequestLogger(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := middleware.GetReqID(r.Context())
			log := logger.With(slog.String("request_id", requestID))
			ctx := context.WithValue(r.Context(), loggerKey, log)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// FromContext retrieves the request-scoped logger from context.
// If no logger is found, returns the default slog logger.
func FromContext(ctx context.Context) *slog.Logger {
	if log, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return log
	}
	return slog.Default()
}

// NewContextWithLogger creates a new context with the given logger attached.
// Useful for tests or background jobs where there's no HTTP request.
func NewContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// LogBuilder provides a fluent API for building structured log entries.
// It reduces verbosity by eliminating repeated slog.String() calls inside log.Error, log.Warn, and log.Info functions.
type LogBuilder struct {
	logger *slog.Logger
	attrs  []any
}

// Log creates a LogBuilder from context, extracting the request-scoped logger.
func Log(ctx context.Context) *LogBuilder {
	return &LogBuilder{
		logger: FromContext(ctx),
		attrs:  make([]any, 0, logAttributesNumber),
	}
}

// With creates a LogBuilder from an existing slog.Logger.
func With(logger *slog.Logger) *LogBuilder {
	return &LogBuilder{
		logger: logger,
		attrs:  make([]any, 0, logAttributesNumber),
	}
}

// Layer adds the "layer" field (e.g., "routes", "handler", "repository").
func (b *LogBuilder) Layer(layer string) *LogBuilder {
	b.attrs = append(b.attrs, slog.String("layer", layer))
	return b
}

// Op adds the "operation" field (e.g., "AddFavourite", "GetUser").
func (b *LogBuilder) Op(operation string) *LogBuilder {
	b.attrs = append(b.attrs, slog.String("operation", operation))
	return b
}

// User adds the "user_id" field.
func (b *LogBuilder) User(userID string) *LogBuilder {
	b.attrs = append(b.attrs, slog.String("user_id", userID))
	return b
}

// Asset adds the "asset_id" field.
func (b *LogBuilder) Asset(assetID string) *LogBuilder {
	b.attrs = append(b.attrs, slog.String("asset_id", assetID))
	return b
}

// AssetType adds the "asset_type" field.
func (b *LogBuilder) AssetType(assetType string) *LogBuilder {
	b.attrs = append(b.attrs, slog.String("asset_type", assetType))
	return b
}

// Str adds a custom string field.
func (b *LogBuilder) Str(key, value string) *LogBuilder {
	b.attrs = append(b.attrs, slog.String(key, value))
	return b
}

// Int adds a custom int field.
func (b *LogBuilder) Int(key string, value int) *LogBuilder {
	b.attrs = append(b.attrs, slog.Int(key, value))
	return b
}

// Time adds a custom time field.
func (b *LogBuilder) Time(key string, value time.Time) *LogBuilder {
	b.attrs = append(b.attrs, slog.Time(key, value))
	return b
}

// Bool adds a custom bool field.
func (b *LogBuilder) Bool(key string, value bool) *LogBuilder {
	b.attrs = append(b.attrs, slog.Bool(key, value))
	return b
}

// Any adds a custom field of any type.
func (b *LogBuilder) Any(key string, value any) *LogBuilder {
	b.attrs = append(b.attrs, slog.Any(key, value))
	return b
}

// Err adds the "error" field from an error.
func (b *LogBuilder) Err(err error) *LogBuilder {
	if err != nil {
		b.attrs = append(b.attrs, slog.String(ErrorKey, err.Error()))
	}
	return b
}

// Info logs at INFO level.
func (b *LogBuilder) Info(msg string) {
	b.logger.Info(msg, b.attrs...)
}

// Warn logs at WARN level.
func (b *LogBuilder) Warn(msg string) {
	b.logger.Warn(msg, b.attrs...)
}

// Error logs at ERROR level.
func (b *LogBuilder) Error(msg string) {
	b.logger.Error(msg, b.attrs...)
}

// Debug logs at DEBUG level.
func (b *LogBuilder) Debug(msg string) {
	b.logger.Debug(msg, b.attrs...)
}
