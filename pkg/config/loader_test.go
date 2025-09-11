package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Unset environment variable to ensure defaults are used
	os.Setenv("JWT_SECRET", "secret")

	// Call the function under test
	cfg, err := LoadConfig()

	// Assert that there's no error and the config is not nil
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Assert the default values
	assert.Equal(t, "0.0.0.0:8080", cfg.ServerAddr)
	assert.Equal(t, 30*time.Second, cfg.ShutdownTimeout)
	assert.Equal(t, 5*time.Second, cfg.ReadHeaderTimeout)
	assert.Equal(t, 10*time.Second, cfg.ReadTimeout)
	assert.Equal(t, 10*time.Second, cfg.WriteTimeout)
	assert.Equal(t, 60*time.Second, cfg.IdleTimeout)
}

func TestLoadConfig_EnvVariable(t *testing.T) {
	// Set a mock environment variable
	t.Setenv("JWT_SECRET", "test_env_secret")

	// Call the function under test
	cfg, err := LoadConfig()

	// Assert that there's no error and the JWTSecret is set from the env var
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "test_env_secret", cfg.JWTSecret)
}

func TestLoadConfig_ValidationFailure(t *testing.T) {
	os.Unsetenv("JWT_SECRET")

	cfg, err := LoadConfig()
	assert.Nil(t, cfg)
	assert.ErrorContains(t, err, "JWT_SECRET environment variable or config value is required")
}

func TestLoadConfig_ConfigFileFound(t *testing.T) {
	// Create a temporary directory and a mock config file
	tempDir := t.TempDir()
	viper.AddConfigPath(tempDir)

	mockConfig := map[string]interface{}{
		"server_addr": "127.0.0.1:9090",
		"jwt_secret":  "test_config_secret",
	}
	file, err := json.Marshal(mockConfig)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(tempDir, "config.json"), file, 0644)
	assert.NoError(t, err)

	// Call the function under test
	cfg, err := LoadConfig()

	// Assert that there's no error and values are from the file
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "127.0.0.1:9090", cfg.ServerAddr)
	assert.Equal(t, "test_config_secret", cfg.JWTSecret)

	// Cleanup the temporary directory
	defer os.RemoveAll(tempDir)
}
