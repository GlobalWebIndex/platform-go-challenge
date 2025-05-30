package postgres_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/uptrace/bun/driver/pgdriver"
	"pgc/domain"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	internalpostgres "pgc/persistence/postgres"
)

func TestFavoriteRepository_SearchByUserID(t *testing.T) {
	tests := []struct {
		name              string
		inputUserID       int
		expectedResult    []domain.Favorite
		expectedError     error
		useMockRepository bool
	}{
		{
			name:        "valid user ID with favorites",
			inputUserID: 1,
			expectedResult: []domain.Favorite{
				{
					ID:          1,
					UserID:      1,
					AssetID:     1,
					Description: "Primary sales chart",
				},
				{
					ID:          2,
					UserID:      1,
					AssetID:     4,
					Description: "Key Q4 insight",
				},
			},
			expectedError: nil,
		},
		{
			name:           "valid user ID with no favorites",
			inputUserID:    4,
			expectedResult: []domain.Favorite{},
			expectedError:  nil,
		},
		{
			name:           "invalid user ID - zero",
			inputUserID:    0,
			expectedResult: nil,
			expectedError:  domain.ErrInvalidUserID,
		},
		{
			name:           "invalid user ID - negative",
			inputUserID:    -1,
			expectedResult: nil,
			expectedError:  domain.ErrInvalidUserID,
		},
		{
			name:              "database connection error",
			inputUserID:       1,
			expectedResult:    nil,
			expectedError:     domain.ErrInternalServerError,
			useMockRepository: true,
		},
		{
			name:           "user ID not found",
			inputUserID:    999,
			expectedResult: []domain.Favorite{},
			expectedError:  nil, // No error, just an empty result
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result []domain.Favorite
			var err error
			if tt.useMockRepository {
				result, err = mockRepo.SearchByUserID(context.Background(), tt.inputUserID)
			} else {
				result, err = favoriteRepo.SearchByUserID(context.Background(), tt.inputUserID)
			}
			// Assert results
			if tt.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedError.Error(), err.Error())
				require.Equal(t, tt.expectedResult, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
			}

		})
	}
}

