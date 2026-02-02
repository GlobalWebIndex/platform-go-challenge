# GlobalWebIndex Engineering Challenge – Favorites API (Go)

## Overview
This project implements a lightweight HTTP API that allows users to manage a list of favorite assets.

The solution focuses on:
- **Clean API design** (simple REST endpoints)
- **Type-safe validation** (type-specific payload structs per asset type)
- **Concurrency-safe in-memory storage** (RWMutex-protected store)
- **No external dependencies** (standard library only)

---

## Quickstart

### Prerequisites
- Go 1.20+
- Git
- Docker 

--

### Run with Docker

Build the image:
```bash
docker build -t gwi-favorites .

Run the container:

docker run --rm -p 8080:8080 gwi-favorites

Health check:

curl http://localhost:8080/health

### Run locally
```bash
go run ./cmd/server

Server runs at:
http://localhost:8080

Health check:
curl http://localhost:8080/health

Run tests
go test ./...

API Endpoints

List favorites
GET /users/{userId}/favorites

Response:

{
  "items": [ ... ]
}

Add favorite
POST /users/{userId}/favorites

Request Body: an asset with supported type and valid payload.

Update description
PATCH /users/{userId}/favorites/{assetId}

Request Body:

{
  "description": "updated description"
}

Remove favorite
DELETE /users/{userId}/favorites/{assetId}

Returns:

- 204 No Content on success

- 404 Not Found if the asset does not exist

Asset Model
All asset types share a common wrapper:

{
  "id": "a1",
  "type": "insight",
  "description": "optional description",
  "payload": { ... }
}
The payload structure depends on type and is validated accordingly.

Supported Asset Types & Payload Examples

1) Insight
Represents a short piece of text.

Example:

{
  "id": "a1",
  "type": "insight",
  "description": "Example insight",
  "payload": {
    "text": "Users spend more time on mobile devices"
  }
}
Validation:

text must be non-empty

2) Chart
Represents a chart with minimal metadata.

Example:

{
  "id": "c1",
  "type": "chart",
  "description": "Sales performance",
  "payload": {
    "title": "Monthly Sales",
    "axes": {
      "x": "Month",
      "y": "Revenue"
    },
    "data": [120, 150, 180]
  }
}

Validation:

- title must be present

- both axis labels (x, y) must be non-empty

- data is intentionally flexible to support different chart shapes

3) Audience
Represents an audience definition based on one or more dimensions.

Example:

{
  "id": "aud1",
  "type": "audience",
  "description": "Young adults in EU",
  "payload": {
    "ageGroups": ["18-24", "25-34"],
    "birthCountry": ["GR", "DE"]
  }
}

Validation:

- At least one audience dimension must be provided

Validation Strategy

- Payloads are decoded into type-specific structs

- Only essential fields are validated (pragmatic correctness)

- Schema is kept flexible to avoid over-constraining future extensions

Invalid payloads return:

- 400 Bad Request

Storage & Concurrency

- The service uses an in-memory store backed by:

- userId -> assetId -> asset

Concurrency guarantees:

- reads via RWMutex to allow parallel access

- writes serialized to ensure consistency

- Production alternatives (discussion points):

- PostgreSQL (JSONB payloads)

- Document store (e.g., MongoDB)

- Redis cache + persistent storage

Testing

- Handler-level tests using httptest cover:

- add + list favorites

- validation failures

- update description

- remove favorite

Run:

go test ./...

Trade-offs & Next Steps

Trade-offs:

- No external router library (minimal dependencies)

- Validation focuses on essentials, not strict schemas

- In-memory store for clarity and speed

Possible next steps:

- pagination for list endpoint

- authentication/authorization middleware

- persistent storage

- structured logging + metrics