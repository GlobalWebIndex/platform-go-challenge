# Go Favorites App

![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![Go Version](https://img.shields.io/badge/go-1.25-blue)
![Coverage](https://img.shields.io/badge/coverage-85%25-green)

A high-performance, **production-grade** Microservice for managing user favorites (Charts, Insights, Audiences).

This project demonstrates **System Thinking** and **Principal Engineering** patterns in Go, focusing on:

1. **Memory Safety**: O(1) memory usage via end-to-end Streaming (Iterators).
2. **Latency**: Write-Through Caching with Async Read-Repair.
3. **Resilience**: Graceful Shutdown, Circuit Breakers, and Backpressure handling.
4. **Observability**: Metric instrumentation and structured logging from Day 1.

---

## 📚 Documentation

Detailed documentation is available in the `docs/` directory:

* **[System Architecture](docs/ARCHITECTURE.md)**: Deep dive into Hexagonal Architecture, Data Flow, and C4 Diagrams.
* **[Architecture Decision Records (ADRs)](docs/DECISIONS.md)**: Why we chose JSONB, Go Iterators, and our Caching Strategy.
* **[Project Structure](docs/PROJECT_STRUCTURE.md)**: A map of the codebase for new contributors.

---

## Features

* **Go 1.25 Ready**: Utilizes modern features like `iter.Seq` (Iterators) and `omitzero` struct tags.
* **O(1) Memory Streaming**: End-to-end streaming from Database -> Service -> HTTP Response.
* **Sealed Interfaces**: Domain modeling using strict polymorphism.
* **Authentication**: JWT-based stateless authentication with Bcrypt password hashing.
* **Write-Through Caching**: Redis-backed `ZSET` caching for high-speed pagination.
* **Distroless Docker**: Secure, minimal production images.

## Technology Stack

* **Language**: Go 1.25
* **Database**: PostgreSQL 16 (pgx/v5)
* **Cache**: Redis 7 (go-redis/v9)
* **Auth**: JWT (v5), Bcrypt
* **Observability**: Prometheus, Slog, OpenTelemetry
* **Testing**: Testcontainers-go, k6

## Getting Started

### Prerequisites

* [Go 1.25+](https://go.dev/doc/install)
* [Docker & Docker Compose](https://docs.docker.com/get-docker/)
* [Make](https://www.gnu.org/software/make/)
* [k6](https://k6.io/docs/get-started/installation/) (for load testing)

### Running Locally

1. **Configuration**:

    Copy the example configuration file:

    ```bash
    cp .env.example .env
    ```

    *Tip: If port 8080 is in use, verify the `PORT` variable in `.env`.*

2. **Start Dependencies**:

    ```bash
    make docker-up
    ```

3. **Run the Server**:

    ```bash
    make run
    ```

4. **Test the API**:

    The API is protected. You must first acquire a token:

    ```bash
    # 1. Create User
    curl -X POST http://localhost:8080/signup -d '{"email":"test@example.com","password":"password123"}'

    # 2. Login
    TOKEN=$(curl -X POST http://localhost:8080/login -d '{"email":"test@example.com","password":"password123"}' | jq -r .token)

    # 3. Access Favorites
    curl -H "Authorization: Bearer $TOKEN" "http://localhost:8080/favorites?limit=10"
    ```

### Observability

* **Metrics**

Metrics are exposed at `http://localhost:8080/metrics`.

* **Distributed Tracing**

To enable tracing:

1. Ensure Jaeger is running (included in `make docker-up`).
2. Run the application with tracing enabled:

    ```bash
    make run-trace
    ```

3. View traces at [http://localhost:16686](http://localhost:16686).

### Running with Docker

```bash
docker build -t favorites-app .
docker run -p 8080:8080 -e DATABASE_URL=... favorites-app
```

## Testing

**Integration Tests** (using Testcontainers):

```bash
make test-integration
```

**Load Testing** (using k6):

```bash
k6 run tests/k6/load.js
```

## API Documentation

The OpenAPI spec is available at [api/openapi.yaml](api/openapi.yaml).
