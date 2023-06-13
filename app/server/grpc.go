package server

import (
	"errors"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"x-gwi/app/logs"
	"x-gwi/app/pki"
)

type GRPC struct {
	config         *ConfigServer
	grpcServer     *grpc.Server
	netListener    net.Listener
	grpcServerOpts []grpc.ServerOption
}

func (s *GRPC) initGRPC(config *ConfigServer, pki *pki.PKI) error { //nolint:unparam
	s.config = config

	s.grpcServerOpts = []grpc.ServerOption{
		grpc.Creds(credentials.NewTLS(pki.TLSConfigServerGRPC())),
		// grpc.UnaryInterceptor(authUnaryInterceptor),
	}

	s.grpcServer = grpc.NewServer(s.grpcServerOpts...)

	// create a gRPC server
	// insecure,  no []grpc.ServerOption for TLS and Auth
	// s.GRPC.grpcServer = grpc.NewServer()
	// pb Register...ServiceServer

	return nil
}

func (s *GRPC) serveGRPC() {
	var err error
	// create a listener on TCP port (def 9090)
	s.netListener, err = net.Listen("tcp", s.config.GRPC.Address)
	if err != nil {
		logs.Error().Err(err).Send()

		return
	}

	// create a gRPC server ^^^
	// start the server
	logs.Info().
		Str("addr", s.netListener.Addr().String()).
		Msg("starting gRCP server on TLS (no auth)")

	// Serve will return a non-nil error unless Stop or GracefulStop is called.
	if err = s.grpcServer.Serve(s.netListener); err != nil && errors.Is(err, grpc.ErrServerStopped) {
		logs.Error().Err(err).Send()

		return
	}
}

func (s *GRPC) closeGRPC() {
	if s.grpcServer == nil {
		return
	}

	s.grpcServer.GracefulStop()
}
