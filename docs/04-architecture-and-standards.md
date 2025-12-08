# Architecture & Coding Standards

This document describes the intended architecture and basic coding standards for the Assets Favourites API. It is scoped to the take-home challenge but structured as if it were a small, real production service.

## Design principles

The architecture follows domain-driven design principles with a layered, hexagonal-style structure that separates business logic from infrastructure concerns:

- The **domain layer** (`internal/domain`) contains pure business concepts independent of transport and storage.
- The **API layer** (`internal/api/v1`) handles HTTP concerns through handlers, DTOs, and middleware.
- **Repository interfaces** live in the domain layer and are implemented in `internal/repository`, allowing storage to be swapped (in-memory → database) without changing business logic.

Additional design considerations:

- The service is optimized around the idea that users may have a **large number of favourite assets**. Data structures are keyed by `userId` so listing a user’s favourites is proportional only to that user’s data.
- Changing storage (e.g. moving to Postgres or Redis) should not affect handlers or services.
- Observability (logging, metrics, health) and authentication are treated as **first-class concerns**, not afterthoughts.

## High-level architecture

The service follows a layered architecture with clear separation of concerns.

### Folder structure

```text
platform-go-challenge/
├── cmd/
│   └── api/                    # Application entrypoint
│       └── main.go
├── internal/
│   ├── app/                    # Application setup and wiring
│   ├── config/                 # Configuration management
│   ├── api/                    # API layer
│   │   └── v1/                 # API version 1
│   │       ├── handler/        # HTTP handlers
│   │       ├── dto/            # Data Transfer Objects
│   │       └── middleware/     # API-specific middleware (auth, logging, etc.)
│   ├── domain/                 # Domain layer
│   │   ├── favourites/         # Favourites domain package
│   │   │   ├── model.go        # Domain models (Favourite, Asset)
│   │   │   ├── service.go      # Service (business logic)
│   │   │   └── repository.go   # Repository interface
│   ├── repository/             # Repository implementations
│   │   └── memory/             # In-memory implementations
│   ├── infrastructure/         # Infrastructure concerns
│   │   ├── logger/             # Logging infrastructure
│   │   └── metrics/            # Metrics infrastructure
│   └── shared/                 # Shared utilities
│       ├── errors/             # Shared error types
│       └── httpx/              # HTTP utilities
├── api/                        # API specification
│   └── openapi.yaml
├── scripts/                    # Utility scripts
├── test/                       # Integration tests
│   └── favourites_integration_test.go
├── Dockerfile                  # Docker image definition
├── docker-compose.yaml         # Docker Compose configuration
├── Makefile                    # Build and development commands
├── go.mod
└── README.md
```

### Package descriptions

- **cmd/api**
  Application entrypoint: reads configuration, builds the app, and starts the HTTP server.

- **internal/app**
  Application initialization, HTTP server setup, routing, and dependency wiring:

  - `app.go` – constructs repositories, services, handlers.
  - `router.go` – registers routes and middleware.
  - `server.go` – HTTP server configuration and graceful shutdown.

- **internal/config**
  Configuration management (environment variables, defaults). Exposes a `Config` struct and loader (e.g. `config.Load()`).

