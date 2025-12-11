# Assets Favourites API

A small Go service for managing user favourites of assets (charts, insights, audiences).  
Think "favourites for a data platform": simple to run locally, easy to extend, and ready to plug into a bigger system.

> **📋 For Reviewers:** See [REVIEW_SCENARIO.md](./REVIEW_SCENARIO.md) for a walk-through with copy-paste curl commands. This is the recommended starting point, after reading this document.

---

## 1. Overview

Users on a data platform can browse a large catalogue of assets. This service lets them:

- Mark assets as **favourites** for quick access
- **List** all favourites for a given user
- **Add** new favourites by defining the asset inline
- **Edit** an asset’s **description**
- **Remove** a favourite

Supported asset types:

- **Chart** – title, axis titles, and data points
- **Insight** – a short text insight
- **Audience** – gender, birth country (2-letter code), age groups, hours on social media, purchases last month

Key characteristics:

- Domain-first design (`internal/domain/favourites`)
- Clear separation between HTTP layer, domain, and storage
- In-memory storage with concurrency safety (no external DB required)
- Prometheus metrics and a pre-built Grafana dashboard
- Tests across domain, repository, handlers, and basic integration

### Design highlights

Key decisions:

- **Domain layer independence**: No HTTP or storage dependencies. The domain is a plain Go package usable from any transport.
- **Repository interface**: Storage abstracted behind an interface, easy to swap implementations without changing business logic.
- **Authentication abstraction**: Pluggable `Authenticator` interface. Current mock reads `X-User-ID` headers; swapping to JWT/OAuth2 only requires changing middleware.
- **Pagination**: Ensures predictable response times even when users have thousands of favourites.
- **Observability**: Prometheus metrics and structured logging built in for production readiness.
- **Concurrency safety**: Uses `sync.RWMutex` with `-race` tests to catch data races.

---

## 2. Design documentation

Design documentation:

- **[Context & Constraints](./docs/01-context.md)** – Project context, constraints, and design assumptions
- **[Domain & Use Cases](./docs/02-domain-and-use-cases.md)** – Domain model, core use cases, and business logic
- **[API Design](./docs/03-api-shape.md)** – API endpoints, request/response formats, and standards (HTTP, JSON, error handling)
- **[Architecture & Standards](./docs/04-architecture-and-standards.md)** – System architecture, coding standards, and design principles

These documents explain the design decisions and codebase structure.

---

## 3. Getting started

### 3.1. Prerequisites

- Go **1.24+**
- Docker & Docker Compose (optional, for Prometheus + Grafana)
- `make` (optional but recommended)

### 3.2. Run the API

From the project root:

```bash
# Using Go directly
HTTP_ADDR=:8080 go run ./cmd/api

# Or using Make
make run

# Or using Docker Compose (includes Prometheus + Grafana)
docker compose up --build
```

For complete examples and a walk-through, see [REVIEW_SCENARIO.md](./REVIEW_SCENARIO.md).

### 3.3. Configuration

Configuration is intentionally minimal and optional:

- `HTTP_ADDR` – HTTP listen address (default: `:8080`)

Example:

```bash
HTTP_ADDR=:9090 go run ./cmd/api
```

An `.env.example` file documents available variables. Environment loading is handled by your shell, docker-compose, or CI and it doesn't use .env files.

### 3.4. Makefile

A `Makefile` is provided to keep the development workflow consistent:

```bash
make         # gofmt, vet, test, build
make run     # run the API locally
make test    # run tests (including race detector)
make fmt     # gofmt
make vet     # go vet
```

In practice, the expectation is:

> Run `make` before pushing or committing changes so that formatting, vetting, and tests are always in a good state.

That’s the bar this project assumes for day-to-day development.

---

## 4. API overview

The API endpoints are versioned under `/v1`, except for system endpoints (`/health`, `/metrics`) which are unversioned. Authentication is simulated via an `X-User-ID` header and enforced by middleware; the domain layer only deals with opaque `userId` strings.

### 4.1. Endpoints

**System (unversioned)**

- `GET /health` – Health check endpoint
- `GET /metrics` – Prometheus metrics endpoint

**Favourites (v1)**

- `GET /v1/users/{userId}/favourites` – List user's favourites
  - Query params: `limit` (default: 50, max: 100), `offset` (default: 0), `type` (optional filter)
- `POST /v1/users/{userId}/favourites` – Add a new favourite (creates asset inline)
- `PATCH /v1/users/{userId}/favourites/{assetId}` – Update favourite description
- `DELETE /v1/users/{userId}/favourites/{assetId}` – Remove favourite (idempotent)

For complete API examples with request/response details, see [REVIEW_SCENARIO.md](./REVIEW_SCENARIO.md) or the [OpenAPI spec](./api/openapi.yaml).

### 4.2. Authentication model

This service does **not** manage users. It assumes:

- An upstream identity provider issues a user ID.
- Requests include `X-User-ID: <userId>` header.
- The `{userId}` path parameter must match the authenticated user ID.

Auth is handled via a pluggable `Authenticator` interface in the HTTP layer. The domain works with plain `userID string` and doesn't know about headers or tokens. Swapping to JWT, OAuth2, or SSO later means changing the
`Authenticator` implementation, not the domain.

### 4.3. OpenAPI spec

OpenAPI spec: `api/openapi.yaml` (import into Postman or Swagger UI).

---

## 5. Domain, storage & data model

### 5.1. Package layout

