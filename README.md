# Platform Go Challenge

This is a compact, interview‑style implementation that focuses on clarity, correctness and testing rather than infrastructure heft. I chose an in‑memory design for now to keep iteration fast and to showcase handler logic, auth flows, and concurrency safety. In a real‑world system I would swap the stores for a persistent backend and introduce stronger operational guarantees (revocation, observability, migrations).

## Prerequisites
- Go 1.24+: build and run (`go run .`; `go test ./...`).
- `curl`: used for terminal examples.
- `jq` (Bash/Zsh only): parses JSON in examples.
- Bash/Zsh or PowerShell: Bash/Zsh is primary; PowerShell alternatives provided.
- Docker (optional): container builds and `docker compose` workflows.

## Table of Contents
1. Architecture Overview
2. Data Models
3. Authentication & Authorization
4. Configuration
5. Endpoints & Usage
6. Example Workflows
7. Running the Service
8. Testing
9. Design Decisions & Limitations
10. Future Improvements

---
## 1. Architecture Overview

| Layer | Responsibility |
|-------|---------------|
| `internal/storage` | Thread‑safe in‑memory stores (`UserStore`, `FavoritesStore`) using `sync.RWMutex`. |
| `internal/models` | Domain models: `User`, `Chart`, `Insight`, `Audience` + creation helpers. |
| `internal/auth` | JWT generation (`RealJWTGenerator`) + middleware (`TokenMiddleware`, `AdminAPIKeyMiddleware`). |
| `internal/handlers` | HTTP handlers for user registration/login, asset favorites CRUD, admin user operations. |
| `internal/config` | Central configuration with environment overrides. |
| `main.go` | Service bootstrap (router setup, HTTP server). |

Key choices and why:
- **In‑memory storage**: Keeps the footprint minimal and speeds up testing. I used `sync.RWMutex` to demonstrate thread‑safe access and added targeted concurrency tests. In production I’d introduce a repository layer backed by PostgreSQL or SQLite and keep the handler/storage interfaces unchanged.
- **Favorites ownership**: Normal users can only see/touch their own data; admins can optionally pass `user_id` to act on behalf of a user. This keeps authorization logic explicit at the boundary.
- **Polymorphic assets**: A lightweight registry instantiates `chart|insight|audience` models from `asset.type`. It’s pragmatic here; with more time I’d formalize this via generics or well‑defined protobuf/JSON schemas.
- **JWT + Admin API key**: JWT for identity/role on `/api/*`, and an API key on `/admin/*` to make high‑privilege actions very obvious and separable. In production I might collapse these into RBAC via JWT + policy, or use mTLS/service accounts for admin surfaces.
- **Stateless tokens**: Simpler and sufficient here. In a real system I’d add a token blacklist (e.g., Redis) or short‑lived tokens + refresh to immediately invalidate access when users are deleted or permissions change.

---
## 2. Data Models

### User
```json
{
	"id": "<uuid>",
	"username": "alice",
	"role": "admin" | "user"
}
```
Roles validated via `models.IsValidRole` (`admin`, `user`). Passwords are hashed with bcrypt during registration; the stores never handle raw passwords.

### Asset (Favorites Entry)
Each favorite wraps an asset with a type discriminator.
Types supported: `chart`, `insight`, `audience`.

Generic response shape (handler wraps asset with its type for clients):
```json
{
	"type": "chart",
	"asset": {
		"id": "chart_1",
		"title": "Revenue Q1",
		"description": "Quarterly revenue",
		"data": [{"label": "Q1", "value": 1500}],
		"created_at": "2025-12-01T10:15:30Z",
		"updated_at": "2025-12-01T10:15:30Z"
	}
}
```

### Favorites Listing
```json
{
	"assets": [ { "type": "chart", "asset": { ... } }, ... ]
}
```

---
## 3. Authentication & Authorization

| Mechanism | Used For | Header |
|-----------|----------|--------|
| JWT (HS256) | `/api/*` protected routes | `Authorization: Bearer <token>` |
| Admin API Key | `/admin/*` routes | `X-Admin-API-Key: <key>` |

Flow:
1. Client `POST /register` → obtains user record.
2. Client `POST /login` → receives JWT token (`{"token": "..."}`).
3. Use token on `/api/favorites/*` routes; include optional `user_id` when admin wants to act on another user.
4. Admin routes require the API key only; they do not require a JWT, but JWT can still be supplied for downstream mixed calls.

