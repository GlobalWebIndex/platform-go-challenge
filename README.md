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
- `GET /users/{id}` → list favorites
- `GET /users/{id}/{assetId}` → get a single favorite asset by id
- `POST /users/{id}` → add favorite
- `PUT /users/{id}/{assetId}` → edit favorite
- `DELETE /users/{id}/{assetId}` → remove favorite

## Run with Docker
```bash
docker build -t platform-go-challenge .
docker run -p 8080:8080 platform-go-challenge
```

## Initial curl commands to test GET, POST, PUT and DELETE apis

### Create an Insight asset
```bash
curl -X POST http://localhost:8080/users/user123 \
-H "Content-Type: application/json" \
-d '{
"id": "12",
"type": "Insight",
"description": "New insight asset"
}'
```
### Retrieve all assets of a user
```bash
curl -X GET http://localhost:8080/users/user123
```

### Update an existing asset
```bash
curl -X PUT http://localhost:8080/users/user123/12 \
  -H "Content-Type: application/json" \
  -d '{
        "description": "Updated insight description"
      }'
```


### Delete an existing asset
```bash
curl -X DELETE http://localhost:8080/users/user123/12
```