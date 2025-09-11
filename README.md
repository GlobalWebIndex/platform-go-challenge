# Favorites Service

A small Go web service that lets users manage their favorite assets (Charts, Insights, Audiences).

## Features
- Add, remove, update, and list favorites
- In-memory store (thread-safe)
- Basic tests included
- Dockerized

## Running locally
```bash
go run main.go
```

## API Endpoints
- `GET /users/{id}/favorites` → list favorites
- `POST /users/{id}/favorites` → add favorite
- `PUT /users/{id}/favorites/{assetId}` → edit favorite
- `DELETE /users/{id}/favorites/{assetId}` → remove favorite

## Run with Docker
```bash
docker build -t favorites-service .
docker run -p 8080:8080 favorites-service
```