func TestFavoriteRepository_Create(t *testing.T) {
	tests := []struct {
		name              string
		input             domain.Favorite
		expectedResult    domain.Favorite
		expectedError     error
		useMockRepository bool
	}{
		{
			name: "valid favorite creation",
			input: domain.Favorite{
				UserID:      1,
				AssetID:     2,
				Description: "Test favorite",
			},
			expectedResult: domain.Favorite{
				ID:          6,
				UserID:      1,
				AssetID:     2,
				Description: "Test favorite",
			},
			expectedError: nil,
		},
		{
			name: "valid favorite with empty description",
			input: domain.Favorite{
				UserID:      1,
				AssetID:     3,
				Description: "",
			},
			expectedResult: domain.Favorite{
				ID:          7,
				UserID:      1,
				AssetID:     3,
				Description: "",
			},
			expectedError: nil,
		},
		{
			name: "invalid user ID - zero",
			input: domain.Favorite{
				UserID:      0,
				AssetID:     2,
				Description: "Test favorite",
			},
			expectedResult: domain.Favorite{},
			expectedError:  domain.ErrInvalidUserID,
		},
		{
			name: "invalid user ID - negative",
			input: domain.Favorite{
				UserID:      -1,
				AssetID:     2,
				Description: "Test favorite",
			},
			expectedResult: domain.Favorite{},
			expectedError:  domain.ErrInvalidUserID,
		},
		{
			name: "invalid asset ID - zero",
			input: domain.Favorite{
				UserID:      1,
				AssetID:     0,
				Description: "Test favorite",
			},
			expectedResult: domain.Favorite{},
			expectedError:  domain.ErrInvalidAssetID,
		},
		{
			name: "invalid asset ID - negative",
			input: domain.Favorite{
				UserID:      1,
				AssetID:     -1,
				Description: "Test favorite",
			},
			expectedResult: domain.Favorite{},
			expectedError:  domain.ErrInvalidAssetID,
		},
		{
			name: "duplicate favorite error",
			input: domain.Favorite{
				UserID:      1,
				AssetID:     1,
				Description: "Duplicate favorite",
			},
			expectedResult: domain.Favorite{},
			expectedError:  domain.ErrDuplicateFavorite,
		},
		{
			name: "foreign key violation - invalid user ID",
			input: domain.Favorite{
				UserID:      999,
				AssetID:     2,
				Description: "Test favorite",
			},
			expectedResult: domain.Favorite{},
			expectedError:  domain.ErrInvalidUserID,
		},
		{
			name: "foreign key violation - invalid asset ID",
			input: domain.Favorite{
				UserID:      1,
				AssetID:     999,
				Description: "Test favorite",
			},
			expectedResult: domain.Favorite{},
			expectedError:  domain.ErrInvalidAssetID,
		},
		{
			name: "database connection error",
			input: domain.Favorite{
				UserID:      1,
				AssetID:     2,
				Description: "Test favorite",
			},
			expectedResult:    domain.Favorite{},
			expectedError:     domain.ErrInternalServerError,
			useMockRepository: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result domain.Favorite
			var err error
			if tt.useMockRepository {
				result, err = mockRepo.Insert(context.Background(), tt.input)
			} else {
				result, err = favoriteRepo.Insert(context.Background(), tt.input)

			}

			// Assert results
			if tt.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedError.Error(), err.Error())
				require.Equal(t, tt.expectedResult, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestFavoriteRepository_UpdateDescription(t *testing.T) {
	tests := []struct {
		name              string
		inputFavoriteID   int
		inputUserID       int
		inputDescription  string
		expectedResult    domain.Favorite
		expectedError     error
		useMockRepository bool
	}{
		{
			name:             "valid favorite description update",
			inputFavoriteID:  1,
			inputUserID:      1,
			inputDescription: "Updated description",
			expectedResult: domain.Favorite{
				ID:          1,
				UserID:      1,
				AssetID:     1,
				Description: "Updated description",
			},
			expectedError: nil,
		},
		{
			name:             "valid favorite description update with empty description",
			inputFavoriteID:  1,
			inputUserID:      1,
			inputDescription: "",
			expectedResult: domain.Favorite{
				ID:          1,
				UserID:      1,
				AssetID:     1,
				Description: "",
			},
			expectedError: nil,
		},
		{
			name:             "invalid favorite ID - zero",
			inputFavoriteID:  0,
			inputUserID:      1,
			inputDescription: "Test description",
			expectedResult:   domain.Favorite{},
			expectedError:    domain.ErrInvalidFavoriteID,
		},
		{
			name:             "invalid favorite ID - negative",
			inputFavoriteID:  -1,
			inputUserID:      1,
			inputDescription: "Test description",
			expectedResult:   domain.Favorite{},
			expectedError:    domain.ErrInvalidFavoriteID,
		},
		{
			name:             "invalid user ID - zero",
			inputFavoriteID:  1,
			inputUserID:      0,
			inputDescription: "Test description",
			expectedResult:   domain.Favorite{},
			expectedError:    domain.ErrInvalidUserID,
		},
		{
			name:             "invalid user ID - negative",
			inputFavoriteID:  1,
			inputUserID:      -1,
			inputDescription: "Test description",
			expectedResult:   domain.Favorite{},
			expectedError:    domain.ErrInvalidUserID,
		},
		{
			name:              "database connection error",
			inputFavoriteID:   1,
			inputUserID:       1,
			inputDescription:  "Test description",
			expectedResult:    domain.Favorite{},
			expectedError:     domain.ErrInternalServerError,
			useMockRepository: true,
		},
		{
			name:             "favorite not found",
			inputFavoriteID:  1,
			inputUserID:      2,
			inputDescription: "Test description",
			expectedResult:   domain.Favorite{},
			expectedError:    domain.ErrFavoriteNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result domain.Favorite
			var err error
			if tt.useMockRepository {
				result, err = mockRepo.UpdateDescription(context.Background(), tt.inputFavoriteID, tt.inputUserID, tt.inputDescription)
			} else {
				result, err = favoriteRepo.UpdateDescription(context.Background(), tt.inputFavoriteID, tt.inputUserID, tt.inputDescription)
			}

			// Assert results
			if tt.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedError.Error(), err.Error())
				require.Equal(t, tt.expectedResult, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
			}

		})
	}
}

func TestFavoriteRepository_Remove(t *testing.T) {
	tests := []struct {
		name              string
		inputID           int
		inputUserID       int
		expectedError     error
		useMockRepository bool
	}{
		{
			name:          "valid favorite removal",
			inputUserID:   1,
			inputID:       1,
			expectedError: nil,
		},
		{
			name:          "invalid user ID - zero",
			inputUserID:   0,
			inputID:       1,
			expectedError: domain.ErrInvalidUserID,
		},
		{
			name:          "invalid user ID - negative",
			inputUserID:   -1,
			inputID:       1,
			expectedError: domain.ErrInvalidUserID,
		},
		{
			name:          "invalid asset ID - zero",
			inputUserID:   2,
			inputID:       0,
			expectedError: domain.ErrInvalidFavoriteID,
		},
		{
			name:          "invalid asset ID - negative",
			inputUserID:   2,
			inputID:       -1,
			expectedError: domain.ErrInvalidFavoriteID,
		},
		{
			name:              "database connection error",
			inputUserID:       2,
			inputID:           1,
			expectedError:     domain.ErrInternalServerError,
			useMockRepository: true,
		},
		{
			name:          "favorite not found",
			inputUserID:   2,
			inputID:       1,
			expectedError: domain.ErrFavoriteNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.useMockRepository {
				err = mockRepo.Remove(context.Background(), tt.inputID, tt.inputUserID)

			} else {
				err = favoriteRepo.Remove(context.Background(), tt.inputID, tt.inputUserID)

			}

			// Assert results
			if tt.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

		})
	}
}

