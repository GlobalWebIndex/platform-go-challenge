# GlobalWebIndex Favorites Service

A Go web service that allows users to manage their list of favorite assets (Charts, Insights, Audiences).
The service supports:

- Fetching all favorites for a user
- Adding a favorite
- Removing a favorite
- Editing an asset’s description
- JWT authentication
- OpenAPI documentation
- Docker support

This project was built as part of the GlobalWebIndex Engineering Challenge.

## Asset Types

### Chart

- id: integer
- title: string
- x_axis: string
- y_axis: string
- data: array of points (each point has x and y values)
- description: string

### Insight

- id: integer
- text: string
- description: string

### Audience

- id: integer
- gender: string
- birth_country: string
- age_group: string
- hours_social_media_daily: integer
- number_of_purchases_last_month: integer
- description: string

## Favorite Management

Users can add, remove, and edit their favorite assets. Each user is identified by a unique user ID and must authenticate using JWT tokens.

- GET /favorites: Retrieve all favorite assets for the authenticated user.
- POST /favorites: Add a new favorite asset.
- DELETE /favorites: Remove a favorite asset.
- PATCH /favorites: Edit the description of a favorite asset.

## Storage

For the challenge, storage is implemented using a JSON file acting as a simple local database.

## Authentication

The service uses JWT tokens for authentication. Each request to the favorites endpoints must include a valid JWT token in the Authorization header.

## Directory Structure

- `controllers/`: Contains the HTTP handlers for the API endpoints.
- `customerrors/`: Defines custom error types for better error handling.
- `middleware/`: Contains middleware for authentication
- `models/`: Defines the data models for assets, favorites, and users.
- `openapi/`: OpenAPI documentation files.
- `routes/`: Defines the API routes and associates them with controllers and middleware.
- `services/`: Business logic for managing favorites and assets.
- `storage/`: Handles data storage and retrieval from the JSON file.
- `utils/`: Utility functions including JWT token generation and validation.
- `test/`: Contains unit tests for the service.
- `validators/`: Input validation logic for API requests.
- `data.json`: Sample data file acting as a local database.
- `Dockerfile`: Docker configuration for containerizing the application.
- `go.mod` and `go.sum`: Go module files for dependency management.
- `main.go`: Entry point of the application.
- `test_data.json`: Sample data file used for testing.

## Environment Variables

The service requires environment variables for configuration. Create a `.env` file in the root directory with the following variables:

- `JWT_SECRET`: Secret key used for signing JWT tokens
- `DATA_FILE`: Path to the JSON file used for data storage (`data.json`)

**Note**: The `.env` file is not included in the repository for security reasons. Make sure to create your own `.env` file with appropriate values before running the service.

## Running the Service

1. Ensure you have Go installed on your machine.
2. Clone the repository.
3. Navigate to the project directory.
4. **Create a `.env` file with required environment variables (see Environment Variables section above).**
5. Run `go mod tidy` to install dependencies.
6. Start the service by running `go run main.go`.
7. The service will be available at `http://localhost:8090`.

### Running with Docker

1. Ensure you have Docker installed on your machine.
2. Build the Docker image:
   ```bash
   docker build -t gwx-favorites-service .
   ```
3. Run the Docker container:
   ```bash
   docker run -p 8090:8090 gwx-favorites-service
   ```

## API Documentation

The API documentation is available via OpenAPI and can be accessed at `http://localhost:8090/docs/`.

## Testing

The service includes comprehensive unit tests to ensure reliability and correctness of the API endpoints and business logic.

### Running Tests

To run the tests, execute the following command in the project directory:
`bash
    go test ./test/
    `

## Test structure

- `favorites_test.go`: Contains tests for the favorites endpoints, including authentication scenarios.
- `auth_test.go`: Contains tests for authentication mechanisms.