Claims include `user_id` and registered claims (expiry ≈ 1 hour, configurable). Middleware consistently injects `userID` and `role` into request context.

---
## 4. Configuration

Environment variables (with defaults):
| Var | Default | Purpose |
|-----|---------|---------|
| `JWT_SECRET` | `supersecretkey` | HMAC secret for token signing. |
| `ADMIN_API_KEY` | `supersecret-admin-key` | Required for `/admin/*`. |
| `SERVER_PORT` | `:8080` | Port binding when starting the server. |

Set vars before running (Bash/Zsh):
```bash
export JWT_SECRET="change-me"
export ADMIN_API_KEY="my-admin-key"
export SERVER_PORT=":9090"
go run .
```

<details>
<summary>PowerShell alternative</summary>

```powershell
$env:JWT_SECRET = "change-me";
$env:ADMIN_API_KEY = "my-admin-key";
$env:SERVER_PORT = ":9090";
go run .
```

</details>

---
## 5. Endpoints & Usage

### Public

#### POST `/register`
Register a new user.
Request:
```json
{ "username": "alice", "password": "s3cur3", "role": "user" }
```
Response `201 Created`:
```json
{ "id": "<uuid>", "username": "alice", "role": "user" }
```
Errors: `400` invalid role / empty fields, `409` duplicate username.

#### POST `/login`
Authenticate and retrieve JWT.
Request:
```json
{ "username": "alice", "password": "s3cur3" }
```
Response `200 OK`:
```json
{ "token": "<jwt>" }
```
Errors: `401` wrong credentials, `400` malformed JSON.

### Authenticated (JWT required)

Common Header: `Authorization: Bearer <jwt>`

#### GET `/api/favorites/list` (optional `user_id` query param)
List favorites for current user or specified user (admin only).
Example: `/api/favorites/list?user_id=<uuid>`
Response:
```json
{ "assets": [ { "type": "chart", "asset": { ... } } ] }
```
Errors: `403` forbidden (non-admin requesting other user), `500` internal on storage failure.

#### POST `/api/favorites/add`
Add a favorite; admin may supply `user_id` to add for another user.
Request:
```json
{
	"user_id": "<optional uuid>",
	"asset": {
		"type": "chart",
		"id": "chart_1",
		"title": "Revenue Q1",
		"description": "Quarterly revenue",
		"data": [{"label":"Q1","value":1500}]
	}
}
```
Response `201 Created` (echoes asset or minimal success message).
Errors: `403` cross-user without admin, `409` duplicate (currently returns `500` – improvement noted).

#### DELETE `/api/favorites/remove`
Removes a favorite by `asset_id`; admin can target another user using `user_id` query param.
Example: `/api/favorites/remove?asset_id=chart_1` or `/api/favorites/remove?user_id=<uuid>&asset_id=chart_1`
Response `200 OK` on success.
Errors: `404` asset not found, `403` forbidden.

#### PATCH `/api/favorites/update`
Update existing favorite asset fields.
Request:
```json
{
	"asset": {
		"type": "chart",
		"id": "chart_1",
		"title": "Updated Title",
		"description": "New description",
		"data": [{"label":"Q1","value":1600}]
	}
}
```
Response `200 OK` with updated asset.
Errors: `404` when asset missing (currently may surface as `500` in this simplified handler; noted below), `403` forbidden cross-user.

### Admin (API Key required)

Header: `X-Admin-API-Key: <key>`

#### GET `/admin/users/list`
Returns JSON object mapping usernames to user IDs:
```json
{ "alice": "<uuid>", "bob": "<uuid>" }
```

#### DELETE `/admin/users/{id}`
Deletes a user and cleans associated favorites.
Response: `204 No Content`
Errors: `400` invalid UUID, `404` user not found.

---
## 6. Example Workflows

### Register + Login + Add Favorite (Bash/Zsh)
Requires `jq` for parsing JSON.
```bash
# Register
curl -sS -X POST http://localhost:8080/register \
	-H 'Content-Type: application/json' \
	-d '{"username":"alice","password":"pw","role":"user"}'

# Login
TOKEN=$(curl -sS -X POST http://localhost:8080/login \
	-H 'Content-Type: application/json' \
	-d '{"username":"alice","password":"pw"}' | jq -r '.token')

# Add favorite
curl -sS -X POST http://localhost:8080/api/favorites/add \
	-H "Authorization: Bearer $TOKEN" \
	-H 'Content-Type: application/json' \
	-d '{"asset":{"type":"chart","id":"c1","title":"Chart","description":"Desc","data":[{"label":"Q1","value":100}]}}'

# List favorites
curl -sS -H "Authorization: Bearer $TOKEN" \
	http://localhost:8080/api/favorites/list
```

