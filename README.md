# 🧩 GWI Platform Go Challenge

This repository contains a small REST API built in Go for managing user favourites (charts, insights, and audiences).
It demonstrates a production‑ready, layered design with optional authentication, pagination, Swagger documentation, and environment‑based configuration.

---

## ⚙️ Tech Overview

**Language:** Go 1.25  
**Architecture:** layered (handler → service → repository) with `cmd/` + `internal/` layout  
**Storage:** in‑memory (thread‑safe with `sync.RWMutex`)  
**Containerization:** Docker + Docker Compose  
**Documentation:** Swagger UI (`swaggerapi/swagger-ui` container)  
**Testing:** built‑in Go test framework (`go test ./...`)  
**Configuration:** via `.env` (12‑factor style)  

---

## 🌐 API Overview

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET`  | `/users/{userID}/favourites` | List all favourites for a user (supports pagination) |
| `POST` | `/users/{userID}/favourites` | Create a new favourite |
| `PATCH`| `/users/{userID}/favourites/{favID}` | Update the description of a favourite |
| `DELETE` | `/users/{userID}/favourites/{favID}` | Delete a favourite |
| `GET`  | `/healthz` | Liveness probe |
| `GET`  | `/readyz` | Readiness probe |

---

## 🧱 Data Models (examples)

### Chart
```json
{
  "asset": {
    "type": "chart",
    "description": "Monthly sales chart",
    "title": "Sales 2024",
    "axis_x_title": "Month",
    "axis_y_title": "€",
    "data": [100, 150, 120, 200]
  }
}
```

### Insight
```json
{
  "asset": {
    "type": "insight",
    "description": "Insight on consumer behaviour",
    "text": "40% of users prefer mobile checkout."
  }
}
```

### Audience
```json
{
  "asset": {
    "type": "audience",
    "description": "Frequent social media users",
    "gender": "female",
    "birth_country": "Greece",
    "age_groups": ["25-34", "35-44"],
    "hours_social_daily": 2.5,
    "purchases_last_month": 3
  }
}
```

---

## 🔐 Authentication (optional)

The API supports **optional API key authentication** via middleware.  
By default, authentication is **disabled** (empty `API_KEY` in `.env`).  
To enable it, set `API_KEY` and include the header:

```
X-API-Key: <your_key>
```

Example:
```bash
curl -H "X-API-Key: topsecretkey" http://localhost:8080/healthz
```

---

## 📄 Pagination for Large Datasets

The service supports **pagination** to ensure fast response times even with thousands of favourites per user.

```
GET /users/{userID}/favourites?limit=100&offset=200
```

| Parameter | Description | Default | Max |
|------------|--------------|----------|------|
| `limit` | Number of results to return | 100 | 1000 |
| `offset` | Index to start returning results from | 0 | — |

Example response:
```json
{
  "favourites": [ ... ],
  "total": 1500,
  "limit": 100,
  "offset": 200
}
```

- Deterministic ordering (newest first)  
- Safe slicing via in-memory repository  
- Backward compatible (works with or without query params)

---

## 📘 API Documentation (Swagger UI)

Interactive API documentation is available at:

👉 **http://localhost:8081**

### How to enable
Swagger UI is automatically started as part of Docker Compose (`swaggerapi/swagger-ui` container).  
It uses the OpenAPI spec file (`openapi.yaml`) from the project root.

### CORS Integration
- The Go API (8080) exposes CORS for `http://localhost:8081`
- Swagger UI can directly call live API endpoints
- Fully testable via browser

### Example
Run:
```bash
curl http://localhost:8080/healthz
# {"status":"ok"}
```
Or execute `/healthz` directly inside Swagger UI — you'll see the same JSON response.

---

## 🗂️ Project Layout

```
platform-go-challenge/
├── cmd/
│   └── api/
│       └── main.go              # entrypoint
├── internal/
│   ├── config/                  # env-driven configuration
│   ├── middleware/              # logger, request id, security headers, rate limiter, body limit, api key
│   ├── models/                  # domain models
│   ├── repo/                    # repository interface + in-memory impl (thread-safe)
│   ├── service/                 # business logic + validation
│   └── server/                  # http handlers, routes, composition
├── Dockerfile
├── docker-compose.yml
├── .env
├── .env.example
└── README.md
```

---

## 🔧 Running Locally

### Prerequisites
- Go 1.25+ or Docker

### Option A — Go
```bash
go mod tidy
go run ./cmd/api
```
Server starts on `APP_PORT` (default `8080`).

### Option B — Docker
```bash
docker compose up --build -d
```

Test health:
```bash
curl http://localhost:8080/healthz
# {"status":"ok"}
```

---

## ⚙️ Configuration

The application reads runtime configuration from `.env`. Example:

```dotenv
APP_PORT=8080
APP_ENV=development
ENABLE_HTTP_LOG=true
RATE_LIMIT_MS=50
MAX_BODY_BYTES=1048576
READ_TIMEOUT=5
WRITE_TIMEOUT=10
IDLE_TIMEOUT=60
LOG_LEVEL=info
API_KEY=      # leave empty to disable auth
```

> After modifying `.env`, restart the container:  
> `docker compose down && docker compose up -d`

---

## 🧩 Example Usage

Create a Favourite
```bash
curl -X POST http://localhost:8080/users/kostas/favourites   -H "Content-Type: application/json"   -d '{"asset":{"type":"insight","description":"market trend","text":"40% of users..."}}'
```

List Favourites (paged)
```bash
curl "http://localhost:8080/users/kostas/favourites?limit=3&offset=0"
```

Update Description
```bash
curl -X PATCH http://localhost:8080/users/kostas/favourites/<favID>   -H "Content-Type: application/json"   -d '{"description":"updated insight"}'
```

Delete Favourite
```bash
curl -X DELETE http://localhost:8080/users/kostas/favourites/<favID>
```

---

## 🧪 Running Tests

```bash
go test ./... -v
```
Covers:
- Service layer validation & CRUD logic
- HTTP endpoints including pagination and validation
- Edge cases and invalid payloads

---

## 🗒️ Design notes

- Concurrency safety: in-memory repository guarded with `sync.RWMutex` (parallel reads, single writer).
- Separation of concerns: handlers → service → repository; middleware for cross‑cutting concerns.
- Production hygiene: env-based config, timeouts, request size limit, basic rate limiting, health/readiness.
- Pagination ensures scalability for large datasets while keeping latency minimal.

---

## 🚀 Future Enhancements

- **Persistent storage** — Replace the in-memory repository with PostgreSQL or Redis, adding proper indexing, migrations, and connection pooling for scalability.  
- **Advanced authentication & authorization** — Extend the current API-key approach with JWTs and role-based access control for multi-tenant setups.  
- **Observability & tracing** — Add structured logging and distributed tracing via **OpenTelemetry**, exporting metrics to Prometheus and traces to Jaeger or Grafana Tempo.  
- **Operational insights** — Expose a `/metrics` endpoint (Prometheus format) for request latency, throughput, and error-rate monitoring. Combine with Grafana dashboards for real-time health.  
- **Async job processing** — Introduce a **worker-pool pattern** for background or heavy tasks such as bulk favourites export.  
  - Jobs would be queued in RabbitMQ, Redis Streams, or AWS SQS.  
  - Workers consume from the queue with bounded concurrency, support retries with exponential back-off, and expose job status via `/exports/{id}` (HTTP 202 → poll for result).  
  - This approach keeps API latency low while handling large-scale data exports safely.  

---

**Author:** Konstantinos Dasios  
**Date:** November 2025  
**Challenge:** GWI Engineering Manager – Platform Go Challenge
