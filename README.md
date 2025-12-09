# Assets Favourites API

A REST API service for managing user favourites of assets (charts, insights, audiences), implemented as part of the GWI Engineering Challenge.

## Overview

This service provides endpoints for users to manage their personal list of favourite assets. Assets can be of different types (charts, insights, audiences) with varying structures, and are created inline when adding favourites.

## Architecture

The service follows a domain-driven design with a layered, hexagonal-style architecture:

- **Domain layer** (`internal/domain`) - Pure business concepts independent of transport and storage
- **API layer** (`internal/api/v1`) - HTTP handlers, DTOs, and middleware
- **Repository pattern** - Storage abstractions allowing easy swapping of implementations (in-memory → database)

The architecture emphasizes separation of concerns, testability, and extensibility. Observability (logging, metrics, health checks) and authentication are treated as first-class concerns. See [04-architecture-and-standards.md](docs/04-architecture-and-standards.md) for detailed architecture and coding standards.

## Documentation

Detailed documentation is available in the `docs/` directory. Instead of full PRD (Product Requirements Document) and Technical Design Document templates, a lighter free format is used to keep things simple and reviewable without the overhead of formal templates.

- [00-challenge.md](docs/00-challenge.md) - Original GWI Engineering Challenge description
- [01-context.md](docs/01-context.md) - Project context, constraints, and non-functional requirements
- [02-domain-and-use-cases.md](docs/02-domain-and-use-cases.md) - Domain model and use cases
- [03-api-shape.md](docs/03-api-shape.md) - API design and specifications
- [04-architecture-and-standards.md](docs/04-architecture-and-standards.md) - Architecture and coding standards

## API Specification

The OpenAPI specification is available at [api/openapi.yaml](api/openapi.yaml).

## Getting Started

### Prerequisites

- Go 1.25.5 or later
- Make (optional, for convenience commands) - [Install Make](https://www.gnu.org/software/make/)

### Building

```bash
go build -o bin/api ./cmd/api
```

This will compile the application to `bin/api`.

### Running

Run the server with default settings (listens on `:8080`):

```bash
go run ./cmd/api
```

Or run the compiled binary:

```bash
./bin/api
```

### Configuration

The service can be configured using environment variables. See `.env.example` for available options. Note that `.env` files are not used; set environment variables directly if you want to override defaults.

**Available environment variables:**

- `HTTP_ADDR` - Server address (default: `:8080`)

Example:

```bash
HTTP_ADDR=:9090 go run ./cmd/api
```

### Testing

Run all tests:

```bash
go test ./... -race
```

Run specific checks:

```bash
go fmt ./...           # Format code
go vet ./...           # Run go vet
go mod tidy            # Tidy dependencies
```

### Makefile

A `Makefile` is provided for convenience and alignment. **It is strongly recommended to run `make` before committing code** to ensure proper formatting, vetting, and that all tests pass:

```bash
make              # Run all checks: tidy, fmt, vet, build, test
make run          # Run the server
make build        # Build the binary
make test         # Run tests
make fmt          # Format code
make vet          # Run go vet
make tidy         # Tidy dependencies
```

Running `make` executes: `go mod tidy`, `gofmt -s -w .`, `go vet ./...`, `go build`, and `go test ./...`

### Health Check

Once the server is running, verify it's healthy:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{ "status": "healthy" }
```

## Status

This project is in the initial scaffold phase. Core functionality is being implemented.
