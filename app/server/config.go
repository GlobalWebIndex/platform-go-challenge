package server

import "x-gwi/internal/env"

type ConfigServer struct {
	GRPC   ConfigGRPC
	RESTGW ConfigRESTGW
}

type ConfigGRPC struct {
	Address string
}

type ConfigRESTGW struct {
	Address string
}

func NewConfigServer() *ConfigServer {
	return &ConfigServer{
		GRPC: ConfigGRPC{
			Address: env.Env("SERVER_GRPC_ADDRESS", ":9090"),
		},
		RESTGW: ConfigRESTGW{
			Address: env.Env("SERVER_RESTGW_ADDRESS", ":9080"),
		},
	}
}

func (s *ConfigServer) Valid() bool {
	return s.GRPC.Address != "" &&
		s.RESTGW.Address != ""
}
