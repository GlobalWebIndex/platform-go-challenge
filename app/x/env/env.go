package env

import (
	"os"
	"strings"
)

func Env(envField, defValue string) string {
	if s := os.Getenv(envField); s != "" {
		return s
	}

	return defValue
}

func Envs(envField, defValues string) []string {
	if os.Getenv(envField) != "" {
		return strings.Split(os.Getenv(envField), ",")
	}

	return strings.Split(defValues, ",")
}