var (
	favoriteRepo *internalpostgres.FavoriteRepository
	mockRepo     *internalpostgres.FavoriteRepository
)

func TestMain(m *testing.M) {
	// Setup repositories BEFORE running tests
	t := &testing.T{}
	favoriteRepo = setupRepository(t)

	// Setup mock repository
	sqlDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	mockRepo, err = internalpostgres.NewFavoriteRepository(bun.NewDB(sqlDB, pgdialect.New()))
	require.NoError(t, err)

	// Now run the tests
	exitCode := m.Run()

	// Exit with the appropriate code
	os.Exit(exitCode)
}

func setupPostgres(ctx context.Context) (testcontainers.Container, string, error) {
	// PostgreSQL container configuration
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort("5432/tcp"),
		),
	}

	// Start container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to create postgres container: %w", err)
	}

	// Get connection details
	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return container, "", fmt.Errorf("failed to get mapped port: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return container, "", fmt.Errorf("failed to get host: %w", err)
	}

	// Connection URI
	connStr := fmt.Sprintf("postgresql://testuser:testpass@%s:%s/testdb?client_encoding=UTF8&sslmode=disable", host, mappedPort.Port())

	// Add a small delay to ensure container is fully ready
	time.Sleep(2 * time.Second)

	return container, connStr, nil
}

func setupRepository(t *testing.T) *internalpostgres.FavoriteRepository {
	ctx := context.Background()

	// Use the new setupPostgres function
	container, _, err := setupPostgres(ctx)
	require.NoError(t, err)

	// Parse connection string to get individual components for pgdriver
	// For simplicity, we'll extract the components from the connection string
	// But you might want to use a proper URL parser in production
	mappedPort, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)

	// Create pgdriver connection
	sqlDB := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(fmt.Sprintf("%s:%s", host, mappedPort.Port())),
		pgdriver.WithTLSConfig(nil),
		pgdriver.WithUser("testuser"),
		pgdriver.WithPassword("testpass"),
		pgdriver.WithDatabase("testdb"),
		pgdriver.WithTimeout(5*time.Second),
		pgdriver.WithDialTimeout(5*time.Second),
		pgdriver.WithReadTimeout(5*time.Second),
		pgdriver.WithWriteTimeout(5*time.Second),
	))

	db := bun.NewDB(sqlDB, pgdialect.New())

	// Wait for connection to be ready
	err = db.PingContext(ctx)
	require.NoError(t, err)

	err = initSchema(db.DB)
	if err != nil {
		require.NoError(t, err, "Failed to initialize database schema")
	}

	repo, err := internalpostgres.NewFavoriteRepository(db)
	require.NoError(t, err)

	t.Cleanup(func() {
		// Ensure cleanup happens even if test fails
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := container.Terminate(ctx); err != nil {
			t.Logf("Warning: Failed to terminate container: %v", err)
		}
	})

	return repo
}

// initSchema sets up the database schema from SQL file
func initSchema(db *sql.DB) error {
	// Read the schema file
	schemaSQL, err := os.ReadFile("./../../infra/postgres/migrations/01.schema.up.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Execute the schema SQL
	_, err = db.Exec(string(schemaSQL))
	if err != nil {
		return fmt.Errorf("failed to execute schema SQL: %w", err)
	}

	populateSQL, err := os.ReadFile("./../../infra/postgres/migrations/02.populate.up.sql")
	if err != nil {
		return fmt.Errorf("failed to read populate file: %w", err)
	}

	_, err = db.Exec(string(populateSQL))

	return nil
}
