package client

import (
	"fmt"
	"net/http"

	"x-gwi/app/x/env"
)

type ServiceMethod string

const (
	UserAccountCreate ServiceMethod = "UserAccountCreate"
	UserAccountGet    ServiceMethod = "UserAccountGet"
)

// HTTPMethodAndPath returns http method and http path assigned to provided ServiceMethod
func (sm *ServiceMethod) HTTPMethodAndPath() (string, string, error) {
	switch *sm {
	case UserAccountCreate:
		return http.MethodPost, "/proto.useraccountpb.UserAccountService/Create", nil
	case UserAccountGet:
		return http.MethodPost, "/proto.useraccountpb.UserAccountService/Get", nil
	}

	return "", "", fmt.Errorf("unsupported ServiceMethod: %v", sm) //nolint:goerr113
}

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
