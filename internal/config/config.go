package config

import (
	"os"
	"path/filepath"
)

func GetConfigPath() (string, error) {
	executablePath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	configPath := filepath.Join(filepath.Dir(executablePath), "config")
	return configPath, nil
}