<details>
<summary>PowerShell alternative</summary>

```powershell
# Register
curl -X POST http://localhost:8080/register -H "Content-Type: application/json" -d '{"username":"alice","password":"pw","role":"user"}'
# Login
$token = (curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{"username":"alice","password":"pw"}' | ConvertFrom-Json).token
# Add favorite
curl -X POST http://localhost:8080/api/favorites/add -H "Authorization: Bearer $token" -H "Content-Type: application/json" -d '{"asset":{"type":"chart","id":"c1","title":"Chart","description":"Desc","data":[{"label":"Q1","value":100}]}}'
# List favorites
curl -X GET http://localhost:8080/api/favorites/list -H "Authorization: Bearer $token"
```

</details>

### Admin List & Delete User (Bash/Zsh)
```bash
export ADMIN_API_KEY="supersecret-admin-key"
curl -sS -H "X-Admin-API-Key: $ADMIN_API_KEY" \
	http://localhost:8080/admin/users/list

curl -sS -X DELETE -H "X-Admin-API-Key: $ADMIN_API_KEY" \
	http://localhost:8080/admin/users/<uuid>
```

<details>
<summary>PowerShell alternative</summary>

```powershell
$env:ADMIN_API_KEY = "supersecret-admin-key"
curl -H "X-Admin-API-Key: $env:ADMIN_API_KEY" http://localhost:8080/admin/users/list
curl -X DELETE -H "X-Admin-API-Key: $env:ADMIN_API_KEY" http://localhost:8080/admin/users/<uuid>
```

</details>

Tip: Install jq if missing: macOS `brew install jq`, Debian/Ubuntu `sudo apt-get install -y jq`.

---
## 7. Running the Service

### Local (Go)
```bash
go run .
```
Server listens on `SERVER_PORT` (default `:8080`).

### Docker
```bash
docker build -t platform-go-challenge .
docker run -p 8080:8080 -e JWT_SECRET=change-me -e ADMIN_API_KEY=my-admin-key platform-go-challenge
```

### docker-compose
```bash
docker compose up --build
```

---
## 8. Testing
Run unit + integration + advanced scenarios:
```bash
go test ./... -v
```
Key tests:
- `internal/handlers/e2e_test.go` basic flow.
- `internal/handlers/e2e_advanced_test.go` concurrency, permissions, edge cases.
- `internal/storage/asset_store_test.go` concurrency safety.

---
## 9. Design Decisions & Limitations (intentionally candid)

| Topic | Notes |
|-------|-------|
| In‑memory only | Chosen for a take‑home to keep focus on code quality and tests; not persistent. I would switch to Postgres/SQLite with migrations and transactional semantics. |
| Stateless JWT | Fast and simple. I would add revocation (Redis list) or shorter TTL + refresh tokens in production. |
| Error codes | Some edges bubble up as `500`; I would add a small error mapping layer to normalize domain → HTTP (`404`, `409`, etc.). |
| Admin surface | Kept as API key to clearly separate privileged controls. In production, RBAC via JWT claims or an IAM‑style policy engine is preferable. |
| Concurrency | Demonstrated with tests; in production I would add metrics and tracing to observe lock contention and latencies. |
| Asset types | Lightweight and flexible now; I would define stricter schemas and validation, plus pagination for large lists. |

---
## 10. Future Improvements (what I would do next)
- Swap stores for Postgres or SQLite behind the current interfaces and add migrations.
- Introduce token revocation / blacklist (Redis) or shorter JWT TTL + refresh.
- Normalize errors with a tiny mapper to ensure `404/409/403` where appropriate.
- Add pagination + filters to favorites.
- Add rate limiting and request logging; wire tracing (OTel) for observability.
- Generate an OpenAPI spec and publish it with a Swagger UI.
- CI pipeline with coverage gates and static analysis.

---
## Quick Reference

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | /register | none | Create user |
| POST | /login | none | Obtain JWT |
| GET | /api/favorites/list | JWT | List favorites (self or other if admin) |
| POST | /api/favorites/add | JWT | Add favorite asset |
| DELETE | /api/favorites/remove | JWT | Remove favorite asset |
| PATCH | /api/favorites/update | JWT | Update favorite asset |
| GET | /admin/users/list | API Key | List all users |
| DELETE | /admin/users/{id} | API Key | Delete user |

---

