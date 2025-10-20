package api

import (
	"strconv"
)

func parseIntOrDefault(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}

	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}

	return int(val)
}
