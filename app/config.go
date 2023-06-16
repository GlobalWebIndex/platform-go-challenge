package app

import (
	"x-gwi/app/auth"
	"x-gwi/app/client"
	"x-gwi/app/instance"
	"x-gwi/app/server"
	"x-gwi/app/storage"
)

// Config represents global app configuration.
type Config struct {
	Inst    *instance.Instance
	Auth    *auth.ConfigAuth
	Storage *storage.ConfigAppStorage
	Server  *server.ConfigServer
	Client  *client.ConfigClient
	// Services
}

func newConfig() *Config {
	return &Config{
		Inst:    instance.NewInstance(),
		Auth:    auth.NewConfigAuth(),
		Server:  server.NewConfigServer(),
		Client:  client.NewConfigClient(),
		Storage: storage.NewConfigStorage(),
	}
}

func (c *Config) Valid() bool {
	return c.Inst.Valid() &&
		c.Auth.Valid() &&
		c.Storage.Valid() &&
		c.Server.Valid() &&
		c.Client.Valid()
}
