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
	s := os.Getenv(envField)
	if s == "" {
		if defValues == "" {
			return []string{}
		}

		s = defValues
	}

	s = strings.TrimSpace(s)
	if strings.Contains(s, " ") {
		s = strings.ReplaceAll(s, " ", ",")
	}

	return strings.Split(s, ",")
}
