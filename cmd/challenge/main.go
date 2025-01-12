package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	_ "platform-go-challenge/internal/api/docs"

	"platform-go-challenge/internal/api"

	"platform-go-challenge/internal/app"
	"platform-go-challenge/internal/domain"
	"platform-go-challenge/internal/repository"
	"platform-go-challenge/internal/service"
)

func main() {
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)

	api := api.New(
		withService(
			withSqlRepo()))

	app := app.New(api)
	go func() {
		errCh <- app.Start(ctx)
	}()
	defer app.Stop()

	select {
	case err := <-errCh:
		slog.Error("service shutdown with error", slog.Any("error", err))
	case sig := <-interruptCh:
		slog.Info("received signal", "signal", sig)
	}
}

func withService(r domain.Repository) domain.AssetService {
	return service.AssetService(r)
}

func withSqlRepo() domain.Repository {
	repo, err := repository.New()
	if err != nil {
		panic(err)
	}
	return repo
}
