package server

import (
	"context"
	"fmt"
	"runtime"

	"google.golang.org/grpc"

	"x-gwi/app/client"
	"x-gwi/app/instance"
	"x-gwi/app/pki"
)

type Server struct {
	config  *ConfigServer
	inst    *instance.Instance
	gRPC    *GRPC
	rESTGW  *RESTGW
	cliGRPC *client.GRPC
	cliHTTP *client.HTTP
}

func NewServer(ctx context.Context, config *ConfigServer, inst *instance.Instance, pki *pki.PKI) (*Server, error) {
	var err error

	s := &Server{
		config:  config,
		inst:    inst,
		gRPC:    new(GRPC),
		rESTGW:  new(RESTGW),
		cliGRPC: nil, // new(client.GRPC),
		cliHTTP: nil, // new(client.HTTP),
	}

	if err = s.gRPC.initGRPC(s.config, pki); err != nil {
		return nil, fmt.Errorf("s.newGRPC: %w", err)
	}

	if err = s.rESTGW.initRESTGW(ctx, s.config, pki); err != nil {
		return nil, fmt.Errorf("s.newRESTGW: %w", err)
	}

	s.cliGRPC = client.NewGRPC(client.NewConfigClient(), pki)

	s.cliHTTP = client.NewHTTP(client.NewConfigClient(), pki)

	return s, nil
}

// ServiceRegistrar to register service instance on gRPC
func (s *Server) ServiceRegistrar() *grpc.Server {
	return s.gRPC.grpcServer
}

func (s *Server) Serve(ctx context.Context) {
	go s.serve(ctx)
	runtime.Gosched()
}

func (s *Server) serve(ctx context.Context) {
	go s.gRPC.serveGRPC()
	runtime.Gosched()
	s.cliGRPC.DialContext(ctx)

	go s.rESTGW.serveRESTGW(ctx)
	runtime.Gosched()

	<-ctx.Done()
	s.Close(ctx)
}

func (s *Server) Close(ctx context.Context) {
	s.rESTGW.closeRESTGW(ctx)
	s.cliGRPC.Close()
	s.gRPC.closeGRPC()
}

func (s *Server) Valid() bool {
	return s.config.Valid() &&
		s.inst.Valid() &&
		s.gRPC != nil &&
		s.gRPC.netListener != nil &&
		s.rESTGW != nil &&
		s.rESTGW.httpServer != nil &&
		s.cliGRPC != nil &&
		s.cliGRPC.ValidConnection()
}

func (s *Server) ClientGRPC() *client.GRPC {
	return s.cliGRPC
}

func (s *Server) ClientHTTP() *client.HTTP {
	return s.cliHTTP
}
