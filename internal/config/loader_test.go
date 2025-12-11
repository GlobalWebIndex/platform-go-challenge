package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	_ = os.Unsetenv(envHTTPAddr)

	cfg := Load()
	if cfg.HTTPAddr != defaultHTTPAddr {
		t.Fatalf("expected HTTPAddr %q, got %q", defaultHTTPAddr, cfg.HTTPAddr)
	}
}

func TestLoad_FromEnv(t *testing.T) {
	const expected = ":9090"
	_ = os.Setenv(envHTTPAddr, expected)
	defer os.Unsetenv(envHTTPAddr)

	cfg := Load()
	if cfg.HTTPAddr != expected {
		t.Fatalf("expected HTTPAddr %q, got %q", expected, cfg.HTTPAddr)
	}
}
