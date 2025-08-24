// Package config loads configuration from the config.yml file 
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Logs LogsConfig
	Database DatabaseConfig
	Server ServerConfig
}

type LogsConfig struct {
	LogsDir string `mapstructure:"logsDir"`
} 

type DatabaseConfig struct {
	MongoURI string `mapstructure:"mongoURI"`
	FavCollection string `mapstructure:"favCollection"`
	AssetsCollection string `mapstructure:"assetsCollection"`
}

type ServerConfig struct {
	TLS bool `mapstructure:"tls"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	// Config file setup
	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.AddConfigPath("/app/config") 

	// Read the config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable ro decode into struct: %v", err)
	}

	return &cfg, nil

	//return &Config{}, nil
}




