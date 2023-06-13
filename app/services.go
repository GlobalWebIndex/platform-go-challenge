package app

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	assetsrv1 "x-gwi/service/asset/v1"
	assetsrv2 "x-gwi/service/asset/v2"
	favouritesrv1 "x-gwi/service/favourite/v1"
	favouritesrv2 "x-gwi/service/favourite/v2"
	opinionsrv1 "x-gwi/service/opinion/v1"
	opinionsrv2 "x-gwi/service/opinion/v2"
	usersrv1 "x-gwi/service/user/v1"
	usersrv2 "x-gwi/service/user/v2"
)

// register instances of implementations of Services
// after initialisation of grpc.Server and before starting it
func (app *App) registerServices(registrar *grpc.Server) error {
	// app.server.ServiceRegistrar()
	if registrar == nil {
		return fmt.Errorf("registrar is nil") //nolint:goerr113
	}

	assetsrv1.RegisterGRPC(
		registrar,
		app.storage,
	)

	assetsrv2.RegisterGRPC(
		registrar,
		app.storage,
	)

	usersrv1.RegisterGRPC(
		registrar,
		app.storage,
	)

	usersrv2.RegisterGRPC(
		registrar,
		app.storage,
	)

	favouritesrv1.RegisterGRPC(
		registrar,
		app.storage,
	)

	favouritesrv2.RegisterGRPC(
		registrar,
		app.storage,
	)

	opinionsrv1.RegisterGRPC(
		registrar,
		app.storage,
	)

	opinionsrv2.RegisterGRPC(
		registrar,
		app.storage,
	)

	// Register reflection service on gRPC server.
	// google.golang.org/grpc/reflection
	reflection.Register(registrar)

	return nil
}
