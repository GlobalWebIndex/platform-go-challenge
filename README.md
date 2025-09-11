# GWI Favourites API (Go Challenge)

This project is a Go web server for managing user favourites (assets: Chart, Insight, Audience) with full CRUD endpoints, JWT authentication, pagination, Docker support, and automated testing.

## Features
- RESTful API for user favourites
- Asset types: Chart, Insight, Audience
- JWT-based authentication for all endpoints
- Pagination support for listing favourites
- In-memory storage with concurrency safety
- Dockerfile for containerization
- Makefile for build, run, lint, and test automation
- Postman collection for API testing

## Getting Started

### Prerequisites
- Go 1.20+
- Docker (optional)

### Build & Run (Locally)
```bash
make build      # Build the Go binary
make run        # Run the server locally (default: :8080)
```

### Build & Run (Docker)
```bash
docker build -t gwi-favourites .
docker run -p 8080:8080 gwi-favourites
```

### Test
```bash
make test           # Run all unit tests
```

## API Endpoints

All endpoints (except `/token`) require JWT authentication via the `Authorization: Bearer <TOKEN>` header.

### Authentication
- **POST /token**
  - Request: `{ "user_id": "<USER_UUID>" }`
  - Response: `{ "token": "<JWT_TOKEN>" }`
  - Use this endpoint to obtain a JWT for a user.

### Favourites
- **GET /favourites?limit=10&offset=0**
  - List all favourite assets for the authenticated user (supports pagination).
- **POST /favourites/add**
  - Add a new asset (Chart, Insight, Audience) to favourites.
  - Request body: `{ "type": "chart|insight|audience", "favorite": true|false, "asset": { ... } }`
- **PUT /favourites/remove?asset_id=<ASSET_UUID>**
  - Update the `favorite` status of an asset.
  - Request body: `{ "favorite": true|false }`
- **PUT /favourites/edit?asset_id=<ASSET_UUID>**
  - Edit the description of an asset.
  - Request body: `{ "description": "..." }`
- **DELETE /favourites/delete?asset_id=<ASSET_UUID>**
  - Delete an asset from favourites.

## Asset Types
- **Chart**: `{ "Title", "XAxisTitle", "YAxisTitle", "Data", "Description" }`
- **Insight**: `{ "Text", "Description" }`
- **Audience**: `{ "Gender", "BirthCountry", "AgeGroup", "SocialHours", "Purchases", "Description" }`

## Authentication Flow
- Obtain a JWT via `/token` by providing a valid user UUID.
- Include the JWT in the `Authorization` header for all other requests.
- The server extracts the user ID from the token and uses it to scope all data access.

## Running Tests
```bash
make test
```
All tests use JWT authentication and validate the full API contract.

## Postman Collection
A ready-to-use Postman collection (`collection.json`) is included for testing all endpoints, including authentication.

## Notes
- The server uses in-memory storage; data resets on restart.
- For production, use a secure JWT secret and persistent storage.

---
