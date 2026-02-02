package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	DefaultHTTPAddr           = ":8080"
	DefaultPublicKeyPath      = "public.pem"
	DefaultPrivateKeyPath     = "private.pem"
	DefaultJWTIssuer          = "favourites-api"
	DefaultJWTAudience        = "web"
	DefaultPaginationPage     = 1
	DefaultPaginationLimit    = 20
	DefaultPaginationMaxLimit = 100
)

type AppConfig struct {
	HTTPAddr   string
	Auth       AuthConfig
	Pagination PaginationConfig
	EnableDocs bool
}

type AuthConfig struct {
	PublicKeyPath  string
	PrivateKeyPath string
	Issuer         string
	Audience       string
}

type PaginationConfig struct {
	DefaultPage  int
	DefaultLimit int
	MaxLimit     int
}

func Load() (AppConfig, error) {
	return LoadWithEnv("")
}

// LoadWithEnv loads configuration using (in order):
// 1) defaults,
// 2) config file at configs/<env>.yaml (if present),
// 3) environment variables overrides.
// If env is empty, APP_ENV is used, falling back to "local".
func LoadWithEnv(env string) (AppConfig, error) {
	if env == "" {
		env = envOrDefault("APP_ENV", "local")
	}

	cfg := AppConfig{
		HTTPAddr: DefaultHTTPAddr,
		Auth: AuthConfig{
			PublicKeyPath:  DefaultPublicKeyPath,
			PrivateKeyPath: DefaultPrivateKeyPath,
			Issuer:         DefaultJWTIssuer,
			Audience:       DefaultJWTAudience,
		},
		Pagination: PaginationConfig{
			DefaultPage:  DefaultPaginationPage,
			DefaultLimit: DefaultPaginationLimit,
			MaxLimit:     DefaultPaginationMaxLimit,
		},
		EnableDocs: false,
	}

	if err := applyFileConfig(&cfg, env); err != nil {
		return AppConfig{}, err
	}

	// Env overrides
	cfg.HTTPAddr = envOrDefault("HTTP_ADDR", cfg.HTTPAddr)
	cfg.Auth.PublicKeyPath = envOrDefault("PUBLIC_KEY_PATH", cfg.Auth.PublicKeyPath)
	cfg.Auth.PrivateKeyPath = envOrDefault("PRIVATE_KEY_PATH", cfg.Auth.PrivateKeyPath)
	cfg.Auth.Issuer = envOrDefault("JWT_ISSUER", cfg.Auth.Issuer)
	cfg.Auth.Audience = envOrDefault("JWT_AUDIENCE", cfg.Auth.Audience)

	if v, ok, err := intFromEnv("PAGINATION_DEFAULT_LIMIT"); err != nil {
		return AppConfig{}, err
	} else if ok {
		cfg.Pagination.DefaultLimit = v
	}
	if v, ok, err := intFromEnv("PAGINATION_MAX_LIMIT"); err != nil {
		return AppConfig{}, err
	} else if ok {
		cfg.Pagination.MaxLimit = v
	}
	if v, ok, err := intFromEnv("PAGINATION_DEFAULT_PAGE"); err != nil {
		return AppConfig{}, err
	} else if ok {
		cfg.Pagination.DefaultPage = v
	}

	if b := boolFromEnv("ENABLE_DOCS"); b {
		cfg.EnableDocs = b
	}

	if err := validate(cfg); err != nil {
		return AppConfig{}, err
	}

	return cfg, nil
}

type fileConfig struct {
	HTTPAddr   string              `yaml:"httpAddr"`
	Auth       *fileAuthConfig     `yaml:"auth"`
	Pagination *filePaginationConf `yaml:"pagination"`
	EnableDocs *bool               `yaml:"enableDocs"`
}

type fileAuthConfig struct {
	PublicKeyPath  string `yaml:"publicKeyPath"`
	PrivateKeyPath string `yaml:"privateKeyPath"`
	Issuer         string `yaml:"issuer"`
	Audience       string `yaml:"audience"`
}

type filePaginationConf struct {
	DefaultPage  int `yaml:"defaultPage"`
	DefaultLimit int `yaml:"defaultLimit"`
	MaxLimit     int `yaml:"maxLimit"`
}

func applyFileConfig(cfg *AppConfig, env string) error {
	path := filepath.Join("configs", env+".yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read config file %s: %w", path, err)
	}

	var fc fileConfig
	if err := yaml.Unmarshal(data, &fc); err != nil {
		return fmt.Errorf("parse config file %s: %w", path, err)
	}

	if fc.HTTPAddr != "" {
		cfg.HTTPAddr = fc.HTTPAddr
	}
	if fc.Auth != nil {
		if fc.Auth.PublicKeyPath != "" {
			cfg.Auth.PublicKeyPath = fc.Auth.PublicKeyPath
		}
		if fc.Auth.PrivateKeyPath != "" {
			cfg.Auth.PrivateKeyPath = fc.Auth.PrivateKeyPath
		}
		if fc.Auth.Issuer != "" {
			cfg.Auth.Issuer = fc.Auth.Issuer
		}
		if fc.Auth.Audience != "" {
			cfg.Auth.Audience = fc.Auth.Audience
		}
	}
	if fc.Pagination != nil {
		if fc.Pagination.DefaultPage > 0 {
			cfg.Pagination.DefaultPage = fc.Pagination.DefaultPage
		}
		if fc.Pagination.DefaultLimit > 0 {
			cfg.Pagination.DefaultLimit = fc.Pagination.DefaultLimit
		}
		if fc.Pagination.MaxLimit > 0 {
			cfg.Pagination.MaxLimit = fc.Pagination.MaxLimit
		}
	}
	if fc.EnableDocs != nil {
		cfg.EnableDocs = *fc.EnableDocs
	}

	return nil
}

func validate(cfg AppConfig) error {
	if strings.TrimSpace(cfg.Auth.Issuer) == "" {
		return fmt.Errorf("jwt issuer is required")
	}
	if strings.TrimSpace(cfg.Auth.Audience) == "" {
		return fmt.Errorf("jwt audience is required")
	}
	if strings.TrimSpace(cfg.Auth.PublicKeyPath) == "" {
		return fmt.Errorf("public key path is required")
	}
	if cfg.Pagination.DefaultPage <= 0 {
		return fmt.Errorf("pagination default page must be > 0")
	}
	if cfg.Pagination.DefaultLimit <= 0 {
		return fmt.Errorf("pagination default limit must be > 0")
	}
	if cfg.Pagination.MaxLimit <= 0 {
		return fmt.Errorf("pagination max limit must be > 0")
	}
	if cfg.Pagination.DefaultLimit > cfg.Pagination.MaxLimit {
		return fmt.Errorf("pagination default limit (%d) cannot exceed max limit (%d)", cfg.Pagination.DefaultLimit, cfg.Pagination.MaxLimit)
	}
	return nil
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// intFromEnv returns (value, true, nil) if key exists and is a valid int, (0, false, nil) if unset.
func intFromEnv(key string) (int, bool, error) {
	v := os.Getenv(key)
	if v == "" {
		return 0, false, nil
	}
	num, err := strconv.Atoi(v)
	if err != nil {
		return 0, false, fmt.Errorf("%s must be an integer: %w", key, err)
	}
	return num, true, nil
}

func boolFromEnv(key string) bool {
	switch strings.ToLower(os.Getenv(key)) {
	case "1", "true", "yes", "y", "on":
		return true
	default:
		return false
	}
}
