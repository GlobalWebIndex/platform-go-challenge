# API Shape – Assets Favourites API

Base URL: `/v1`

The surface is kept small and focused on the core use cases.

## Standards

Industry standards are followed to keep things consistent and maintainable:

**Format & Protocol:**

- **JSON:** All requests and responses use JSON (Content-Type: `application/json`)
- **HTTP:** RFC 9110 (HTTP Semantics) - HTTP methods (GET, POST, PATCH, DELETE) and status codes (200, 201, 204, 400, 401, 404, 409, 422, 500)
- **Timestamps:** ISO 8601 format in UTC (e.g., `2025-01-01T12:00:00Z`)
- **Error format:** RFC 9457 (Problem Details for HTTP APIs)

**REST Design (following Google API Design Guide principles):**

- Resource-based URLs (e.g., `/v1/users/{userId}/favourites`)
- HTTP methods: GET (read), POST (create), PATCH (update), DELETE (remove)
- Versioning in path (`/v1/...`)
- Plural nouns for collections (`favourites`, not `favourite`)
- camelCase for JSON keys

**Idempotency:**

- DELETE operations are idempotent (204 is returned even if resource doesn't exist)
- PATCH operations are idempotent

**References:**

- [RFC 9110](https://www.rfc-editor.org/rfc/rfc9110) - HTTP Semantics (HTTP methods: Section 9.3, Status codes: Section 15)
- [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457) - Problem Details for HTTP APIs
- [RFC 8259](https://tools.ietf.org/html/rfc8259) - JSON Format
- [Google API Design Guide](https://cloud.google.com/apis/design) - REST API design principles
- [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) - Date and time representation

---

## 1. List user favourites

**Endpoint**

- `GET /v1/users/{userId}/favourites`

**Query params (optional)**

- `type` – filter by asset type (`chart`, `insight`, `audience`)
- `limit` – maximum number of results (accepted but not yet implemented; all results are returned)
- `offset` – number of results to skip (accepted but not yet implemented)

**Response 200**

```json
{
  "items": [
    {
      "userId": "user-123",
      "assetId": "asset-abc",
      "createdAt": "2025-01-01T12:00:00Z",
      "asset": {
        "id": "asset-abc",
        "type": "chart",
        "description": "Chart description",
        "chart": {
          "title": "Example chart",
          "xAxis": "Age",
          "yAxis": "Score",
          "data": [1.2, 3.4, 5.6]
        }
      }
    }
  ],
  "meta": {
    "total": 42,
    "count": 42
  }
}
```

**Error Responses**

- `400 Bad Request` – Invalid query parameters
- `401 Unauthorized` – Authentication required or invalid
- `404 Not Found` – User not found
- `422 Unprocessable Entity` – Validation error
- `500 Internal Server Error` – Server error

---

## 2. Add favourite

**Endpoint**

- `POST /v1/users/{userId}/favourites`

**Request body**

Create a new asset inline:

```json
{
  "asset": {
    "type": "insight",
    "description": "Key finding",
    "insight": {
      "text": "People who do X are likely to do Y."
    }
  }
}
```

**Response 201**

```json
{
  "userId": "user-123",
  "assetId": "asset-abc",
  "createdAt": "2025-01-01T12:00:00Z"
}
```

The asset is stored in the asset store and an `id` is assigned (auto-generated UUID).

**Error Responses**

- `400 Bad Request` – Invalid request body or missing required fields
- `401 Unauthorized` – Authentication required or invalid
- `404 Not Found` – User not found
- `409 Conflict` – Favourite already exists for this user-asset pair
- `422 Unprocessable Entity` – Validation error (invalid asset data or format)
- `500 Internal Server Error` – Server error

---

## 3. Remove favourite

**Endpoint**

- `DELETE /v1/users/{userId}/favourites/{assetId}`

**Behavior**

- Idempotent: success is returned even if the favourite did not exist.

**Response 204**

- Empty body.

**Error Responses**

- `400 Bad Request` – Invalid parameters
- `401 Unauthorized` – Authentication required or invalid
- `404 Not Found` – User not found
- `500 Internal Server Error` – Server error

---

## 4. Edit favourite description

**Endpoint**

- `PATCH /v1/users/{userId}/favourites/{assetId}`

**Request body**

```json
{
  "description": "New description"
}
```

**Response 200**

```json
{
  "userId": "user-123",
  "assetId": "asset-abc",
  "createdAt": "2025-01-01T12:00:00Z",
  "asset": {
    "id": "asset-abc",
    "type": "chart",
    "description": "New description",
    "chart": {
      "title": "Example chart",
      "xAxis": "Age",
      "yAxis": "Score",
      "data": [1.2, 3.4, 5.6]
    }
  }
}
```

**Error Responses**

- `400 Bad Request` – Invalid request body
- `401 Unauthorized` – Authentication required or invalid
- `404 Not Found` – User or favourite not found
- `422 Unprocessable Entity` – Validation error
- `500 Internal Server Error` – Server error

---

## 5. Health check

**Endpoint**

- `GET /health`

**Behavior**

- Returns the health status of the service. This endpoint is kept simple for the take-home assignment.
- No authentication required (typically public for load balancer and orchestration system health checks).

**Response 200**

```json
{
  "status": "healthy"
}
```

**Error Responses**

- This endpoint does not return error responses. If the service is running, it returns 200 with `{"status": "healthy"}`.

---

## 6. Metrics

**Endpoint**

- `GET /metrics`

**Behavior**

- Returns Prometheus-formatted metrics for monitoring and observability.
- Exposes HTTP request counts, latency, error rates, and operation metrics.
- Requires authentication (same `mockAuth` as other endpoints).

**Response 200**

Content-Type: `text/plain`

Prometheus exposition format:

```
# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",endpoint="/v1/users/{userId}/favourites",status="200"} 42
http_requests_total{method="POST",endpoint="/v1/users/{userId}/favourites",status="201"} 10
# HELP http_request_duration_seconds HTTP request duration in seconds
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{method="GET",endpoint="/v1/users/{userId}/favourites",le="0.1"} 35
http_request_duration_seconds_bucket{method="GET",endpoint="/v1/users/{userId}/favourites",le="0.5"} 40
http_request_duration_seconds_bucket{method="GET",endpoint="/v1/users/{userId}/favourites",le="+Inf"} 42
http_request_duration_seconds_sum{method="GET",endpoint="/v1/users/{userId}/favourites"} 8.5
http_request_duration_seconds_count{method="GET",endpoint="/v1/users/{userId}/favourites"} 42
```

**Error Responses**

- `401 Unauthorized` – Authentication required or invalid

**Note:** In production environments, metrics endpoints should be protected with stronger security measures (e.g., network-level restrictions, service mesh policies, or dedicated monitoring credentials). For this take-home assignment, the same authentication mechanism (`mockAuth`) is used as the rest of the API.

---

## Error Format

Errors follow RFC 9457 format. Every error includes `type`, `title`, `status`, `detail`, and `instance`.

**Example (400 Bad Request):**

```json
{
  "type": "https://api.example.com/problems/bad-request",
  "title": "Bad Request",
  "status": 400,
  "detail": "Invalid request parameters",
  "instance": "/v1/users/user-123/favourites"
}
```

**Example (404 Not Found):**

```json
{
  "type": "https://api.example.com/problems/not-found",
  "title": "Not Found",
  "status": 404,
  "detail": "User 'user-123' not found",
  "instance": "/v1/users/user-123/favourites"
}
```

---

## Validation Rules

Basic validation is applied to ensure data integrity. For this challenge, validation is kept simple and focused on required fields and basic type checking.

### Required Fields

- All required fields must be present and non-empty
- Asset `type` must be one of: `chart`, `insight`, `audience`
- Type-specific fields must match the asset type (e.g., `chart` asset must have `chart` object)

### Basic Type Validation

- User IDs and Asset IDs: non-empty strings
- Asset descriptions: non-empty strings
- Chart: `title`, `xAxis`, `yAxis` (strings), `data` (array of numbers)
- Insight: `text` (string)
- Audience: `gender` (enum: `male`, `female`, `other`), `birthCountry` (string), `ageGroup` (string), `hoursOnSocialMediaPerDay` (integer), `purchasesLastMonth` (integer)

Note: More detailed validation (length limits, patterns, ranges) can be added later as needed. For this challenge, focus is placed on ensuring required fields are present and types are correct.

---

## Notes

- **Asset ID Generation:** When creating assets inline, IDs are auto-generated (UUID v4 format)
- **Pagination:** The `limit` and `offset` parameters are accepted but not yet implemented. All results are returned regardless of these parameters.
- **Idempotency:** DELETE operations are idempotent and 204 is returned even if the resource doesn't exist.
- **Polymorphism:** Asset types are strictly enforced. Only the relevant type-specific field (`chart`, `insight`, or `audience`) should be present based on the `type` field.
