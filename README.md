# Platform Go Challenge - Favorites Management API

A RESTful API service for managing user favorites (charts, insights, and audiences) built with Go, featuring JWT authentication, rate limiting, and in-memory storage.

## Table of Contents

- [Platform Go Challenge - Favorites Management API](#platform-go-challenge---favorites-management-api)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Architecture](#architecture)
  - [Prerequisites](#prerequisites)
  - [Quick Start](#quick-start)
    - [Option 1: Using Make (Recommended)](#option-1-using-make-recommended)
    - [Option 2: Using Docker Compose](#option-2-using-docker-compose)
    - [Option 3: Manual Setup](#option-3-manual-setup)
  - [Environment Variables](#environment-variables)
  - [API Endpoints](#api-endpoints)
    - [Authentication](#authentication)
      - [Login](#login)
    - [Health Check](#health-check)
      - [Get Server Health](#get-server-health)
    - [Favorites Management](#favorites-management)
      - [Get User Favorites](#get-user-favorites)
      - [Add Favorite](#add-favorite)
      - [Update Favorite Description](#update-favorite-description)
      - [Remove Favorite](#remove-favorite)
  - [Testing with Postman](#testing-with-postman)
    - [Import Collection](#import-collection)
    - [Setup](#setup)
    - [Available Requests](#available-requests)
  - [Testing](#testing)
    - [Run All Tests](#run-all-tests)
    - [Run Tests with Coverage](#run-tests-with-coverage)
    - [Generate HTML Coverage Report](#generate-html-coverage-report)
    - [Run Tests with Race Detector](#run-tests-with-race-detector)
    - [Test Structure](#test-structure)
    - [What's Tested](#whats-tested)
  - [Docker Support](#docker-support)
    - [Build Image](#build-image)
    - [Run with Docker Compose](#run-with-docker-compose)
    - [Docker Configuration](#docker-configuration)
  - [Error Handling](#error-handling)
  - [Rate Limiting](#rate-limiting)
  - [Middleware Stack](#middleware-stack)
  - [Concurrency \& Thread Safety](#concurrency--thread-safety)
  - [Design Decisions](#design-decisions)
    - [In-Memory Storage](#in-memory-storage)
    - [JWT Authentication](#jwt-authentication)
    - [Asset Union Type](#asset-union-type)
    - [Rate Limiting](#rate-limiting-1)
    - [Graceful Shutdown](#graceful-shutdown)
  - [Makefile Commands](#makefile-commands)
  - [Project Structure Philosophy](#project-structure-philosophy)
  - [Future Enhancements](#future-enhancements)
  - [Development](#development)
    - [Code Style](#code-style)
    - [Adding New Asset Types](#adding-new-asset-types)
  - [Problems \& Troubleshooting](#problems--troubleshooting)
    - [Port Already in Use](#port-already-in-use)
    - [Docker Issues](#docker-issues)
    - [Tests Failing](#tests-failing)
  - [License](#license)
  - [Author](#author)
  - [Acknowledgments](#acknowledgments)

## Features

- JWT-based authentication
- Rate limiting (5 req/min for login, 50 req/min for other endpoints)
- Support for multiple asset types (Chart, Insight, Audience)
- Thread-safe in-memory storage
- 100% test coverage on core business logic
- Docker support with Docker Compose
- Request metrics and health monitoring
- Clean architecture with separation of concerns
- Comprehensive error handling
- Graceful shutdown

## Architecture

```
platform-go-challenge/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── api_errors/          # Custom error definitions
│   |   └── errors.go
│   ├── auth/                # JWT authentication
│   │   └── auth.go          
│   ├── handler/             # HTTP request handlers
│   │   ├── handler.go       # Favorites handlers
│   │   ├── auth.go          # Login handler
│   │   ├── handler_test.go  # Favorites tests
│   │   └── auth_test.go     # Login tests
│   ├── middleware/          # HTTP middleware (metrics, auth, rate limiting)
│   │   ├── auth.go          
│   │   ├── logger.go        
│   │   └── rate_limit.go    
│   ├── model/               # Data models
│   │   ├── asset.go
│   │   └── asset_test.go
│   ├── storage/             # In-memory storage implementation
│   │   ├── store.go
│   │   └── store_test.go
│   └── util/                # Utility functions
│       └── response.go
├── postman                  # Postman API collection
|   ├── postman_collection.json
|   └── README.md
├── .dockerignore            # Docker ignore file
├── .gitignore              # Git ignore file
├── docker-compose.yml       # Docker Compose setup
├── Dockerfile               # Docker build configuration
├── go.mod                   # Go module file
├── go.sum                   # Go module checksum
├── makefile                 # Build automation
└── README.md
```

## Prerequisites

- Go 1.21 or higher
- Docker & Docker Compose (optional)
- Make (optional, for using Makefile commands)

## Quick Start

### Option 1: Using Make (Recommended)

```bash
# View available commands
make help

# Run locally
make run

# Run tests
make test

# Build binary
make build

# Docker
make docker-up
make docker-down
```

### Option 2: Using Docker Compose

```bash
# Start the application
docker compose up

# Stop the application
docker compose down
```

### Option 3: Manual Setup

```bash
# Install dependencies
go mod download

# Run the application
go run cmd/server/main.go

# Or build and run
go build -o bin/server ./cmd/server
./bin/server
```

The server will start on `http://localhost:8080`

## Environment Variables

```bash
JWT_SECRET=your-secret-key-here    # JWT signing key (optional, defaults to empty)
PORT=8080                          # Server port (optional, defaults to 8080)
```

## API Endpoints

### Authentication

#### Login
```http
POST /api/login
Content-Type: application/json

{
  "user_id": "user-123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": "user-123"
}
```

**Rate Limit:** 5 requests per minute per IP

### Health Check

#### Get Server Health
```http
GET /api/health
```

**Response:**
```json
{
  "status": "healthy",
  "total_requests": 42,
  "active_requests": 0
}
```

**Note:** No authentication required

### Favorites Management

All endpoints below require JWT authentication via `Authorization: Bearer <token>` header.

**Rate Limit:** 50 requests per minute per IP

#### Get User Favorites
```http
GET /api/users/{id}/favorites
Authorization: Bearer <token>
```

**Response:**
```json
{
  "user_id": "user-123",
  "favorites": [
    {
      "asset": {
        "chart": {
          "id": "chart-1",
          "type": "chart",
          "title": "Sales Chart",
          "description": "Q4 sales data",
          "x_axis_key": "month",
          "y_axis_key": "revenue",
          "data": [...],
          "created_at": 1234567890,
          "updated_at": 1234567890
        }
      },
      "added_at": 1234567890
    }
  ],
  "count": 1
}
```

#### Add Favorite
```http
POST /api/users/{id}/favorites
Authorization: Bearer <token>
Content-Type: application/json
```

**Asset Types:**

1. **Chart**
```json
{
  "asset": {
    "chart": {
      "id": "chart-001",
      "type": "chart",
      "title": "Sales by Region",
      "description": "Monthly sales performance",
      "x_axis_key": "month",
      "y_axis_key": "sales",
      "data": [
        {"month": "Jan", "sales": 1200},
        {"month": "Feb", "sales": 1900}
      ]
    }
  }
}
```

2. **Insight**
```json
{
  "asset": {
    "insight": {
      "id": "insight-001",
      "type": "insight",
      "title": "Millennial Social Media Usage",
      "description": "Social media engagement trend",
      "content": "40% of millennials spend more than 3 hours on social media daily"
    }
  }
}
```

3. **Audience**
```json
{
  "asset": {
    "audience": {
      "id": "audience-001",
      "type": "audience",
      "description": "High-value male customers",
      "gender": "Male",
      "countries": ["US", "UK", "Canada"],
      "age_min": 24,
      "age_max": 35,
      "min_hours_daily_social_media": 3,
      "min_purchases_last_month": 5
    }
  }
}
```

**Response:**
```json
{
  "message": "Favorite added successfully",
  "asset_id": "chart-001"
}
```

#### Update Favorite Description
```http
PUT /api/users/{id}/favorites/{assetId}
Authorization: Bearer <token>
Content-Type: application/json

{
  "description": "Updated description"
}
```

**Response:**
```json
{
  "message": "Favorite updated successfully",
  "asset_id": "chart-001"
}
```

#### Remove Favorite
```http
DELETE /api/users/{id}/favorites/{assetId}
Authorization: Bearer <token>
```

**Response:**
```json
{
  "message": "Favorite removed successfully",
  "asset_id": "chart-001"
}
```

## Testing with Postman

A complete Postman collection is included in `postman_collection.json`.

### Import Collection

1. Open Postman
2. Click **Import** → **File**
3. Select `postman_collection.json`
4. Collection "GWI API" will appear

### Setup

1. Run **Login (Get Token)** request
2. Copy the returned `token` value
3. Click on the "GWI API" collection
4. Go to **Variables** tab
5. Set `{{token}}` variable to the copied token
6. All authenticated requests will now use this token automatically

### Available Requests

- Login (Get Token)
- Health Check (No Auth)
- Get User Favorites
- Add Chart Favorite
- Add Insight Favorite
- Add Audience Favorite
- Edit Favorite Description
- Remove Favorite
- Test Missing Token (401)
- Test Invalid Token (401)

## Testing

The project includes comprehensive unit and integration tests with 100% coverage on core business logic.

### Run All Tests
```bash
make test
# or
go test ./...
```

### Run Tests with Coverage
```bash
go test -cover ./...
```


### Generate HTML Coverage Report
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run Tests with Race Detector
```bash
go test -race ./...
```

### Test Structure

```
internal/
├── handler/
│   ├── handler_test.go    # Favorites handlers (100% coverage)
│   └── auth_test.go       # Authentication (100% coverage)
├── storage/
│   └── store_test.go      # Storage layer (100% coverage)
└── model/
    └── asset_test.go      # Models (100% coverage)
```

### What's Tested

**Handler Tests:**
- All CRUD operations for favorites
- JWT authentication flows
- Input validation
- Error handling (404, 400, 401, 500)
- Store error scenarios
- Integration workflow tests

**Storage Tests:**
- Thread-safe operations
- Concurrent access patterns
- CRUD operations
- Error cases
- Timestamp management

**Model Tests:**
- GetBaseAsset() for all asset types
- Edge cases (nil assets)

## Docker Support

### Build Image
```bash
make docker-build
# or
docker compose build
```

### Run with Docker Compose
```bash
# Start in foreground
docker compose up

# Start in background
docker compose up -d

# View logs
docker compose logs -f

# Stop
docker compose down
```

### Docker Configuration

- **Multi-stage build** for optimized image size
- **Health checks** every 10 seconds
- **Auto-restart** unless stopped
- **Port 8080** exposed

## Error Handling

The API returns standard HTTP status codes:

| Status Code | Description |
|------------|-------------|
| `200 OK` | Successful GET/PUT/DELETE |
| `201 Created` | Successful POST |
| `400 Bad Request` | Invalid input/JSON |
| `401 Unauthorized` | Missing or invalid JWT token |
| `404 Not Found` | Resource not found |
| `429 Too Many Requests` | Rate limit exceeded |
| `500 Internal Server Error` | Server error |

**Error Response Format:**
```json
{
  "error": "Error message description"
}
```

## Rate Limiting

Rate limits are applied per IP address:

| Endpoint | Limit |
|----------|-------|
| `POST /api/login` | 5 requests/minute |
| All other endpoints | 50 requests/minute |

When rate limit is exceeded:
```json
{
  "error": "Rate limit exceeded. Please try again later."
}
```

## Middleware Stack

The application uses the following middleware chain:

1. **Request Logger** - Logs all incoming requests
2. **Rate Limiter** - Enforces rate limits per endpoint
3. **Metrics Collector** - Tracks request counts
4. **Auth Middleware** - Validates JWT tokens (protected routes only)

## Concurrency & Thread Safety

The in-memory storage implementation uses `sync.RWMutex` to ensure thread-safe operations:
- Multiple concurrent reads are allowed
- Writes are exclusive and block other operations
- Tested with concurrent access patterns using `-race` flag

## Design Decisions

### In-Memory Storage
- Simple and fast for MVP/prototype
- Thread-safe implementation with RWMutex
- Easy to replace with persistent storage (implements `Store` interface)

### JWT Authentication
- Stateless authentication
- 24-hour token expiry
- Standard claims with user identification
- Configurable secret key

### Asset Union Type
- Flexible design supporting multiple asset types
- Type-safe with compile-time checks
- Easy to extend with new asset types
- Single `GetBaseAsset()` method for common fields

### Rate Limiting
- IP-based rate limiting
- Different limits for sensitive endpoints
- Token bucket algorithm for smooth traffic control

### Graceful Shutdown
- Handles SIGINT and SIGTERM signals
- Allows in-flight requests to complete
- Clean resource cleanup

## Makefile Commands

```bash
make help          # Show available commands
make build         # Build the Go binary
make run           # Run the server locally
make test          # Run tests with coverage
make docker-build  # Build Docker image
make docker-up     # Start with Docker Compose
make docker-down   # Stop Docker Compose
make clean         # Remove binaries and containers
```

## Project Structure Philosophy

The project follows these principles:

- **`cmd/`** - Application entry points
- **`internal/`** - Private application code
  - **`handler/`** - HTTP layer (controllers)
  - **`storage/`** - Data access layer
  - **`model/`** - Domain models
  - **`middleware/`** - Cross-cutting concerns
  - **`auth/`** - Authentication logic
  - **`util/`** - Shared utilities

Pros:
- Clear separation of concerns
- Easy to test each layer independently
- Simple to swap implementations (e.g., storage)
- Maintainable and scalable

## Future Enhancements

- [ ] Persistent storage (PostgreSQL/MongoDB)
- [ ] Pagination for favorites list
- [ ] Filtering and sorting favorites
- [ ] User registration and management
- [ ] Refresh token mechanism
- [ ] CI/CD pipeline
- [ ] Redis for distributed rate limiting
- [ ] API documentation (Swagger/OpenAPI) & docs endpoint
- [ ] Prometheus metrics export
- [ ] Request tracing (OpenTelemetry)
- [ ] WebSocket support for real-time updates

## Development

### Code Style
The project follows standard Go conventions:
- `gofmt` for formatting
- Descriptive function and variable names
- Comprehensive error handling
- Clear separation of concerns

### Adding New Asset Types

1. Add new type to `internal/model/asset.go`:
```go
type NewAsset struct {
    BaseAsset `json:",inline"`
    CustomField string `json:"custom_field"`
}
```

2. Update `Asset` struct:
```go
type Asset struct {
    Chart    *Chart    `json:"chart,omitempty"`
    Insight  *Insight  `json:"insight,omitempty"`
    Audience *Audience `json:"audience,omitempty"`
    NewAsset *NewAsset `json:"new_asset,omitempty"` // Add this
}
```

3. Update `GetBaseAsset()` method:
```go
func (a *Asset) GetBaseAsset() *BaseAsset {
    switch {
    case a.Chart != nil:
        return &a.Chart.BaseAsset
    case a.Insight != nil:
        return &a.Insight.BaseAsset
    case a.Audience != nil:
        return &a.Audience.BaseAsset
    case a.NewAsset != nil:
        return &a.NewAsset.BaseAsset // Add this
    default:
        return nil
    }
}
```

4. Add tests in `internal/model/asset_test.go`

## Problems & Troubleshooting

### Port Already in Use
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

### Docker Issues
```bash
# Clean up Docker resources
make clean

# Rebuild from scratch
docker compose build --no-cache
docker compose up
```

### Tests Failing
```bash
# Clear test cache
go clean -testcache

# Run specific test
go test -v ./internal/handler -run TestLogin
```

## License

MIT License

## Author

Built for a Go REST API technical challenge.

## Acknowledgments

- Built with Go standard library
- JWT implementation using [golang-jwt/jwt](https://github.com/golang-jwt/jwt)
- Testing with [testify](https://github.com/stretchr/testify)