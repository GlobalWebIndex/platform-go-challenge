# GWI Favourites Service

A small Go service that lets users save their favourite assets (charts, insights, audiences). It connects to PostgreSQL for storage and is meant to be run with Docker Compose. It is has been developed and verified that it is working in Linux with Go v1.25.6.

## How it works

There are two main concepts:

- **Assets** — these are platform entities like Charts, Insights, and Audiences.
- **Favourites** — a user picks an asset they like. When users fetch their favourites, each one comes back with its full asset data and a description. The user is able to add and change a personalized description. The user can remove an asset from favourites.

### Convention to simplify the demonstration of this project

The task assumes that the assets are stored in a pre-existing database table:
assets(asset_id PK, asset_type, description, data JSONB)

The user picks favourites in the dashboard and they are related in the database:
favourites(user_id, asset_id, description, PK(user_id, asset_id))

Setting up a pre-existing assets table every time the Favourites Service runs would render the demonstration of this project much more complex. It would add overhead, because a database setup script would have to be executed after Postgres is healhty and running with docker compose. To make things as simple as possible for the demonstration, I omit this step and let the POST operation add an asset and also mark it as favourite for the user.

Asset IDs also are not UUIDs in this implementation. A simple check that a user does not have duplicate favourite IDs is implemented.

## API

Every request needs a JWT token in the `Authorization: Bearer <token>` header. The user ID is pulled from the token's `sub` claim — there's no user ID in the URL.

### Required Headers

| Header | Required For | Value |
|--------|--------------|-------|
| `Authorization` | All requests | `Bearer <token>` |
| `Accept` | All requests | Must include `application/json` (or `*/*`) |
| `Content-Type` | POST, PATCH | Must be `application/json` |