```text
cmd/api/                # Entrypoint (main.go)
internal/
  app/                  # Wiring, router, server startup
  config/               # Configuration loading
  api/
    handler/            # Non-versioned handlers (health)
    v1/
      handler/          # Favourites HTTP handlers
      dto/              # Request/response DTOs
      middleware/       # Auth & other API middleware
  domain/
    favourites/         # Domain models, service, repository interface
  repository/
    memory/             # In-memory favourites repository
  infrastructure/
    logger/             # Logging abstraction
    metrics/            # Prometheus metrics + middleware
  shared/
    errors/             # Shared error helpers
    httpx/              # JSON + Problem Details utilities
api/
  openapi.yaml          # API spec
docs/
  00-challenge.md
  01-context.md
  02-domain-and-use-cases.md
  03-api-shape.md
  04-architecture-and-standards.md
deploy/
  prometheus/           # Prometheus config
  grafana/              # Grafana provisioning & dashboards
```

### 5.2. Favourites domain

The domain (`internal/domain/favourites`) defines:

- `Asset` – a tagged union with:
  - `type`: `chart | insight | audience`
  - `description`: shared
  - `chart`, `insight`, `audience`: optional payloads, exactly one set depending on `type`
- `FavouriteEntry` – one favourite for one user:
  - `userId`
  - `assetId` (UUIDv4)
  - `createdAt`
  - `asset` (the `Asset` struct above)
- `Repository` – abstract interface for storage operations:
  - create favourite
  - list favourites by user
  - remove favourite
  - update description
- `Service` – orchestrates validation and repository calls:
  - validates that the asset shape matches its `type`
  - ensures operations are scoped to a user’s own favourites
  - can be extended with domain-level metrics or auditing without touching handlers
  - The listing endpoint is **paginated** using `limit` and `offset` query parameters. The handler enforces a default limit of 50 and a maximum limit of 100 to keep response times stable even for users with very large favourites lists.

The domain is a plain Go package usable from any transport, not just HTTP.

### 5.3. In-memory storage (current implementation)

The challenge explicitly allows **no real database**, but expects a clear story about storage and data representation. The current implementation uses:

- In-memory maps keyed by `userId` and `assetId`
- UUIDv4 for `assetId`
- `sync.RWMutex` to guard access and avoid race conditions when multiple goroutines add/remove/list favourites

The main design focus is:

- **Per-user scalability**: listing favourites is proportional to that user’s data only, not the global number of favourites.
- **Safety under concurrency**: tests run with `-race` to validate the locking strategy.

For local development and tests this is simple, fast, and has no external dependencies.

### 5.4. Production storage options

The `Repository` interface is storage-agnostic. The same domain model can map to:

- **Document store** (like MongoDB or Firestore): appears optimal for this service because:
  - Access patterns are per-user (no cross-user queries needed)
  - One document per user with embedded favourites array matches the data model naturally
  - Simple queries and good horizontal scaling via sharding by `userId`
- **SQL databases** (like Postgres or Google Cloud SQL): could be alternatives, especially if the platform already uses SQL, requires cross-user analytics, or needs strong ACID guarantees for related operations
- **Caching layer** (Redis): can serve as a read-through cache for high-traffic scenarios

The choice depends on platform context and requirements. All implementations can be swapped without changing handlers or service logic.

---

## 6. Observability

### 6.1. Metrics

Metrics are exposed at the unversioned endpoint:

- `GET /metrics`

**Note:** The `/metrics` endpoint is intentionally unauthenticated for this exercise to allow easy review and Prometheus scraping in docker-compose. In production, this should be protected with internal auth or special admin gateway.

The service uses a dedicated Prometheus registry (`internal/infrastructure/metrics`) and publishes:

- `http_requests_total{method, path, status}`  
  Total HTTP requests, labelled by:

  - `method` – HTTP method
  - `path` – chi route pattern (e.g. `/v1/users/{userId}/favourites`)
  - `status` – HTTP status code

- `http_request_duration_seconds_bucket{method, path, status, le}`  
  Histogram buckets for request durations

The metrics middleware:

- Wraps `http.ResponseWriter` to capture the status code
- Uses `chi.RouteContext` to get the route pattern (keeps label cardinality under control)
- Records duration for every request

### 6.2. Prometheus + Grafana stack

You can start a full local observability stack with one command:

```bash
docker compose up --build
```

This starts:

- API on `http://localhost:8080`
- Prometheus on `http://localhost:9090`
- Grafana on `http://localhost:3000`

Grafana includes an **API Overview** dashboard with:

- Requests per second by endpoint & status
- p95 latency per endpoint
- Table of endpoint + status + req/s (last 5 minutes)
- Global error rate (4xx + 5xx)

To see the graphs move, send traffic using the examples in [REVIEW_SCENARIO.md](./REVIEW_SCENARIO.md).

---

## 7. Testing & quality

Tests cover:

- Domain and repository behaviour
- Service logic (validation, ownership checks)
- HTTP handlers (using `httptest` and the real router)
- Metrics middleware and `/metrics` endpoint
- Basic integration flows (add + list via HTTP)
- Concurrency of the in-memory repository (run under `-race`)

Run everything:

```bash
go test ./... -race
```

The code follows idiomatic Go practices:

- Effective Go
- Go Code Review Comments
- Google Go Style Guide

Aim: readable, understandable, and easy to change.

---

## 8. Possible next steps

This project is intentionally small, but it's structured to grow. Possible next steps can include:

- Replace the in-memory repository with a production storage backend.
- Plug in real authentication (JWT/OAuth2/SSO) behind the `Authenticator` interface
- Add rate limiting and quotas for favourites operations
- Add OpenTelemetry tracing and richer structured logging
- Set up GitHub Actions workflow for continuous integration (build, test, lint)

This is a working end-to-end favourites service you can run, test, and inspect locally.
