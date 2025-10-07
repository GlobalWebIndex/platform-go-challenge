# GWI platform-go-challenge (Go)

A small Go service for managing **per-user favourites** of polymorphic assets (Chart, Insight, Audience).  
Built for the GWI engineering challenge.

- Language: Go (net/http, chi)
- Storage: Postgres (JSONB for polymorphic `data`)
- Extras: pagination, graceful shutdown, Postman collection

---

## Features

- Users can **add/list/edit/remove** favourites.
- Assets are **polymorphic** via `type` + JSONB `data`.
- **Pagination** on listing favourites.
- Clean architecture: `http → service → repo`.
- Includes demo user in migration (gwi@demo.local).

---

## Project Structure

- cmd/server       # main entrypoint
- internal/config  # env config
- internal/http    # router, handlers
- internal/service # business logic
- internal/repo    # postgres repos
- migrations/      # schema + demo user

## Database Schema (Postgres)

Tables:
- `users(id, email)`
- `assets(id, type ENUM('chart','insight','audience'), description, data JSONB, created_at)`
- `favourites(user_id, asset_id, custom_description, created_at, PK(user_id, asset_id))`

Indexes:
- `idx_assets_type` on `assets(type)`
- `idx_favourites_user` on `favourites(user_id)`

> Schema SQL is in migrations/001_init.sql (creates tables + inserts demo user).

---

## How to run

### Dockerfile

First build image: docker build -t platform-go-challenge:latest .

Then run container (with database url):  docker run --rm -p 8080:8080 -e DATABASE_URL="postgres://gwi_user:gwi_pass@host.docker.internal:5432/gwi?sslmode=disable" platform-go-challenge:latest

### Local Dev
$env:PORT="8080"
$env:DATABASE_URL="postgres://gwi_user:gwi_pass@localhost:5432/gwi?sslmode=disable"
go run .\cmd\server

## Postman Collection

Import GWI_Favourites_API.postman_collection.json.
It contains ready-made requests for:

- POST /users (create user)
- GET /users?limit=&offset= (list all users, with pagination)
- GET /users/{id} (get one user)
- GET /users/{id}/favourites?limit=&offset= (list favourites for user)
- POST /users/{id}/favourites (add favourite: Insight / Chart / Audience)
- PATCH /users/{id}/favourites/{assetId} (edit description)
- DELETE /users/{id}/favourites/{assetId} (delete one favourite)
- DELETE /users/{id}/favourites (clear all users favourites)
