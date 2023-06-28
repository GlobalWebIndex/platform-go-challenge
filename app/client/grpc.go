package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"x-gwi/app/logs"
)

const (
	timeoutDial = 10 * time.Second
)

func ConnGRPC(ctx context.Context, config *ConfigClient, tlsConfigDial *tls.Config) (*grpc.ClientConn, context.CancelFunc, error) { //nolint:lll
	if tlsConfigDial == nil {
		//nolint:exhaustruct,gosec
		tlsConfigDial = &tls.Config{
			// Certificates:       []tls.Certificate{certTLSdial},
			MinVersion: tls.VersionTLS13,
			// RootCAs:    p.CertPool(),
			InsecureSkipVerify: true,
		}
	}

	grpcDialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		// grpc.WithTransportCredentials(insecure.NewCredentials()), // grpc.WithInsecure(),
		// grpc.WithTransportCredentials(credentials.NewTLS(pki.TLSConfigDial())),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfigDial)),
	}

	t := time.Now()
	ctxDial, cancelDial := context.WithTimeout(ctx, timeoutDial)
	// defer cancelDial()

	// config.GRPC.Address = "localhost:9090"
	grpcClientConn, err := grpc.DialContext(ctxDial, config.GRPC.Address, grpcDialOpts...)
	if err != nil {
		cancelDial()

		return nil, nil, fmt.Errorf("grpc.DialContext: (t: %v) %w", time.Since(t), err)
	}

	return grpcClientConn, cancelDial, nil
}

func TryConnGRPC(ctx context.Context, config *ConfigClient, tlsConfigDial *tls.Config) bool {
	_, cancelDial, err := ConnGRPC(ctx, config, tlsConfigDial)
	if err != nil {
		return false
	}

	defer cancelDial()

	logs.Info().
		Str("client-grpc", config.GRPC.Address).
		Msg("OK")

	return true
}
