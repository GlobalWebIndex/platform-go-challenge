# System Architecture

This document details the architectural design decisions, patterns, and structure of the Favorites Service.

## Architectural Pattern: Hexagonal Architecture (Ports & Adapters)

The application follows the **Hexagonal Architecture** pattern (also known as Ports & Adapters) to ensure separation of concerns and testability.

* **Core (`internal/core`)**: Contains the pure business logic and domain entities. It has **zero dependencies** on frameworks, databases, or HTTP.
* **Ports (`internal/core/ports`)**: Interfaces that define how data enters and leaves the core. This is the boundary of the hexagon.
* **Adapters (`internal/adapter`)**: Implementations of the ports (e.g., `pgx` for Postgres, `net/http` for REST).

This structure allows us to swap the database, cache, or transport layer without touching the business logic.

```mermaid
block-beta
columns 3
  block:adapters_left
    columns 1
    HTTP["HTTP / REST"]
    CLI["CLI / Main"]
  end

  block:core
    columns 1
    block:ports
      columns 3
      InPort(("Input Ports")) space OutPort(("Output Ports"))
    end
    block:domain
      columns 1
      Service["Service Layer (Business Logic)"]
      Model["Domain Models (Asset, User)"]
    end
  end

  block:adapters_right
    columns 1
    Postgres["PostgreSQL (Persistence)"]
    Redis["Redis (Cache)"]
    External["Enrichment Service"]
  end

  HTTP --> InPort
  Service --> OutPort
  OutPort --> Postgres
  OutPort --> Redis
  style core fill:#f9f,stroke:#333,stroke-width:2px
```

## Zero-Allocation Streaming Pattern

To satisfy the requirement of "unlimited favorites" without crashing memory (OOM), we utilize **Go 1.25 Iterators (`iter.Seq2`)**.

### Data Flow

```mermaid
sequenceDiagram
    participant DB as PostgreSQL (Cursor)
    participant Repo as Repository (iter.Seq2)
    participant Svc as Service Layer
    participant HTTP as HTTP Handler
    participant Client as HTTP Client

    Note over DB, Client: O(1) Memory Usage Flow

    HTTP->>Svc: FindAll(limit=N)
    Svc->>Repo: FindAll(limit=N)
    Repo->>DB: Query (Rows)
    DB-->>Repo: Row Cursor
    Repo-->>Svc: return iter.Seq2[Asset]
    Svc-->>HTTP: return iter.Seq2[Asset]

    loop Pull-based Streaming
        HTTP->>Svc: next()
        Svc->>Repo: next()
        Repo->>DB: Scan Row
        DB-->>Repo: Raw Data
        Repo-->>Svc: Domain Object
        Svc-->>HTTP: Domain Object
        HTTP->>Client: Write JSON Chunk
    end
```

### Explanation

1. **Cursor-Based Query**: The Database Adapter acts as a cursor. It does not load all rows into a slice.
2. **Iterator Yielding**: The Repository layer converts SQL rows into domain objects one by one, yielding them to the Service layer via `iter.Seq2[domain.Asset, error]`.
3. **HTTP Streaming**: The HTTP Handler loops over this iterator. For each item, it marshals it to JSON and writes it immediately to the `http.ResponseWriter`.

**Result:** Memory usage remains flat (approx. size of 1 Asset + Buffer) regardless of whether the user has 100 or 1,000,000 favorites.

## Container Diagram (C4)

The system is composed of a core Go REST API, backed by PostgreSQL for persistence and Redis for high-speed caching and sorting. It interacts with an external Enrichment Service to hydrate asset data.

```mermaid
C4Context
    title C4 Container Diagram - Favorites System

    Person(user, "API Client", "Web Frontend or External System")
    
    System_Boundary(favorites_boundary, "Favorites System") {
        Container(api, "Favorites API", "Go 1.25, net/http", "Handles REST requests, Authentication (JWT), validation, and orchestration.")
        ContainerDb(postgres, "PostgreSQL", "v16", "Primary persistent storage. Stores Users, JSONB assets + GIN Indexes.")
        ContainerDb(redis, "Redis", "v7", "Write-through cache. Stores ZSETs for pagination and Hashes for asset data.")
    }

    System_Ext(enricher, "Enrichment Service", "External HTTP API", "Provides additional metadata for assets.")

    Rel(user, api, "Uses", "HTTPS/JSON")
    Rel(api, postgres, "Reads/Writes", "TCP/5432 (pgx)")
    Rel(api, redis, "Reads/Writes", "TCP/6379")
    Rel(api, enricher, "Enriches Assets", "HTTPS")
```

## Request Flow: Get Favorites (Async Write-Through)

This sequence diagram illustrates the **Async Write-Through** pattern used to decouple write performance from read performance.

```mermaid
sequenceDiagram
    autonumber
    participant Client
    participant API
    participant Redis
    participant Postgres
    participant Enricher as Enrichment Svc

    Note over Client, API: GET /favorites?limit=100&offset=0

    Client->>API: Request Favorites List (Bearer Token)
    
    rect rgb(255, 230, 230)
    Note right of API: 0. Auth Check
    API->>API: Validate JWT Signature
    end
    
    rect rgb(240, 248, 255)
    Note right of API: 1. Cache Check
    API->>Redis: ZREVRANGE (Get IDs)
    Redis-->>API: (Cache Miss / Empty)
    end

    rect rgb(255, 250, 240)
    Note right of API: 2. DB Streaming
    API->>Postgres: Query (Cursor/Iterator)
    Postgres-->>API: Returns iter.Seq2 (Stream)
    end

    loop For Each Batch (Streaming Response)
        API->>Enricher: Enrich Asset (Parallel/errgroup)
        Enricher-->>API: Enriched Data

        par Async Cache Population
            API->>Redis: ZADD (ID, Score) & SET (Data)
        and Stream to Client
            API-->>Client: Stream JSON Chunk (Data)
        end
    end

    Note right of Client: Connection Closed
```
