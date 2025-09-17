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

### Create a Chart asset
```bash
curl -X POST http://localhost:8080/users/user123 \
-H "Content-Type: application/json" \
-d '{
  "id": "12",
  "type": "Chart",
  "description": "Monthly sales chart",
  "data": {
    "title": "Sales",
    "x_axis": "Month",
    "y_axis": "Revenue",
    "data": [100, 200, 300]
  }
}'
```

### Create an Insight asset
```bash
curl -X POST http://localhost:8080/users/user123 \
-H "Content-Type: application/json" \
-d '{
  "id": "13",
  "type": "Insight",
  "description": "Social media insight",
  "data": {
    "text": "How many hours people spend on their phones"
  }
}'
```

### Create an Audience asset
```bash
curl -X POST http://localhost:8080/users/user123 \
-H "Content-Type: application/json" \
-d '{
  "id": "14",
  "type": "Audience",
  "description": "Male audience",
  "data": {
    "gender": "Male",
    "birth_country": "Greece",
    "age_group": "24-35",
    "social_hours": "3",
    "purchases_last": "5"
  }
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