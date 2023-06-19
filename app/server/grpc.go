package server

import (
	"errors"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

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
	logger := logs.Logger
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			// logging.StartCall,
			// logging.FinishCall,
			// logging.PayloadReceived,
			logging.PayloadSent,
		),
	}
	loggingOptsStream := []logging.Option{
		logging.WithLogOnEvents(
		// logging.StartCall,
		// logging.FinishCall,
		// logging.PayloadReceived,
		// logging.PayloadSent,
		),
	}
	// Define customfunc to handle panic
	// api servive skip 12,
	// return status.Errorf(codes.Unknown, "panic triggered: %v @at: %s", p, logs.LogTraceX(12))
	// core method skip 11,
	// return status.Errorf(codes.Unknown, "panic triggered: %v @at: %s", p, logs.LogTraceX(11))
	// const skip = 11

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) { //nolint:nonamedreturns
			return status.Errorf(codes.Unknown, "panic triggered: %v", p)
		}),
	}
	// Create a server. Recovery handlers should typically be last in the chain so that other middleware
	// (e.g. logging) can operate on the recovered state instead of being directly affected by any panic

	s.config = config

	s.grpcServerOpts = []grpc.ServerOption{
		grpc.Creds(credentials.NewTLS(pki.TLSConfigServerGRPC())),
		// grpc.UnaryInterceptor(authUnaryInterceptor),
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(logs.InterceptorLogger(logger), loggingOpts...),
			recovery.UnaryServerInterceptor(recoveryOpts...),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(logs.InterceptorLogger(logger), loggingOptsStream...),
			recovery.StreamServerInterceptor(recoveryOpts...),
		),
		// grpc.UnaryInterceptor(
		// 	recovery.UnaryServerInterceptor(recoveryOpts...),
		// ),
		// grpc.StreamInterceptor(
		// 	recovery.StreamServerInterceptor(recoveryOpts...),
		// ),
	}

	s.grpcServer = grpc.NewServer(s.grpcServerOpts...)
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
