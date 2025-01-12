package app

import (
	"context"
	"log/slog"
	"platform-go-challenge/internal/domain"
	"time"
)

type App struct {
	cancelFunc context.CancelFunc
	api        domain.Api
}

func New(api domain.Api) domain.App {
	return &App{api: api}
}

func (app *App) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	app.cancelFunc = cancel

	if err := app.api.Start(ctx); err != nil {
		slog.Error("api error", slog.Any("error", err))
		return err
	}

	return nil
}

func (app *App) Stop() {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if app.cancelFunc != nil {
		app.cancelFunc()
	}

	done := make(chan struct{})
	go func() {
		if err := app.api.Stop(shutdownCtx); err != nil {
			slog.Error("error stopping api", slog.Any("error", err))
		}
		close(done)
	}()

	select {
	case <-done:
		slog.Info("app stopped successfully")
	case <-shutdownCtx.Done():
		slog.Warn("app stop timed out")
	}
}
