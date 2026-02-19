# Favourites API – Engineering Challenge Submission

## Introduction

This service implements a REST API for managing user-favourited assets. Beyond satisfying the functional requirements, the implementation prioritises production readiness: clean architectural boundaries, efficient pagination for unbounded datasets, a replaceable storage layer, and comprehensive test strategy.

The areas I focused on beyond the functional requirements are:

- A clean layered architecture that separates HTTP, business logic, and storage concerns
- Cursor-based pagination to handle unbounded favourites lists efficiently
- A cache-aside pattern with explicit invalidation to keep reads fast
- A storage design that is in-memory for simplicity but clearly structured for a real database backend
- Testing at multiple levels: unit, integration, and shell-based end-to-end scenarios
- Docker support, structured logging, graceful shutdown, and Swagger documentation



## Assumptions

To keep the scope realistic and focused on backend engineering concerns, the following assumptions were made:

- Users and assets originate from upstream systems, therefore:
  - `user_id` and `asset_id` are external identifiers — this service does not own or validate them
  - This service does not generate or manage user identity
  - Authentication and authorisation would be handled by a dedicated upstream auth service in a real enterprise architecture
- The number of favourites per user is unbounded, therefore performance must be considered — this directly drives the pagination and caching decisions described below


## Architecture

The application follows a layered architecture with clear boundaries:

```
Handler → Service → Repository
```

- **Handlers** deal exclusively with HTTP: reading request parameters, validating input, and writing responses.
- **Services** contain the business logic and coordinate reads and writes.
- **Repositories** abstract the storage backend behind an interface.

The `Repository` interface means the in-memory store used here can be replaced with a PostgreSQL implementation without touching a single line of handler or service code. The same applies to the cache layer.



## Domain Model

The challenge specifies three asset types — `Chart`, `Insight`, and `Audience` — which share common attributes but have completely different structures. I modelled this using struct embedding and a shared interface, which is the idiomatic Go approach:

```go
type BaseAsset struct {
    AssetID      string
    Type         AssetType
    Description  string
    FavouritedAt time.Time
}
```

Each concrete type (`Chart`, `Insight`, `Audience`) embeds `BaseAsset` and adds its own fields. A shared `Asset` interface allows the service layer to handle all types uniformly. Adding a new asset type in the future requires no changes to existing logic — only a new struct and a registration in the unmarshalling layer.


## API

**Base path:** `/api/v1`

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/users/{userID}/favorites` | List user favourites (paginated) |
| `POST` | `/users/{userID}/favorites` | Add an asset to favourites |
| `GET` | `/users/{userID}/favorites/{assetID}` | Retrieve a specific favourite |
| `PUT` | `/users/{userID}/favorites/{assetID}/description` | Update description |
| `DELETE` | `/users/{userID}/favorites/{assetID}` | Remove a favourite |
| `GET` | `/health` | Service and cache health |

Full interactive documentation is available at `/swagger/index.html` when the service is running.



## Pagination

The challenge notes that users have no limit on favourites and that the service must provide reasonable response times. This rules out offset-based pagination, which degrades linearly with dataset size and produces inconsistent results when records are inserted or deleted between requests.

**keyset (cursor) pagination** was implemented instead.

**Ordering:** `favourited_at DESC`, with `asset_id DESC` as a tie-breaker for entries added at the same instant.

The cursor encodes both fields and is passed opaquely between client and server. On each request, the server decodes the cursor and uses it as a range condition rather than a row offset:

```sql
SELECT * FROM favourites
WHERE user_id = $1
  AND (favourited_at, asset_id) < ($cursor_time, $cursor_asset_id)
