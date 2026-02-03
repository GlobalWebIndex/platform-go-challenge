# Architecture Decision Records (ADRs)

This document captures the key architectural decisions made during the development of the Favorites Service, including context, decision, and consequences.

## ADR 001: JSONB for Asset Storage

* **Status**: Accepted
* **Context**: The application handles multiple types of "Assets" (Charts, Insights, Audiences). These assets have vastly different fields and structures. Creating separate tables for each would require complex joins and frequent schema migrations as new asset types are introduced.
* **Decision**: Use PostgreSQL `JSONB` column to store the asset payload, with a top-level `type` discriminator. Use a `GIN` index on the `asset_data` column to support efficient querying if needed.
* **Consequences**:
  * **Pros**: Schema flexibility. New asset types can be added without database migrations. Polyglot persistence within a relational database.
  * **Cons**: Slightly higher storage size compared to normalized columns. Complex queries on JSON fields are more verbose.

## ADR 002: Write-Through Caching Strategy with Background Enrichment

* **Status**: Accepted
* **Context**: Users expect extremely fast read access to their favorites list. The list must be enriched with metadata from an external service, which can be slow. Blocking the "List Favorites" request to call an external API for every item is unacceptable for latency.
* **Decision**: Implement a **Write-Through** pattern with an **Atomic Background Worker**.
    1. When an item is Saved, we trigger an async background job to enrich it and store it in Redis.
    2. When `FindAll` is called, we first check Redis.
    3. If a Cache Miss occurs, we stream from DB and defensively trigger a background cache update for subsequent requests (Read-Repair).
* **Consequences**:
  * **Pros**: Read latency is decoupled from the external Enrichment Service. Cache is kept fresh.
  * **Cons**: Eventual consistency for the first read after a save if the background job is slow. Complexity in managing background goroutines during shutdown.

## ADR 003: Go 1.25 Iterators for Streaming

* **Status**: Accepted
* **Context**: The system must support "unlimited favorites" (returning potentially millions of rows) without exhausting server memory. Traditional slice-based returns (`[]Asset`) load the entire dataset into RAM.
* **Decision**: Adopt the Go 1.25 `iter` package (`iter.Seq2`).
* **Consequences**:
  * **Pros**: True O(1) memory usage for the request pipeline. Standardized streaming interface across layers. Replaces error-prone channel-based streaming implementations.
  * **Cons**: `iter` is a relatively new feature in Go 1.23+, requiring Go 1.25 for full standard library support and ecosystem tooling.

## ADR 004: Sealed Interfaces for Domain Modeling

* **Status**: Accepted
* **Context**: The `Asset` domain object is polymorphic. We need to handle `Chart`, `Insight`, and `Audience` safely without type assertions scattered everywhere.
* **Decision**: Use a Sealed Interface pattern where the `Asset` interface contains a private method, ensuring only defined types in the domain package can implement it.
* **Consequences**:
  * **Pros**: Compile-time safety. Exhaustive type checking in switch statements is easier to enforce.
  * **Cons**: Slightly more boilerplate code for the interface definition.

## ADR 005: JWT Authentication

* **Status**: Accepted
* **Context**: The service requires user attribution for "favorites". We need a way to identify users across requests without maintaining server-side session state (statelessness), as this microservice mimics a high-scale environment.
* **Decision**: Implement **JSON Web Tokens (JWT)** signed with HS256 (HMAC). Passwords are hashed using **Bcrypt** (cost 10) before storage.
* **Consequences**:
  * **Pros**: Stateless authentication scales horizontally. Decouples the Auth verification from the database (once the key is known/distributed, though currently monolithic).
  * **Cons**: Token invalidation (logout) is difficult without a blocklist (not implemented yet). Clients must manage token storage securely.
