# Platform Go Challenge API

Go-Based REST API for users, assets, and favourites with JWT authentication. Ships with Docker, Swagger UI, and seed data so you can try it immediately.
This application was made as a deliverble for GWI Engineering Challenge.
It covers all the basic requirements: 
- Web server 
- Endpoint that receives user Id and returns user's favorite assets
- Endpoints that would add an asset to favourites, remove it, or edit its description
- Data structure of the assets

Extra attributes:
- PostgreSQL database
- Dockerfile and docker-compose.yml
- unit tests
- swagger documentation
- JWT authentication (no user roles specified)

## Stack
- Go 1.25.6
- PostgreSQL 15
- Docker + docker-compose
- JWT via github.com/golang-jwt/jwt/v5
- Password hashing via golang.org/x/crypto/bcrypt
- Swagger UI via github.com/swaggo/http-swagger

## Quickstart (Docker)
```bash
cd docker
docker compose up --build
```
- API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html
- Adminer (DB UI): http://localhost:8081 (System=PostgreSQL, Server=platform-db, Username=user, Password=password as below)

## Authentication

The API uses **JWT (JSON Web Tokens)** for authentication. All protected endpoints require a valid JWT bearer token in the `Authorization` header.

### How JWT Authentication Works
1. **Register or Login**: Call `POST /register` or `POST /login` with credentials to receive a JWT token
2. **Include Token**: Add the token to subsequent requests using the header: `Authorization: Bearer <your_token>`
3. **Token Validation**: The API validates the token signature and expiration on each protected route
4. **Token Expiration**: Tokens expire after 24 hours by default

### Configuration
- **`JWT_SECRET`**: Environment variable for signing tokens (defaults to dev secret if not set)
- **`AUTH_DISABLED=true`**: Bypass authentication for local testing (not recommended for production)
- **Algorithm**: HS256 (HMAC with SHA-256)

### Seeded Users (for testing)
- **User ID**: `u1` / **Password**: `alice123`
- **User ID**: `u2` / **Password**: `bob123`

### Authentication Endpoints
- **POST /register** — Create new user (requires: `id`, `name`, `password`)
- **POST /login** — Authenticate and get JWT token (requires: `id`, `password`)

### Example Usage
```bash
# 1. Login to get token
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"id":"u1","password":"alice123"}'

# Response: {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}

# 2. Use token for protected endpoints
curl -X GET http://localhost:8080/users/u1/favourites \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Protected Routes
All `/users/*` endpoints require valid JWT authentication (unless `AUTH_DISABLED=true`).

## Endpoints (summary)
- Auth: POST /login, POST /register
- Health: GET /health
- Users: GET /users, POST /users
- Favourites: GET /users/{userId}/favourites, POST /users/{userId}/favourites, PATCH /users/{userId}/favourites/{assetId}, DELETE /users/{userId}/favourites/{assetId}
- Assets: GET /assets/{id}, POST /assets

See full request/response schemas in Swagger UI.
http://localhost:8080/swagger/index.html#/

## Database seeding
`db/init/001_init.sql` creates tables and seeds:
- Users u1/u2 (with bcrypt password hashes)
- Sample assets (insight, chart)
- Sample favourites linking users to assets

## Running tests
```bash
go test ./...
```

## Configuration
- `DATABASE_URL` (required): postgres connection string.
- `PORT` (default 8080): HTTP port inside the container/process.
- `JWT_SECRET`: secret for signing tokens.
- `AUTH_DISABLED`: set to `true` to bypass auth checks (use only for local testing).

## Useful URLs (Docker defaults)
- API base: http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html
- Adminer: http://localhost:8081 (server=db, user=user, password=password, database=dashboard)