- **internal/api/v1** – API version 1 layer

  - **handler/**: HTTP handlers that parse/validate input, call services, and map responses.
  - **dto/**: Request/response structs matching `api/openapi.yaml`.
  - **middleware/**: API-specific middleware (authentication, logging, recovery).

- **internal/domain** – Core domain layer organized by domain packages

  - **favourites/**: Favourites domain package containing `Service`, `Repository` interface, and domain models (`Favourite`, `Asset`) – pure business concepts (see `02-domain-and-use-cases.md`).
    Domain packages follow Go naming guidance: package name provides context (e.g. `favourites.Service`, `favourites.Repository`), so exported names can stay short.
  - **asset/**: Reserved for future asset-specific domain logic if the model grows independently.

- **internal/repository** – Repository implementations

  - **memory/**: In-memory implementations (maps + `sync.RWMutex`) for this challenge. These implement the `favourites.Repository` interface.

- **internal/infrastructure** – Infrastructure concerns

  - **logger/**: Structured logging setup (JSON logs).
  - **metrics/**: Prometheus metrics setup (HTTP metrics, domain counters).

- **internal/shared** – Shared utilities

  - **errors/**: Shared error types and error handling utilities.
  - **httpx/**: HTTP utilities (Problem Details/JSON error helpers, write JSON responses).

- **api/** – API contract

  - `openapi.yaml`: Machine-readable API contract for `/v1` endpoints.

- **scripts/** – Utility scripts
  Development helpers (e.g. `run-local.sh`, `test.sh`).

- **test/** – Integration tests
  End-to-end tests that exercise the service through HTTP using the real router.

This separation keeps the codebase testable and extensible (e.g. adding a real DB or a new transport) without breaking existing code.

## Transport models (DTOs)

DTOs in `internal/api/v1/dto` match the JSON contract in `api/openapi.yaml`. Handlers:

1. Decode JSON into request DTOs.
2. Validate inputs (e.g. required fields, enum values).
3. Map DTOs to domain models and call services.
4. Map domain models back to response DTOs.

The domain models are **not** exposed directly over HTTP. This allows the API contract (field names, shapes) to evolve independently from the domain and storage structures.

## User & identity

This service is **not** a user management system; it treats users as external and identifies them via `userId`:

- `userId` is an **opaque identifier** coming from an upstream identity provider (SSO / auth service).
- The favourites domain stores a mapping from `userId` to that user’s favourite assets. No user profile data is stored here; the service only keys data by `userId`.
- Authentication is handled via middleware:

  - A mock `Authenticator` parses a token (e.g. a simple JWT-like structure) and extracts a `user_id` claim.
  - Handlers verify that the `userId` path parameter matches the authenticated `user_id`.

- All domain operations take `userID string` explicitly as an argument.
  The domain layer does not depend on HTTP, headers, or token formats and remains independent from the authentication mechanism.

This keeps **identity and auth** concerns clearly separated while still enforcing per-user access to favourites.

## Dependency injection

Constructor functions are used to inject dependencies:

- Domain services: `favourites.NewService(repo favourites.Repository)`
- Handlers: `handler.NewFavouriteHandler(svc favourites.Service, logger, metrics, authenticator, ...)`

Key points:

- Services depend on **interfaces**, not concrete implementations (e.g. `favourites.Repository`), making unit tests straightforward.
- Infrastructure concerns (logger, metrics, authenticator) are injected at the edges (handlers, app wiring).
- Global mutable state is avoided; shared state is passed via constructors.

## Storage options

For the exercise, repositories use in-memory maps behind interfaces:

- Favourites are stored keyed by `userId`, e.g. `map[userId][]Favourite`, so listing a user’s favourites is `O(n)` in that user’s favourites, independent of total users.
- Assets are stored in an in-memory map keyed by asset ID.

In a real deployment, the same interfaces could be backed by:

- A **relational database** (e.g. Postgres) with tables for `assets` and `favourites`, indexes on `(user_id, asset_id)` for fast lookups.
- A **key/value store** (e.g. Redis) storing per-user favourite lists.
- A **document store** (e.g. MongoDB) where each user document embeds their favourites.

The architecture keeps storage concerns isolated in `internal/repository`, so swapping storage backends does not impact handlers or services.

## Error handling

Domain/service layers return Go `error` values (typed errors like `ErrFavouriteNotFound` when helpful). HTTP handlers:

- Map domain errors to HTTP status codes (`404` for not found, `409` for conflict, etc.).
- Return Problem Details error responses (RFC 9457) from helpers in `internal/shared/httpx`.
- Log errors once at the edge (handler or middleware), avoiding duplicate logging deeper in the stack.

This ensures errors are surfaced consistently to clients and diagnosable in logs.

## Logging & metrics

Logging and metrics are implemented via middleware and infrastructure helpers:

- **Structured logging** (JSON format):

  - Request logging middleware logs method, path, status code, and duration.
  - Business events (e.g. favourite added/removed) can log at info level with structured fields (`userId`, `assetId`).

- **Metrics**:

  - `/metrics` endpoint exposes Prometheus metrics:

    - HTTP request count and latency by route, method, and status code.
    - Domain counters (e.g. favourites added/removed).

  - Health endpoint `/health` returns a simple JSON body:

    ```json
    { "status": "healthy" }
    ```

This supports both basic observability for the exercise and a clear path to production-grade monitoring.

## API lifecycle & versioning

The API is versioned and designed with evolution in mind:

- All endpoints are namespaced under `/v1`, with the contract defined in `api/openapi.yaml`.
- Error responses use a consistent Problem Details format so clients can handle failures predictably.
- Authentication is abstracted via middleware and an `Authenticator` interface, providing a clear place to plug in real OAuth2/JWT/SSO when needed.
- The folder structure (`internal/api/v1`) allows a future `/v2` implementation to coexist alongside `/v1` without breaking existing clients.
- Rate limiting and quotas are out of scope for this exercise but can be added as middleware at the router layer without changing handlers or domain logic.

## Testing approach

- **Unit tests**

  - Located in the same folder as the code they test (e.g. `service.go` has `service_test.go`).
  - Services are tested with in-memory or mock repositories.
  - Handlers are tested using `httptest` with the real router.
  - Focus on testing individual components in isolation, including happy path and key error scenarios.

- **Integration tests**

  - Located in the `test/` folder.
  - At least one end-to-end handler test per main endpoint (happy path + one error case).
  - Integration tests verify the interaction between multiple components (handler → service → repository) through the full HTTP request/response cycle.

- Tests are expected to be **fast, deterministic**, and free of external dependencies (no live DB).

## Coding standards

The codebase aims for idiomatic, "boring" Go that any Go engineer can read and maintain.

Core principles:

- **Formatting & tooling**

  - All code is formatted with `gofmt`.
  - `go vet` (and similar linters) should be clean before committing.
  - Imports are organised and unused code is removed.

- **Idiomatic Go style**

  - Follow the official guidance from:

    - Go Code Review Comments
    - Effective Go
    - Google Go Style Guide

  - Use clear, descriptive names and avoid stutter between package and exported names.
    The package name should provide most of the context, so type names usually don’t need to repeat it. For example, for favourites-related logic, a domain package like `favourites` with `type Service` (used as `favourites.Service`) is preferred over a generic `service` package with `FavouritesService`. Similarly, in the `config` package, `type Loader` (used as `config.Loader`) is better than `type ConfigLoader`. The goal is that code reads well at the call site: `favourites.Service`, `config.Loader`, `asset.Repository`, etc. This follows the guidance in Effective Go and Go Code Review Comments that package names provide context and exported names should avoid stuttering (e.g. `bufio.Reader`, not `bufio.BufReader`).
  - Choose pointer vs value receivers according to the Google style guide (mutating methods or large structs → pointer; small, immutable value types → value).

- **Variable names & abbreviations**

  - The length of a name should roughly match the length of its scope (per Effective Go). Short names are fine for short scopes; longer, descriptive names are preferred for wider scopes and exported identifiers.
  - Conventional Go abbreviations are used for well-known concepts:
    - `ctx` for `context.Context`
    - `w` / `r` for `http.ResponseWriter` / `*http.Request` in handlers
    - `err` for `error`
  - For local variables that clearly represent a service, `svc` is acceptable (e.g. `svc := favourites.NewService(repo)` in `main` or `app.go`). For struct fields or exported names, full words are preferred, e.g. `service *favourites.Service`.
  - Project-specific or unclear abbreviations are avoided; descriptive names are preferred where there is any ambiguity.

- **Errors & control flow**

  - Functions return `error`. `panic` is not used in normal control flow.
  - Check errors explicitly and return early rather than deep nesting.
  - Use well-named sentinel/typed errors (e.g. `ErrFavouriteNotFound`) and map them to HTTP status codes at the edge.

- **APIs & layering**

  - Handlers stay thin: parse/validate input, call services, map domain errors to HTTP.
  - Services implement the use cases and depend on small interfaces defined close to the consumer.
  - Repositories hide storage details (in-memory now, DB later) behind interfaces.
  - Common Go API design patterns are followed (clear separation of handler → service → repository, constructor-based DI, context-aware services).

- **Concurrency & context**

  - Shared in-memory state is protected with `sync.RWMutex`.
  - `context.Context` is passed from handlers down through services and repositories for cancellation/timeouts; contexts are never stored in structs or global variables.

- **Tests**

  - Tests live next to the code they test (`*_test.go` in the same package).
  - Use table-driven tests where appropriate.
  - HTTP tests use `httptest` and the real router; services are tested with in-memory or mocked repositories.

### References

- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- [Effective Go](https://go.dev/doc/effective_go)
- [Google Go Style Guide](https://google.github.io/styleguide/go/)
- [5 API design patterns in Go that solve your biggest problems (2025)](https://cristiancurteanu.com/5-api-design-patterns-in-go-that-solve-your-biggest-problems-2025/)
