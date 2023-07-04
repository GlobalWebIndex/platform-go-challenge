package client

import (
	"x-gwi/app/x/env"
)

type ConfigClient struct {
	GRPC   ConfigGRPC
	RESTGW ConfigRESTGW
}

type ConfigGRPC struct {
	// Network string
	Address string
}

type ConfigRESTGW struct {
	Address string
}

func NewConfigClient() *ConfigClient {
	return &ConfigClient{
		GRPC: ConfigGRPC{
			Address: env.Env("SERVER_GRPC_ADDRESS", ":9090"),
		},
		RESTGW: ConfigRESTGW{
			Address: env.Env("SERVER_RESTGW_ADDRESS", ":9080"),
		},
	}
}

func (s *ConfigClient) Valid() bool {
	return s.GRPC.Address != "" &&
		s.RESTGW.Address != ""
}
