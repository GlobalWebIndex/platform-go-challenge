package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"x-gwi/app/logs"
	"x-gwi/app/pki"
	pbsrvuser "x-gwi/proto/serv/_user/v1"
	pbsrvuser2 "x-gwi/proto/serv/_user/v2"
	pbsrvasset "x-gwi/proto/serv/asset/v1"
	pbsrvasset2 "x-gwi/proto/serv/asset/v2"
	pbsrvfavourite "x-gwi/proto/serv/favourite/v1"
	pbsrvfavourite2 "x-gwi/proto/serv/favourite/v2"
	pbsrvopinion "x-gwi/proto/serv/opinion/v1"
	pbsrvopinion2 "x-gwi/proto/serv/opinion/v2"
)

// asset user favourite opinion

type RESTGW struct {
	config       *ConfigServer
	gwMux        *runtime.ServeMux
	httpServer   *http.Server
	netListener  net.Listener
	grpcDialOpts []grpc.DialOption
}

func (s *RESTGW) initRESTGW(ctx context.Context, config *ConfigServer, pki *pki.PKI) error { //nolint:funlen
	s.config = config

	s.grpcDialOpts = []grpc.DialOption{
		// grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTransportCredentials(credentials.NewTLS(pki.TLSConfigDial())),
	}

	s.gwMux = runtime.NewServeMux()
	// err = gwMux.HandlePath("POST", "/api/file", s.servRESTGW...HandleBinaryFileUploadPortMaps)

	//nolint:exhaustruct,gomnd
	s.httpServer = &http.Server{
		Addr:              s.config.RESTGW.Address,
		Handler:           s.gwMux,
		TLSConfig:         pki.TLSConfigServerRESTGW(),
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       10 * time.Minute, // synch max ctx timeout
		ReadTimeout:       10 * time.Second,
	}

	var err error

	// register grpc rest gateway stubs (proto/serv/*/v*/*.pb.gw.go)

	if err = pbsrvasset.RegisterAssetServiceHandlerFromEndpoint(ctx, s.gwMux, s.config.GRPC.Address, s.grpcDialOpts); err != nil { //nolint:lll
		logs.Error().Err(err).Send()

		return fmt.Errorf("auth_pbsrv.RegisterAuthServiceHandlerFromEndpoint: %w", err)
	}

	if err = pbsrvasset2.RegisterAssetServiceHandlerFromEndpoint(ctx, s.gwMux, s.config.GRPC.Address, s.grpcDialOpts); err != nil { //nolint:lll
		logs.Error().Err(err).Send()

		return fmt.Errorf("auth_pbsrv.RegisterAuthServiceHandlerFromEndpoint: %w", err)
	}

	if err = pbsrvuser.RegisterUserServiceHandlerFromEndpoint(ctx, s.gwMux, s.config.GRPC.Address, s.grpcDialOpts); err != nil { //nolint:lll
		logs.Error().Err(err).Send()

		return fmt.Errorf("auth_pbsrv.RegisterAuthServiceHandlerFromEndpoint: %w", err)
	}

	if err = pbsrvuser2.RegisterUserServiceHandlerFromEndpoint(ctx, s.gwMux, s.config.GRPC.Address, s.grpcDialOpts); err != nil { //nolint:lll
		logs.Error().Err(err).Send()

		return fmt.Errorf("auth_pbsrv.RegisterAuthServiceHandlerFromEndpoint: %w", err)
	}

	if err = pbsrvfavourite.RegisterFavouriteServiceHandlerFromEndpoint(ctx, s.gwMux, s.config.GRPC.Address, s.grpcDialOpts); err != nil { //nolint:lll
		logs.Error().Err(err).Send()

		return fmt.Errorf("auth_pbsrv.RegisterAuthServiceHandlerFromEndpoint: %w", err)
	}

	if err = pbsrvfavourite2.RegisterFavouriteServiceHandlerFromEndpoint(ctx, s.gwMux, s.config.GRPC.Address, s.grpcDialOpts); err != nil { //nolint:lll
		logs.Error().Err(err).Send()

		return fmt.Errorf("auth_pbsrv.RegisterAuthServiceHandlerFromEndpoint: %w", err)
	}

	if err = pbsrvopinion.RegisterOpinionServiceHandlerFromEndpoint(ctx, s.gwMux, s.config.GRPC.Address, s.grpcDialOpts); err != nil { //nolint:lll
		logs.Error().Err(err).Send()

		return fmt.Errorf("auth_pbsrv.RegisterAuthServiceHandlerFromEndpoint: %w", err)
	}

	if err = pbsrvopinion2.RegisterOpinionServiceHandlerFromEndpoint(ctx, s.gwMux, s.config.GRPC.Address, s.grpcDialOpts); err != nil { //nolint:lll
		logs.Error().Err(err).Send()

		return fmt.Errorf("auth_pbsrv.RegisterAuthServiceHandlerFromEndpoint: %w", err)
	}

	return nil
}

func (s *RESTGW) serveRESTGW(_ context.Context) {
	// Note: Make sure the gRPC server is running properly and accessible
	// if !s.cliGRPC.connOK() {
	// 	return
	// }
	var err error
	// create a listener on TCP port (def 9080)
	s.netListener, err = net.Listen("tcp", s.config.RESTGW.Address)
	if err != nil {
		logs.Error().Err(err).Str("RESTGW.Address", s.config.RESTGW.Address).Send()

		return
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	logs.Info().
		Str("restgw", s.config.RESTGW.Address).Str("grpc", s.config.GRPC.Address).
		Str("addr", s.netListener.Addr().String()).
		Msg("starting HTTP REST GATEWAY on TLS (no auth)")

	// After Shutdown or Close, the returned error is ErrServerClosed.
	if err = s.httpServer.ServeTLS(s.netListener, "", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logs.Error().Err(err).Send()

		return
	}
}

func (s *RESTGW) closeRESTGW(ctx context.Context) {
	if s.httpServer == nil {
		return
	}

	s.httpServer.SetKeepAlivesEnabled(false)
	_ = s.httpServer.Shutdown(ctx)
}
