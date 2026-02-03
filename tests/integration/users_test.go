package integration

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/crypto/bcrypt"

	repo "go-favorites-app/internal/adapter/storage/postgres"
	"go-favorites-app/internal/core/service"
)

func TestUserIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	// 1. Start Postgres Container
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start postgres: %v", err)
	}
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate postgres: %v", err)
		}
	}()

	pgConnStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get pg connection string: %v", err)
	}

	// 2. Setup Database Connection
	dbPool, err := pgxpool.New(ctx, pgConnStr)
	if err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}
	defer dbPool.Close()

	// 3. Init Schema
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	`
	if _, err := dbPool.Exec(ctx, schema); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	// 4. Initialize Service
	userRepo := repo.NewUserRepository(dbPool)
	authService := service.NewAuthService(userRepo, "test-secret")

	// 5. Test Scenarios
	t.Run("SignUp Success", func(t *testing.T) {
		email := "newuser@example.com"
		password := "securePass123"

		err := authService.SignUp(ctx, email, password)
		if err != nil {
			t.Fatalf("SignUp failed: %v", err)
		}

		// Verify in DB
		var hash string
		err = dbPool.QueryRow(ctx, "SELECT password_hash FROM users WHERE email = $1", email).Scan(&hash)
		if err != nil {
			t.Fatalf("failed to query user: %v", err)
		}

		// Verify not plaintext
		if hash == password {
			t.Fatal("password stored in plaintext")
		}

		// Verify hash matches
		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
			t.Fatalf("stored hash does not match password: %v", err)
		}
	})

	t.Run("SignUp Duplicate Email", func(t *testing.T) {
		email := "duplicate@example.com"
		password := "password"

		// First creation
		err := authService.SignUp(ctx, email, password)
		if err != nil {
			t.Fatalf("first signup failed: %v", err)
		}

		// Second creation
		err = authService.SignUp(ctx, email, password)
		if err == nil {
			t.Fatal("expected error on duplicate email, got nil")
		}
	})

	t.Run("Login Success", func(t *testing.T) {
		email := "loginuser@example.com"
		password := "loginPass"

		// Setup user
		if err := authService.SignUp(ctx, email, password); err != nil {
			t.Fatalf("signup failed: %v", err)
		}

		token, err := authService.Login(ctx, email, password)
		if err != nil {
			t.Fatalf("login failed: %v", err)
		}
		if token == "" {
			t.Fatal("expected token, got empty string")
		}
	})

	t.Run("Login Failure - Wrong Password", func(t *testing.T) {
		email := "wrongpass@example.com"
		password := "correctPass"

		if err := authService.SignUp(ctx, email, password); err != nil {
			t.Fatalf("signup failed: %v", err)
		}

		_, err := authService.Login(ctx, email, "wrongPass")
		if err == nil {
			t.Fatal("expected error on wrong password, got nil")
		}
	})

	t.Run("Login Failure - Non-existent User", func(t *testing.T) {
		_, err := authService.Login(ctx, "ghost@example.com", "password")
		if err == nil {
			t.Fatal("expected error on missing user, got nil")
		}
	})
}
