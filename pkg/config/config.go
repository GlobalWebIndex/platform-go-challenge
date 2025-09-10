package config

import (
	"time"
)

// Config represents the application's configuration settings.
// 'mapstructure' tags are used to map fields from Viper to the struct.
type Config struct {
	JWTSecret         string        `mapstructure:"jwt_secret"`
	ServerAddr        string        `mapstructure:"server_addr"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdown_timeout"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
	ReadTimeout       time.Duration `mapstructure:"read_timeout"`
	WriteTimeout      time.Duration `mapstructure:"write_timeout"`
	IdleTimeout       time.Duration `mapstructure:"idle_timeout"`
}

func (c *Config) GetJWTSecret() string {
	return c.JWTSecret
}

func (c *Config) GetServerAddr() string {
	return c.ServerAddr
}

func (c *Config) GetShutdownTimeout() time.Duration {
	return c.ShutdownTimeout
}

func (c *Config) GetReadHeaderTimeout() time.Duration {
	return c.ReadHeaderTimeout
}

func (c *Config) GetReadTimeout() time.Duration {
	return c.ReadTimeout
}

func (c *Config) GetWriteTimeout() time.Duration {
	return c.WriteTimeout
}

func (c *Config) GetIdleTimeout() time.Duration {
	return c.IdleTimeout
}
