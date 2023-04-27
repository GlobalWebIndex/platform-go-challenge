package cfg

import (
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Configuration defines the struct with the configurable options passed from
// the environment.
type Configuration struct {
	Address       string
	Port          string
	JwtSecret     []byte
	AdminUser     string
	AdminPass     string
	StorageOption string
	CacheOption   string
	Profiling     bool
	Metrics       bool
}

// ReadConfig retrieves the configuration from the .env file.
func ReadConfig(file string) (*Configuration, error) {
	viper.SetConfigType("yaml")
	if file != "" {
		viper.SetConfigFile(file)
	} else {
		_, file, _, _ := runtime.Caller(1)
		if strings.Contains(file, "test") {
			viper.SetConfigFile("../config.yaml")
		} else {
			viper.SetConfigFile("config.yaml")
		}
	}
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := &Configuration{
		Address:       viper.GetString("address"),
		Port:          strconv.Itoa(viper.GetInt("port")),
		JwtSecret:     []byte(viper.GetString("tokenSecret")),
		AdminUser:     viper.GetString("adminUser"),
		AdminPass:     viper.GetString("adminPass"),
		StorageOption: viper.GetString("storage"),
		CacheOption:   viper.GetString("cache"),
		Profiling:     viper.GetBool("profiler"),
		Metrics:       viper.GetBool("metrics"),
	}
	return config, nil
}
