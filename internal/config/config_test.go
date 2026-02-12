package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setDBEnv sets all required database environment variables for testing.
func setDBEnv(t *testing.T) {
	t.Helper()
	t.Setenv("POSTGRES_HOST", "localhost")
	t.Setenv("POSTGRES_PORT", "5432")
	t.Setenv("POSTGRES_USER", "testuser")
	t.Setenv("POSTGRES_PASSWORD", "testpass")
	t.Setenv("POSTGRES_DB", "testdb")
}

func TestLoad_ErrorWhenNoFileAndNoEnv(t *testing.T) {
	t.Setenv("CONFIG_PATH", filepath.Join(t.TempDir(), "nonexistent.yaml"))
	t.Setenv("API_PORT", "")
	t.Setenv("HEALTH_PORT", "")
	setDBEnv(t)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when no config source provides required ports")
	}
}

func TestLoad_ErrorWhenPartialConfig(t *testing.T) {
	path := writeTempConfig(t, `api_port: "9000"`)
	t.Setenv("CONFIG_PATH", path)
	t.Setenv("API_PORT", "")
	t.Setenv("HEALTH_PORT", "")
	setDBEnv(t)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when health_port is missing")
	}
}

func TestLoad_FromYAMLFile(t *testing.T) {
	path := writeTempConfig(t, `api_port: "9000"
health_port: "9001"
`)
	t.Setenv("CONFIG_PATH", path)
	t.Setenv("API_PORT", "")
	t.Setenv("HEALTH_PORT", "")
	setDBEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.APIPort != "9000" {
		t.Errorf("expected APIPort=9000, got %s", cfg.APIPort)
	}
	if cfg.HealthPort != "9001" {
		t.Errorf("expected HealthPort=9001, got %s", cfg.HealthPort)
	}
}

func TestLoad_EnvOverridesFile(t *testing.T) {
	path := writeTempConfig(t, `api_port: "9000"
health_port: "9001"
`)
	t.Setenv("CONFIG_PATH", path)
	t.Setenv("API_PORT", "7000")
	t.Setenv("HEALTH_PORT", "7001")
	setDBEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.APIPort != "7000" {
		t.Errorf("expected APIPort=7000 (env override), got %s", cfg.APIPort)
	}
	if cfg.HealthPort != "7001" {
		t.Errorf("expected HealthPort=7001 (env override), got %s", cfg.HealthPort)
	}
}

func TestLoad_PartialEnvOverride(t *testing.T) {
	path := writeTempConfig(t, `api_port: "9000"
health_port: "9001"
`)
	t.Setenv("CONFIG_PATH", path)
	t.Setenv("API_PORT", "7000")
	t.Setenv("HEALTH_PORT", "")
	setDBEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.APIPort != "7000" {
		t.Errorf("expected APIPort=7000 (env override), got %s", cfg.APIPort)
	}
	if cfg.HealthPort != "9001" {
		t.Errorf("expected HealthPort=9001 (from file), got %s", cfg.HealthPort)
	}
}

func TestLoad_InvalidYAML_ReturnsError(t *testing.T) {
	path := writeTempConfig(t, "{{invalid yaml}}")
	t.Setenv("CONFIG_PATH", path)
	t.Setenv("API_PORT", "")
	t.Setenv("HEALTH_PORT", "")
	setDBEnv(t)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestLoad_DBConfigPopulated(t *testing.T) {
	path := writeTempConfig(t, `api_port: "9000"
health_port: "9001"
`)
	t.Setenv("CONFIG_PATH", path)
	t.Setenv("API_PORT", "")
	t.Setenv("HEALTH_PORT", "")
	setDBEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.DBHost != "localhost" {
		t.Errorf("expected DBHost=localhost, got %s", cfg.DBHost)
	}
	if cfg.DBPort != "5432" {
		t.Errorf("expected DBPort=5432, got %s", cfg.DBPort)
	}
	if cfg.DBUser != "testuser" {
		t.Errorf("expected DBUser=testuser, got %s", cfg.DBUser)
	}
	if cfg.DBPassword != "testpass" {
		t.Errorf("expected DBPassword=testpass, got %s", cfg.DBPassword)
	}
	if cfg.DBName != "testdb" {
		t.Errorf("expected DBName=testdb, got %s", cfg.DBName)
	}
}

func TestLoad_MissingDBEnvVars(t *testing.T) {
	path := writeTempConfig(t, `api_port: "9000"
health_port: "9001"
`)

	requiredVars := []string{
		"POSTGRES_HOST",
		"POSTGRES_PORT",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_DB",
	}

	for _, missing := range requiredVars {
		t.Run("missing_"+missing, func(t *testing.T) {
			t.Setenv("CONFIG_PATH", path)
			t.Setenv("API_PORT", "")
			t.Setenv("HEALTH_PORT", "")

			// Set all DB vars, then clear the one under test
			setDBEnv(t)
			t.Setenv(missing, "")

			_, err := Load()
			if err == nil {
				t.Fatalf("expected error when %s is missing", missing)
			}
			if !strings.Contains(err.Error(), missing) {
				t.Errorf("expected error to mention %s, got: %v", missing, err)
			}
		})
	}
}

func TestDSN(t *testing.T) {
	cfg := &Config{
		DBHost:     "dbhost",
		DBPort:     "5432",
		DBUser:     "myuser",
		DBPassword: "mypass",
		DBName:     "mydb",
	}

	want := "host=dbhost port=5432 user=myuser password=mypass dbname=mydb sslmode=disable"
	got := cfg.PostgresConnString()
	if got != want {
		t.Errorf("DSN() = %q, want %q", got, want)
	}
}

func TestAddr_Methods(t *testing.T) {
	cfg := &Config{APIPort: "3000", HealthPort: "3001"}

	if cfg.APIAddr() != ":3000" {
		t.Errorf("expected :3000, got %s", cfg.APIAddr())
	}
	if cfg.HealthAddr() != ":3001" {
		t.Errorf("expected :3001, got %s", cfg.HealthAddr())
	}
}

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return path
}
