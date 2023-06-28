package server

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"x-gwi/app/client"
	"x-gwi/app/instance"
	"x-gwi/app/pki"
)

type Server struct {
	config *ConfigServer
	inst   *instance.Instance
	pki    *pki.PKI
	gRPC   *GRPC
	rESTGW *RESTGW
}

func NewServer(ctx context.Context, config *ConfigServer, inst *instance.Instance, pki *pki.PKI) (*Server, error) {
	var err error

	s := &Server{
		config: config,
		inst:   inst,
		pki:    pki,
		gRPC:   new(GRPC),
		rESTGW: new(RESTGW),
	}

	if err = s.gRPC.initGRPC(s.config, pki); err != nil {
		return nil, fmt.Errorf("s.newGRPC: %w", err)
	}

	if err = s.rESTGW.initRESTGW(ctx, s.config, pki); err != nil {
		return nil, fmt.Errorf("s.newRESTGW: %w", err)
	}

	return s, nil
}

// ServiceRegistrar to register service instance on gRPC
func (s *Server) ServiceRegistrar() *grpc.Server {
	return s.gRPC.grpcServer
}

func (s *Server) Serve(ctx context.Context) {
	go s.serve(ctx)
}

func (s *Server) serve(ctx context.Context) {
	go s.gRPC.serveGRPC()

	_ = client.TryInternalConnGRPC(ctx, client.NewConfigClient(), s.pki)

	go s.rESTGW.serveRESTGW(ctx)

	<-ctx.Done()
	s.Close(ctx)
}

func (s *Server) Close(ctx context.Context) {
	s.rESTGW.closeRESTGW(ctx)
	s.gRPC.closeGRPC()
}

/*
func (s *Server) Valid() bool {
	return s.config.Valid() &&
		s.inst.Valid() &&
		s.gRPC != nil &&
		s.gRPC.netListener != nil &&
		s.rESTGW != nil &&
		s.rESTGW.httpServer != nil
}
*/