ORDER BY favourited_at DESC, asset_id DESC
LIMIT $limit;
```

This makes every page fetch a direct index seek — O(1) regardless of how far into the list the client has paginated. Results are stable under concurrent writes.


## Caching

A **cache-aside** pattern was implemented at the user level.

- **On read:** the service checks the cache first. On a hit, the result is returned immediately without touching the store. On a miss, the store is queried and the result is written to the cache.
- **On write:** any mutation (add, update, delete) explicitly invalidates the affected user's cache entry, ensuring subsequent reads always reflect current state.

The current implementation uses an in-memory TTL cache. The cache layer sits behind an interface, so replacing it with Redis in a production deployment — where the cache needs to be shared across multiple API instances — requires no changes to the service logic.

Cache statistics (hit rate, miss count, current entry count) are exposed via `/health`.


# Storage Design (Current + Production-Ready Options)

## Current Implementation

- In-memory store (`map` + `sync.RWMutex`)
- Repository abstraction so the storage backend can be swapped without changing handlers/services
- Optimised for simplicity and testability (as the challenge does not require a deployed DB)


## Production-Ready Options

The data model is essentially:

- Many users
- Unlimited favourites per user
- Each favourite references an upstream `asset_id` and stores:
  - `description` (mutable)
  - `favourited_at` (used for sorting/pagination)
  - asset payload (type-specific fields)

The key design choice is: **how to persist polymorphic asset payloads.**


## Option A — PostgreSQL (Relational)

A relational approach works well because the core operations are:

- **Uniqueness:** "a user cannot favourite the same asset twice"
- **Fast lookup** by `(user_id, asset_id)`
- **Efficient ordered listing** by `(favourited_at, asset_id)`

**Suggested tables:**

```sql
CREATE TABLE favourites (
  user_id        TEXT NOT NULL,
  asset_id       TEXT NOT NULL,
  asset_type     TEXT NOT NULL,
  description    TEXT,
  favourited_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, asset_id)
);

-- Supports keyset pagination and "latest favourites"
CREATE INDEX idx_favourites_user_order
  ON favourites (user_id, favourited_at DESC, asset_id DESC);
```

**Why this is strong:**

- Composite PK guarantees no duplicates
- Index supports keyset pagination efficiently
- Strong consistency and mature observability/tooling


## Option B — PostgreSQL + JSONB (Polymorphic Payload Storage)

Since assets can have different structures (Chart/Insight/Audience), a practical approach is:

- Keep typed relational columns for fields needed for queries/pagination/uniqueness
- Store the "variable" part in `asset_payload JSONB`

**Example:**

```sql
CREATE TABLE favourites (
  user_id        TEXT NOT NULL,
  asset_id       TEXT NOT NULL,
  asset_type     TEXT NOT NULL,
  description    TEXT,
  favourited_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  asset_payload  JSONB NOT NULL,
  PRIMARY KEY (user_id, asset_id)
);
```

**When does JSONB be a good fit:**

- Asset schema evolves frequently
- No need for complex queries inside the payload
- The asset is mainly retrieved and displayd as a blob

**Trade-offs:**

- Harder to enforce strict schema constraints
- Querying deeply inside JSONB is possible, but often slower/less practical than relational columns


## Option C — Separate tables per asset type (Fully Relational)

Another option is:

- `charts`, `insights`, `audiences` tables (source of truth)
- `favourites` table references them

However, in this challenge's assumptions, assets already come from upstream systems, so this service would likely store only references + minimal data, not own the asset tables.

> This option makes sense if this service becomes the owner of asset data.


## Option D — Document store (MongoDB)

MongoDB can model favourites as embedded documents:

```json
{
  "user_id": "user-123",
  "favorites": [
    { "asset_id": "...", "type": "...", "favourited_at": "...", "payload": {} }
  ]
}
```

**Pros:**

- Natural fit for polymorphic payloads
- Easy to store full asset blobs without JSON schema migrations

**Cons / Risks at scale:**

- "Unlimited favourites per user" can lead to very large user documents
- Updates to a single favourite may rewrite large documents
- Pagination over embedded arrays can become awkward and less efficient
- Might require modelling favourites as a separate collection anyway

A more scalable MongoDB approach would be a **separate collection**: favourites documents keyed by `(user_id, asset_id)` with an index on `(user_id, favourited_at, asset_id)`.


## Summary: Recommended Path

For production, the most pragmatic and scalable default is:

- **PostgreSQL** for strong constraints + indexing + keyset pagination
- **Optional JSONB column** for the polymorphic asset payload
- **Redis** for distributed caching when horizontally scaling

This keeps:

- fast reads
- safe writes (uniqueness constraints)
- simple scaling path (cache + DB indexes)
  
#### Production stack

```
Client → API instances → Redis → PostgreSQL
```

Redis replaces the in-memory cache in a multi-instance deployment. The cache interface ensures this requires no changes to the service layer.


## Testing

I implemented testing at three levels:

**Unit tests** cover individual components in isolation: asset JSON marshalling and unmarshalling, cursor encoding and parsing, service logic for all CRUD operations, cache hit/miss behaviour and invalidation, and in-memory store correctness.

**Integration tests** verify the full HTTP lifecycle: correct status codes and response shapes, pagination correctness across multiple pages, cursor validity and consistency, and the absence of duplicate records across page boundaries.

**Shell scripts** simulate realistic usage against a running instance: bulk seeding with hundreds of assets, a full paginated walk through the result set, cache invalidation validation, and a golden path end-to-end scenario.

---

## Observability & Operational Readiness

- **Structured logging** with `zap`, configured per environment — human-readable output in development, structured JSON in production.
- **Health endpoint** at `/health` exposes service status and cache statistics.
- **Graceful shutdown** — the server listens for OS signals and waits for in-flight requests to complete before exiting, preventing dropped connections during restarts or deployments.


# Running the Project

## Option 1 — Local (Go)

```bash
go mod tidy
go run ./cmd/api
```

The server starts at:

```
http://localhost:8080
```

Health check:

```bash
curl http://localhost:8080/api/v1/health | jq
```

Swagger UI:

```
http://localhost:8080/swagger/index.html
```


## Option 2 — Using Make (Recommended)

Common development tasks are wrapped in a `Makefile` for convenience and reproducibility:

```bash
make run              # Run API locally
make test             # Run tests with coverage
make swagger          # Regenerate Swagger docs
make build            # Build local binary
make compose-up       # Docker compose (build + run)
make compose-down     # Stop docker compose
```

This provides a clean workflow for development and review.

## Option 3 — Docker

```bash
docker compose up --build
```

Stop the service:

```bash
docker compose down
```

The container exposes:

```
http://localhost:8080
```


## Configuration

All configuration is provided via environment variables (optionally through a `.env` file). Sensible defaults are used when variables are not defined.

| Variable | Default | Description |
|---|---|---|
| `ADDR` | `:8080` | HTTP bind address |
| `ENV` | `development` | Environment (controls logging mode) |
| `API_URL` | `http://localhost:8080` | Base URL used by scripts |
| `CACHE_TTL` | `5m` | Cache TTL in duration format |

