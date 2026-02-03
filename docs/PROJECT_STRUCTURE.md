# Project Structure

This document provides a map of the codebase to help new developers navigate the project.

## Directory Layout

```text
.
├── api/                    # OpenAPI/Swagger definitions
│   └── openapi.yaml
├── cmd/                    # Application Entry Points
│   └── server/             # The main HTTP server application
│       └── main.go         # Wires up dependencies and starts the server
├── deployments/            # Infrastructure configuration (k8s, docker-compose)
├── docs/                   # Architecture and Design documentation
├── internal/               # Private application code (Library/App code)
│   ├── adapter/            # Infrastructure implementations (Adapters)
│   │   ├── api/            # HTTP/REST Layer (Handlers, DTOs, Router)
│   │   ├── cache/          # Cache implementations (Redis)
│   │   └── storage/        # Database implementations (PostgreSQL/pgx)
│   ├── config/             # Configuration loading and validation
│   ├── core/               # Pure Domain Logic (The "Hexagon")
│   │   ├── domain/         # Domain Entities (Asset, Auth/User)
│   │   ├── ports/          # Interfaces (Repositories, Services)
│   │   └── service/        # Business Logic (Favorites, Auth)
│   └── observability/      # Metrics (Prometheus) and Tracing (OpenTelemetry)
├── tests/                  # External Tests
│   ├── integration/        # Testcontainers-based integration tests
│   └── k6/                 # Load testing scripts
├── Dockerfile              # Distroless Docker build
├── go.mod                  # Go Module definition
├── Makefile                # Build and dev shortcuts
└── README.md               # Quickstart and Overview
```

## Key Components

### `internal/core`

This is the heart of the application. It contains no infra-specific imports.

* **`domain`**: Defines what an `Asset` is, validations, and the Enum types.
* **`ports`**: Defines the contracts. For example, `FavoriteRepository` is an interface here. The implementation lives in `internal/adapter/storage`.

### `internal/adapter`

Bridging the gap between the Core and the outside world.

* **`storage/postgres`**: Implements `ports.FavoriteRepository`. Uses `pgx` for connection pooling.
* **`api/rest`**: Implements the HTTP handler. Converts HTTP requests to Service calls and Domain objects to JSON responses.

### `tests/integration`

We prioritize integration tests over extensive unit mocking. These tests spin up real Postgres and Redis containers using `testcontainers-go` to verify the system works end-to-end.
