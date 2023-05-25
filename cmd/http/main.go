package main

import (
	"github.com/Kercyn/crud_template/configs"
	"github.com/Kercyn/crud_template/internal/app"
	"github.com/spf13/viper"
	"log"
)

func main() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.WatchConfig()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("cannot read config file: %v", err)
	}

	var config configs.Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("cannot unmarshal config file: %v", err)
	}

	if err := app.Execute(config); err != nil {
		log.Fatalf("error executing app: %v", err)
	}
}