Example:

```bash
cp .env.example .env
```

# Seeding & Runtime Validation

To properly exercise pagination and caching behaviour, realistic data volume is required.

All validation scripts are located under `scripts/`.

Before running them for the first time:

```bash
chmod +x scripts/*.sh
```

You may execute them directly or via `make`.


## Script Requirements

The scripts depend on:

* `curl` — for HTTP requests
* `jq` — for JSON parsing and validation

If missing, execution will stop automatically (see `scripts/00_env.sh`).

### macOS

```bash
brew install curl jq
```

### Ubuntu / Debian

```bash
sudo apt update
sudo apt install curl jq
```

### Windows

Recommended: Use WSL (Windows Subsystem for Linux)

```bash
wsl --install
```

Then inside Ubuntu:

```bash
sudo apt install curl jq
```

## Scripts

### 1. Seed Many Assets

Creates a large number of favourites for a user. Useful for scale and pagination testing.

Direct:

```bash
TOTAL=500 USER_ID=user-scale scripts/01_seed_many_assets.sh
```

Via Make:

```bash
TOTAL=500 USER_ID=user-scale make seed
```

| Variable | Description |
|---|---|
| `TOTAL` | Number of favourites to generate |
| `USER_ID` | User identifier |

### 2. Pagination Walk (Cursor Validation)

Simulates a client scrolling through all pages using cursor-based pagination.

The script:
* Calls `GET /favorites`
* Follows `next_cursor`
* Continues until `has_more=false`
* Verifies there are no duplicates across pages

Direct:

```bash
USER_ID=user-scale scripts/02_paginate_walk.sh
```

Via Make:

```bash
USER_ID=user-scale LIMIT=50 EXPECTED=500 make paginate
```

This validates:
* Keyset pagination correctness
* Cursor stability
* Ordering guarantees
* Complete dataset traversal


### 3. Cache Behaviour Validation

Demonstrates the cache lifecycle:

1. Initial read → miss
2. Second read → hit
3. Write → invalidation
4. Next read → miss again

Direct:

```bash
USER_ID=user-scale scripts/03_cache_check.sh
```

Via Make:

```bash
USER_ID=user-scale make cache-check
```

Validates:
* Read-through caching
* Invalidation on writes
* TTL behaviour
* Per-user cache isolation


### 4. Golden Path (End-to-End API Contract)

Runs a complete happy-path scenario:
* `POST` favourite
* `GET` favourites
* `GET` specific favourite
* `PUT` description
* `DELETE` favourite

Direct:

```bash
USER_ID=user-demo scripts/04_golden_path.sh
```

Via Make:

```bash
USER_ID=user-demo make golden
```

This ensures the full API contract behaves correctly from a client perspective.


## Running Tests

Via Make:

```bash
make test
```

Or directly:

```bash
go test ./... -cover -count=1
```