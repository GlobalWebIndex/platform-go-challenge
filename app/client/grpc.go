package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"x-gwi/app/logs"
	"x-gwi/app/pki"
)

const (
	timeoutDial = 10 * time.Second
)

func InternalClientConnGRPC(ctx context.Context, config *ConfigClient, pki *pki.PKI) (*grpc.ClientConn, context.CancelFunc, error) { //nolint:lll
	grpcDialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		// grpc.WithTransportCredentials(insecure.NewCredentials()), // grpc.WithInsecure(),
		grpc.WithTransportCredentials(credentials.NewTLS(pki.TLSConfigDial())),
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

func TryInternalConnGRPC(ctx context.Context, config *ConfigClient, pki *pki.PKI) bool {
	_, cancelDial, err := InternalClientConnGRPC(ctx, config, pki)
	if err != nil {
		return false
	}

	defer cancelDial()

	logs.Info().
		Str("client-grpc", config.GRPC.Address).
		Msg("OK")

	return true
}

/*
type GRPC struct {
	config         *ConfigClient
	grpcClientConn *grpc.ClientConn
	cancel         context.CancelFunc
	grpcDialOpts   []grpc.DialOption
	mu             sync.RWMutex
}

func NewGRPC(config *ConfigClient, pki *pki.PKI) *GRPC {
	c := &GRPC{ //nolint:exhaustruct
		config: config,
	}

	c.grpcDialOpts = []grpc.DialOption{
		grpc.WithBlock(),
		// grpc.WithTransportCredentials(insecure.NewCredentials()), // grpc.WithInsecure(),
		grpc.WithTransportCredentials(credentials.NewTLS(pki.TLSConfigDial())),
	}

	// c.DialContext(ctx)
	// c.ConnRPC()
	// pb call

	return c
}

func (c *GRPC) DialContext(ctx context.Context) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	// opts := []grpc.DialOption{
	// 	grpc.WithBlock(),
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()), // grpc.WithInsecure(),
	// }

	var (
		err     error
		ctxDial context.Context
	)

	ctxDial, c.cancel = context.WithTimeout(ctx, timeoutDial)
	c.grpcClientConn, err = grpc.DialContext(ctxDial, c.config.GRPC.Address, c.grpcDialOpts...)

	if err != nil {
		logs.Error().Err(err).Str("target", c.config.GRPC.Address).Msg("unsuccessful CliGRPC diall")

		return false
	}

	logs.Info().
		Str("client-grpc", c.config.GRPC.Address).
		Msg("grpc.DialContext OK")

	return true
}

func (c *GRPC) ValidConfig() bool {
	return c.config.Valid()
}

func (c *GRPC) ValidConnection() bool {
	return c.config.Valid() &&
		c.grpcClientConn != nil &&
		c.cancel != nil
}

func (c *GRPC) Close() {
	if c.cancel != nil {
		c.cancel()
	}

	if c.grpcClientConn != nil {
		c.grpcClientConn.Close()
	}
}
*/

//nolint:dupword
/*
func (c *GRPC) ConnRPC(ctx context.Context) (context.Context, context.CancelFunc, bool) {
	if !c.ValidConnection() {
		return nil, nil, false
	}

	var (
		rpcContext context.Context
		rpcCancel  context.CancelFunc
	)

	deadline, ok := ctx.Deadline()
	if ok {
		rpcContext, rpcCancel = context.WithDeadline(ctx, deadline)
	} else {
		rpcContext, rpcCancel = context.WithTimeout(ctx, 10*time.Minute) //nolint:gomnd
	}
	// opts = []grpc.CallOption{
	// _, _ = useraccountpb.NewUserAccountServiceClient(c.grpcClientConn).
	// Create(rpcContext, &useraccountpb.CreateRequest{})

	return rpcContext, rpcCancel, true
}
*/

func InsecureClientConnGRPC(ctx context.Context, config *ConfigClient) (*grpc.ClientConn, context.CancelFunc, error) {
	//nolint:exhaustruct,gosec
	tlsConfig := &tls.Config{
		// Certificates:       []tls.Certificate{certTLSdial},
		MinVersion: tls.VersionTLS13,
		// RootCAs:    p.CertPool(),
		InsecureSkipVerify: true,
	}

	grpcDialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		// grpc.WithTransportCredentials(insecure.NewCredentials()), // grpc.WithInsecure(),
		// grpc.WithTransportCredentials(credentials.NewTLS(pki.TLSConfigDial())),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
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
