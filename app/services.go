package app

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"x-gwi/service"
	assetsrv1 "x-gwi/service/api/asset/v1"
	assetsrv2 "x-gwi/service/api/asset/v2"
	favouritesrv1 "x-gwi/service/api/favourite/v1"
	favouritesrv2 "x-gwi/service/api/favourite/v2"
	opinionsrv1 "x-gwi/service/api/opinion/v1"
	opinionsrv2 "x-gwi/service/api/opinion/v2"
	usersrv1 "x-gwi/service/api/user/v1"
	usersrv2 "x-gwi/service/api/user/v2"
)

// register instances of implementations of Services
// after initialisation of grpc.Server and before starting it
func (app *App) registerServices(registrar *grpc.Server) error {
	// app.server.ServiceRegistrar()
	if registrar == nil {
		return fmt.Errorf("registrar is nil") //nolint:goerr113
	}

	_, _ = assetsrv1.RegisterGRPC(
		registrar,
		app.storage.ServiceStore(service.NameAsset),
	)

	_, _ = assetsrv2.RegisterGRPC(
		registrar,
		app.storage.ServiceStore(service.NameAsset),
	)

	_, _ = usersrv1.RegisterGRPC(
		registrar,
		app.storage.ServiceStore(service.NameUser),
	)

	_, _ = usersrv2.RegisterGRPC(
		registrar,
		app.storage.ServiceStore(service.NameUser),
	)

	_, _ = favouritesrv1.RegisterGRPC(
		registrar,
		app.storage.ServiceStore(service.NameFavourite),
	)

	_, _ = favouritesrv2.RegisterGRPC(
		registrar,
		app.storage.ServiceStore(service.NameFavourite),
	)

	_, _ = opinionsrv1.RegisterGRPC(
		registrar,
		app.storage.ServiceStore(service.NameOpinion),
	)

	_, _ = opinionsrv2.RegisterGRPC(
		registrar,
		app.storage.ServiceStore(service.NameOpinion),
	)

	// Register reflection service on gRPC server.
	// google.golang.org/grpc/reflection
	reflection.Register(registrar)

	return nil
}
