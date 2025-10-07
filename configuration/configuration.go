package configuration

import (
	"os"
)

var JwtSecret = []byte(func() string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	return "gwi-jwt-secret"
}())
