package app

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"x-gwi/service"
	userAPIv1 "x-gwi/service/api/_user/v1"
	userAPIv2 "x-gwi/service/api/_user/v2"
	assetAPIv1 "x-gwi/service/api/asset/v1"
	assetAPIv2 "x-gwi/service/api/asset/v2"
	favouriteAPIv1 "x-gwi/service/api/favourite/v1"
	favouriteAPIv2 "x-gwi/service/api/favourite/v2"
	opinionAPIv1 "x-gwi/service/api/opinion/v1"
	opinionAPIv2 "x-gwi/service/api/opinion/v2"
)

// register instances of implementations of Services
// after initialisation of grpc.Server and before starting it
func (app *App) registerServices(registrar *grpc.Server) error {
	if registrar == nil {
		return fmt.Errorf("registrar is nil") //nolint:goerr113
	}

	var err error

	if _, err = assetAPIv1.RegisterGRPC(registrar,
		app.storage.ServiceStore(service.NameAsset),
	); err != nil {
		return fmt.Errorf("assetAPIv1.RegisterGRPC: %w", err)
	}

	if _, err = assetAPIv2.RegisterGRPC(registrar,
		app.storage.ServiceStore(service.NameAsset),
	); err != nil {
		return fmt.Errorf("assetAPIv2.RegisterGRPC: %w", err)
	}

	if _, err = userAPIv1.RegisterGRPC(registrar,
		app.storage.ServiceStore(service.NameUser),
	); err != nil {
		return fmt.Errorf("userAPIv1.RegisterGRPC: %w", err)
	}

	if _, err = userAPIv2.RegisterGRPC(registrar,
		app.storage.ServiceStore(service.NameUser),
	); err != nil {
		return fmt.Errorf("userAPIv2.RegisterGRPC: %w", err)
	}

	if _, err = favouriteAPIv1.RegisterGRPC(registrar,
		app.storage.ServiceStore(service.NameFavourite),
	); err != nil {
		return fmt.Errorf("favouriteAPIv1.RegisterGRPC: %w", err)
	}

	if _, err = favouriteAPIv2.RegisterGRPC(registrar,
		app.storage.ServiceStore(service.NameFavourite),
	); err != nil {
		return fmt.Errorf("favouriteAPIv2.RegisterGRPC: %w", err)
	}

	if _, err = opinionAPIv1.RegisterGRPC(registrar,
		app.storage.ServiceStore(service.NameOpinion),
	); err != nil {
		return fmt.Errorf("opinionAPIv1.RegisterGRPC: %w", err)
	}

	if _, err = opinionAPIv2.RegisterGRPC(registrar,
		app.storage.ServiceStore(service.NameOpinion),
	); err != nil {
		return fmt.Errorf("opinionAPIv2.RegisterGRPC: %w", err)
	}

	// Register reflection service on gRPC server.
	// google.golang.org/grpc/reflection
	reflection.Register(registrar)

	return nil
}
