# GlobalWebIndex Engineering Challenge

This is a simple Go application with middleware for JWT authentication, rate limiting, and basic request handling. The application is containerized using Docker.

## Project Structure

- `main.go`: Entry point of the application.
- `asset.go`: Handler for serving static assets.
- `auth_handler.go`: Handler for authenticated routes.
- `asset_handler.go`: Handler for asset management.
- `jwt_verifier.go`: Middleware for JWT authentication.
- `rate_limiter.go`: Middleware for rate limiting.
- `jwt_verifier_test.go`: Tests for JWT middleware.
- `rate_limiter_test.go`: Tests for rate limiting middleware.
- `asset_handler_test.go`: Tests for basic asset handlers.

## Prerequisites

- Docker installed on your machine.
- Go 1.22 or later (if you plan to run the application locally without Docker).

## Getting Started

### Running with Docker

1. **Build the Docker image:**

    ```bash
    docker build -t go-app .
    ```

2. **Run the Docker container:**

    ```bash
    docker run -p 8080:8080 go-app
    ```

3. The application should now be running on `http://localhost:8080`.

### Running Locally

1. **Clone the repository:**

    ```bash
    git clone <repository-url>
    cd <repository-directory>
    ```

2. **Install dependencies:**

    ```bash
    go mod download
    ```

3. **Run the application:**

    ```bash
    go run main.go
    ```

4. The application should now be running on `http://localhost:8080`.

## Endpoints

- `POST /login`: Authenticates a user and returns a JWT token.
- `GET /favorites/{userID}`: Retrieves a list of favorite assets for a user.
- `POST /favorites/{userID}`: Adds a new favorite asset for a user.
- `DELETE /favorites/{userID}/{assetID}`: Removes a favorite asset for a user.
- `PUT /favorites/{userID}/{assetID}`: Edits a favorite asset for a user.

## Middleware

- **JWT Authentication**: Protects routes by requiring a valid JWT token.
- **Rate Limiting**: Limits the number of requests to prevent abuse.

## Running Tests

To run the tests, use the following command:

```bash
go test ./...