Missing or invalid headers result in:
- **401 Unauthorized** — missing or invalid JWT token
- **406 Not Acceptable** — missing or invalid `Accept` header
- **415 Unsupported Media Type** — missing or invalid `Content-Type` on requests with a body

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/favourites` | Get all favourites for the authenticated user |
| `POST` | `/api/v1/favourites` | Add a new favourite |
| `PATCH` | `/api/v1/favourites/{asset_id}` | Update a favourite's description |
| `DELETE` | `/api/v1/favourites/{asset_id}` | Remove a favourite |
| `GET` | `/health/ready` | Health check (served on a separate port, intended for deployment only) |
| `GET` | `/health/live` | Health check (served on a separate port, intended for deployment only) |

Here's what the request/response bodies look like:

**Adding a favourite (POST):**

Chart asset:

```json
{
  "asset_type": "chart",
  "asset_data": {
    "id": "chart-001",
    "title": "Monthly Revenue",
    "x_axis_title": "Month",
    "y_axis_title": "Revenue (USD)",
    "data": {
      "January": 12000,
      "February": 15000,
      "March": 18000
    }
  },
  "description": "Revenue trend for Q1 2026"
}
```

Audience asset:

```json
{
  "asset_type": "audience",
  "description": "Tech-savvy millennials",
  "asset_data": {
    "id": "audience-123",
    "gender": ["Male", "Female"],
    "birth_country": ["US", "UK"],
    "age_groups": ["25-34"],
    "social_media_hours_daily": "3-5",
    "purchases_last_month": 5
  }
}
```
Insight asset:

```json
{
  "asset_type": "insight",
  "asset_data": {
    "id": "insight-001",
    "text": "Users who engage with social media 3-5 hours daily have a 40% higher rate"
  },
  "description": "Social media engagement insight"
}
```

**Updating a description (PATCH):**
```json
{ "description": "Updated description" }
```

**Listing favourites (GET) — returns something similar to:**
```json
[
  {
    "asset_id": "chart-1",
    "type": "chart",
    "description": "My favourite chart",
    "data": { "title": "Social Media Usage", "x_axis": "Age Group", "y_axis": "Hours", "points": [...] }
  }
]
```

**Remove a favourite:
DELETE /api/v1/favourites/chart-1

There is also a full OpenAPI spec in `api/swagger.yaml`.

## Configuration

The app reads port settings from `config.yaml` and/or environment variables (env vars win if both are set). Database and JWT settings only come from environment variables.

| Setting | Env Variable | Config File Key | Default |
|---------|-------------|-----------------|---------|
| API port | `API_PORT` | `api_port` | `8000` |
| Health port | `HEALTH_PORT` | `health_port` | `8001` |
| DB host | `POSTGRES_HOST` | — | `postgres` (in Compose) |
| DB port | `POSTGRES_PORT` | — | `5432` |
| DB host port | `POSTGRES_HOST_PORT` | — | `5432` (change in case of host port conflict) |
| DB user | `POSTGRES_USER` | — | — |
| DB password | `POSTGRES_PASSWORD` | — | — |
| DB name | `POSTGRES_DB` | — | — |
| JWT secret | `JWT_SECRET` | — | empty |
| Allow unsigned tokens | `ALLOW_UNSIGNED_TOKENS` | — | `false` |

You can point to a different config file by setting the `CONFIG_PATH` env var.

## How to run the service

The service needs PostgreSQL, so Docker Compose is required for local running and testing. A deployment.yaml is not included for Kubernetes support, however this project is designed for a straightforward deployment to Kubernetes as a next step.

**1. Set up your `.env` file:**

```bash
cp .env.example .env
```

Open `.env` and fill in at least `POSTGRES_USER` and `POSTGRES_PASSWORD`. Set `ALLOW_UNSIGNED_TOKENS=true`, otherwise e2e tests will fail and requests will not be authorized.

**2. Start everything:**

```bash
docker compose up --build
```

This builds the Go binary, spins up Postgres, waits for it to be healthy, and then starts the app. Add `-d` to run it in the background.

**3. When you are done:**

```bash
docker compose down
```

If you want a clean state (wipe the database too):

```bash
docker compose down --volumes
```

## Testing

### Unit tests

```bash
go test ./internal/... -v
```

The unit tests mock the Postgres database, so Docker Compose is not required.

### End-to-end tests

The E2E tests hit the real running services, so you need Docker Compose up first and ensure that the .env file has `ALLOW_UNSIGNED_TOKENS=true`:

```bash
docker compose up --build
```

Once both containers are healthy, open another terminal:

```bash
cd e2e_tests
go clean -testcache
go test -v -count=1
```

After testing is complete, run:

```bash
cd ..
docker compose down --volumes
```

## Authentication

The API uses JWTs. There is a tool included to generate tokens:

```bash
go run ./tools/tokengen -user alice
```

Or with a JWT secret, which should also be included in the .env file with JWT_SECRET={SECRET}.

```bash
go run ./tools/tokengen -user alice -secret {SECRET}
```

Use the token with curl or Postman:

```bash
curl -H "Authorization: Bearer <token>" \
     -H "Accept: application/json" \
     http://localhost:8000/api/v1/favourites
```

### Token Modes

The service supports two authentication modes:

| Mode | Configuration | Use Case |
|------|---------------|----------|
| **Signed tokens** | Set `JWT_SECRET` | Production — tokens must be HS256-signed with the secret |
| **Unsigned tokens** | No `JWT_SECRET` + `ALLOW_UNSIGNED_TOKENS=true` | Local development and testing only |

**Important:** Unsigned tokens (`alg=none`) require explicit opt-in via `ALLOW_UNSIGNED_TOKENS=true`. This is a safety measure — if there is a failure to set `JWT_SECRET` in production but don't set `ALLOW_UNSIGNED_TOKENS`, all requests will be rejected.

The Docker Compose setup defaults to `ALLOW_UNSIGNED_TOKENS=true` for easy local development. For production, always set up Kubernetes to fetch a proper `JWT_SECRET` and leave `ALLOW_UNSIGNED_TOKENS` unset or `false`.

## Storage

Favourites are stored in PostgreSQL. The table uses a composite primary key `(user_id, asset_id)` and keeps the polymorphic asset data in a `jsonb` column. The schema creates itself on startup with (`CREATE TABLE IF NOT EXISTS`).

A few things I would consider for production:

- **Caching** — implement Cache-Control and ETag.
- **Pipeline code** - Github Actions, Jenkins etc.
- **Audit Logging**
- **Observability** - OpenTelemetry instrumentation and Grafana 
