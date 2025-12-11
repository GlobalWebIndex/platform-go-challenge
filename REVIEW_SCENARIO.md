# Review Scenario – Assets Favourites API

This document is a quick, copy-paste friendly walk-through for reviewing the API.
It exercises all main endpoints with minimal setup.

---

## 0. Start the API

Choose one of the following methods to start the API:

### Option 1: Using Go directly

```bash
go run ./cmd/api
```
or 
```bash
HTTP_ADDR=:8080 go run ./cmd/api
```

### Option 2: Using Make

```bash
make run
```

### Option 3: Using Docker Compose (includes Prometheus + Grafana)

```bash
docker compose up --build
```

This starts:

- API on `http://localhost:8080`
- Prometheus on `http://localhost:9090`
- Grafana on `http://localhost:3000`

---

## Prerequisites

- API is running locally on `http://localhost:8080`
- You have `curl` installed
- Example user ID: `reviewer-1`
- Auth is mocked via `X-User-ID` header (must match the `{userId}` in the path)

---

## 1. Health & basic checks

### 1.1 Health

```bash
curl -i http://localhost:8080/health
```

Expected: `200 OK` with a small JSON payload like:

```json
{ "status": "healthy" }
```

---

## 2. Happy path – favourites for a single user

We’ll use a single user `reviewer-1` and create three favourites:

- one `chart`
- one `insight`
- one `audience`

### 2.1 Create a chart favourite

```bash
curl -i \
  -X POST \
  -H "Content-Type: application/json" \
  -H "X-User-ID: reviewer-1" \
  http://localhost:8080/v1/users/reviewer-1/favourites \
  -d '{
    "asset": {
      "type": "chart",
      "description": "Social media usage by age group",
      "chart": {
        "title": "Daily social media usage",
        "xAxis": "Age group",
        "yAxis": "Hours per day",
        "data": [1.2, 2.3, 3.8, 2.1]
      }
    }
  }'
```

Expected: `201 Created` and a response like:

```json
{
  "userId": "reviewer-1",
  "assetId": "ASSET_ID_CHART",
  "createdAt": "..."
}
```

> Note: copy the `assetId` value somewhere if you want to use it later.

### 2.2 Create an insight favourite

```bash
curl -i \
  -X POST \
  -H "Content-Type: application/json" \
  -H "X-User-ID: reviewer-1" \
  http://localhost:8080/v1/users/reviewer-1/favourites \
  -d '{
    "asset": {
      "type": "insight",
      "description": "Millennials social media insight",
      "insight": {
        "text": "40% of millennials spend more than 3 hours on social media daily"
      }
    }
  }'
```

Expected: another `201 Created` with a new `assetId`.

### 2.3 Create an audience favourite

```bash
curl -i \
  -X POST \
  -H "Content-Type: application/json" \
  -H "X-User-ID: reviewer-1" \
  http://localhost:8080/v1/users/reviewer-1/favourites \
  -d '{
    "asset": {
      "type": "audience",
      "description": "Males 24–35, heavy social media users",
      "audience": {
        "gender": "male",
        "birthCountry": "GB",
        "ageGroup": "24-35",
        "hoursOnSocialMediaPerDay": 3,
        "purchasesLastMonth": 3
      }
    }
  }'
```

Expected: `201 Created` with another `assetId`.

---

## 3. Listing & pagination

### 3.1 List all favourites (default pagination)

```bash
curl -s \
  -H "X-User-ID: reviewer-1" \
  "http://localhost:8080/v1/users/reviewer-1/favourites" | jq .
```

Expected:

- `items` array of favourites (3 items after the previous steps)
- `meta` object with `total` and `count`:

```json
{
  "items": [
    /* 3 items */
  ],
  "meta": {
    "total": 3,
    "count": 3
  }
}
```

### 3.2 List with pagination

Request first 2 items (limit=2, offset=0):

```bash
curl -s \
  -H "X-User-ID: reviewer-1" \
  "http://localhost:8080/v1/users/reviewer-1/favourites?limit=2&offset=0" | jq .
```

Request next 2 items (limit=2, offset=2):

```bash
curl -s \
  -H "X-User-ID: reviewer-1" \
  "http://localhost:8080/v1/users/reviewer-1/favourites?limit=2&offset=2" | jq .
```

> Note: You can omit `| jq .` if `jq` is not available on your system. The response will still be valid JSON, just not formatted.

Expected:

- First call: `meta.count` = 2
- Second call: `meta.count` = 1
- `meta.total` remains `3` in both responses.

### 3.3 List filtered by type

```bash
curl -s \
  -H "X-User-ID: reviewer-1" \
  "http://localhost:8080/v1/users/reviewer-1/favourites?type=insight" | jq .
```

Expected: only the `insight` favourite in `items`.

---

## 4. Edit description

Pick one `assetId` from the list responses (e.g. the chart one).

```bash
ASSET_ID="PUT_ASSET_ID_HERE"
```

Update the description:

```bash
curl -i \
  -X PATCH \
  -H "Content-Type: application/json" \
  -H "X-User-ID: reviewer-1" \
  "http://localhost:8080/v1/users/reviewer-1/favourites/$ASSET_ID" \
  -d '{
    "description": "Updated description from review scenario"
  }'
```

Expected: `200 OK` and a response item with the updated `description` field.

---

## 5. Remove a favourite (idempotent delete)

Reuse the same `ASSET_ID`:

```bash
curl -i \
  -X DELETE \
  -H "X-User-ID: reviewer-1" \
  "http://localhost:8080/v1/users/reviewer-1/favourites/$ASSET_ID"
```

Expected: `204 No Content`.

Deleting again is a no-op but should still return `204`:

```bash
curl -i \
  -X DELETE \
  -H "X-User-ID: reviewer-1" \
  "http://localhost:8080/v1/users/reviewer-1/favourites/$ASSET_ID"
```

---

## 6. Error examples (optional)

### 6.1 Invalid asset type

```bash
curl -i \
  -X POST \
  -H "Content-Type: application/json" \
  -H "X-User-ID: reviewer-1" \
  http://localhost:8080/v1/users/reviewer-1/favourites \
  -d '{
    "asset": {
      "type": "invalid-type",
      "description": "should fail"
    }
  }'
```

Expected: `422 Unprocessable Entity` with a Problem Details body describing the validation error.

### 6.2 Mismatched user ID (auth check)

```bash
curl -i \
  -H "X-User-ID: someone-else" \
  "http://localhost:8080/v1/users/reviewer-1/favourites"
```

Expected: `403 Forbidden` (user in header does not match `{userId}` in path).

---

## 7. Metrics

```bash
curl -s http://localhost:8080/metrics
```

Expected: Prometheus metrics including:

- `http_requests_total{method=...,path=...,status=...}`
- `http_request_duration_seconds_bucket{...}`

These will reflect the requests you just made in the scenario.
