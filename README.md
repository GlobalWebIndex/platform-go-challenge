# User Favorites API

This is a sample Go application that provides an API for managing user favorites. Users can register, login, and manage their favorite assets (charts, insights, and audiences). The API uses JWT for authentication.

## Project Structure

```
user-favorites-api/
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── handlers/
│   └── handlers.go
├── main.go
├── models/
│   └── models.go
├── router/
│   └── router.go
├── store/
│   └── store.go
├── middleware/
│   └── jwt.go
├── db_scripts/
│   └── dbsetup.sql
└── README.md
```

## Setup

### Prerequisites

- Docker
- Docker Compose

### Running the Application

1. Ensure you have Docker and Docker Compose installed.
2. Build and start the application using Docker Compose:

```sh
docker-compose up --build --remove-orphans
```

This command will build the Go application image, start the PostgreSQL database, and then start the Go application. The application will be available at `http://localhost:8080`, and the PostgreSQL database will be available at `localhost:5432` with the credentials specified.

## API Endpoints

### User Registration

- **URL:** `/register`
- **Method:** `POST`
- **Request Body:**

```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "password123"
}
```

- **Sample Request:**

```sh
curl -X POST http://localhost:8080/register -d '{"username":"john_doe","email":"john@example.com","password":"password123"}' -H "Content-Type: application/json"
```

### User Login

- **URL:** `/login`
- **Method:** `POST`
- **Request Body:**

```json
{
  "username": "john_doe",
  "password": "password123"
}
```

- **Sample Request:**

```sh
curl -X POST http://localhost:8080/login -d '{"username":"john_doe","password":"password123"}' -H "Content-Type: application/json"
```

- **Response:**

```json
{
  "token": "your_jwt_token_here"
}
```

### Get User Favorites

- **URL:** `/api/favorites/{userID}`
- **Method:** `GET`
- **Headers:** `Authorization: Bearer your_jwt_token_here`
- **Sample Request:**

```sh
curl -H "Authorization: Bearer your_jwt_token_here" http://localhost:8080/api/favorites/1
```

### Add Favorite

- **URL:** `/api/favorites/{userID}`
- **Method:** `POST`
- **Headers:** `Authorization: Bearer your_jwt_token_here`
- **Request Body:**

```json
{
  "id": "1",
  "type": "chart",
  "description": "Test Chart",
  "data": "Sample Data"
}
```

- **Sample Request:**

```sh
curl -X POST -H "Authorization: Bearer your_jwt_token_here" -d '{"id":"1","type":"chart","description":"Test Chart","data":"Sample Data"}' -H "Content-Type: application/json" http://localhost:8080/api/favorites/1
```

### Remove Favorite

- **URL:** `/api/favorites/{userID}/{assetID}`
- **Method:** `DELETE`
- **Headers:** `Authorization: Bearer your_jwt_token_here`
- **Sample Request:**

```sh
curl -X DELETE -H "Authorization: Bearer your_jwt_token_here" http://localhost:8080/api/favorites/1/1
```

### Edit Favorite

- **URL:** `/api/favorites/{userID}/{assetID}`
- **Method:** `PUT`
- **Headers:** `Authorization: Bearer your_jwt_token_here`
- **Request Body:**

```json
{
  "description": "Updated Chart Description",
  "data": "Updated Sample Data"
}
```

- **Sample Request:**

```sh
curl -X PUT -H "Authorization: Bearer your_jwt_token_here" -d '{"description":"Updated Chart Description","data":"Updated Sample Data"}' -H "Content-Type: application/json" http://localhost:8080/api/favorites/1/1
```

## Environment Variables

- `DATABASE_URL`: The URL of the PostgreSQL database. Example: `postgresql://user_favorites_user:mysecretpassword@db:5432/user_favorites_db`
- `SECRET_KEY`: The secret key used for signing JWT tokens. It should be a long and random string. Example: `my_super_secret_key_123`
