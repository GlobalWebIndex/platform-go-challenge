package storage

import "x-gwi/app/x/env"

type ConfigAppStorage struct {
	AQL ConfigAQL
	// CQL ConfigCQL
	// SQL ConfigSQL
}

type ConfigAQL struct {
	Credentials CredentialsBasic
	// Endpoints holds 1 or more URL's used to connect to the database.
	// In case of a connection to an ArangoDB cluster, you must provide the URL's of all coordinators.
	Endpoints []string
}

/*
type ConfigCQL struct {
	Credentials CredentialsBasic
	// The supplied hosts are used to initially connect to the cluster then the rest of the ring will
	// be automatically discovered. It is recommended to use the value set in the Cassandra config for
	// broadcast_address or listen_address, an IP address not a domain name. This is because events
	// from Cassandra will use the configured IP address, which is used to index connected hosts.
	// If the domain name specified resolves to more than 1 IP address then the driver may connect
	// multiple times to the same host, and will not mark the node being down or up from events.
	Endpoints []string
}
*/

/*
type ConfigSQL struct {
	Credentials CredentialsBasic
	// db, err := sql.Open("pgx", "postgresql://test:123@localhost:5432/benchmark?sslmode=disable")
	Endpoints []string
}
*/

type CredentialsBasic struct {
	UserName string
	PassWord string
	// PassHash string
}

func NewConfigStorage() *ConfigAppStorage {
	return &ConfigAppStorage{
		AQL: ConfigAQL{
			Credentials: CredentialsBasic{
				UserName: env.Env("APP_STORAGE_AQL_USERNAME", "arango"),
				PassWord: env.Env("APP_STORAGE_AQL_PASSWORD", "arango"),
			},
			Endpoints: env.Envs("APP_STORAGE_AQL_ENDPOINTS", "http://192.168.1.15:8529,http://localhost:8529"),
		},
		/* 		CQL: ConfigCQL{
		   			Credentials: CredentialsBasic{
		   				UserName: env.Env("APP_STORAGE_CQL_USERNAME", "cassandra"),
		   				PassWord: env.Env("APP_STORAGE_CQL_PASSWORD", "cassandra"),
		   			},
		   			Endpoints: env.Envs("APP_STORAGE_CQL_ENDPOINTS", "127.0.0.1,192.168.1.15"),
		   		},
		   		SQL: ConfigSQL{
		   			Credentials: CredentialsBasic{
		   				UserName: env.Env("APP_STORAGE_SQL_USERNAME", "postgresql"),
		   				PassWord: env.Env("APP_STORAGE_SQL_PASSWORD", "postgresql"),
		   			},
		   			Endpoints: env.Envs("APP_STORAGE_SQL_ENDPOINTS", ""),
		   		}, */
	}
}

func (s *ConfigAppStorage) Valid() bool {
	return true
}
