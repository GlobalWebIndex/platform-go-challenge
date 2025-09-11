# GWI Favourites API

This document provides a comprehensive overview of the GWI Favourites API, covering how to build and run the application with Docker, its core architecture, and the testing principles used in development.

---

### 🚀 Getting Started

The project uses a multi-stage Docker build to create a secure and highly efficient production image. The first stage builds the Go binary, and the second copies only the executable into a minimal base image.

1.  **Build the Docker Image**: From the project's root directory, run the following command. The `-t` flag tags the image with a name for easy identification.

    ```sh
    docker build -t gwi-favourites .
    ```

2.  **Run the Container**: This command starts a new container from your image. The `-p` flag maps the container's internal port `8080` to your host machine, and the `-e` flag injects the required JWT secret as an environment variable.

    ```sh
    docker run -p 8080:8080 -e JWT_SECRET="your_secret_key" gwi-favourites
    ```

3.  **Access the API**: Once the container is running, you can access the API endpoints and the Swagger UI for documentation on `http://localhost:8080/swagger/index.html`.

---

### 🔐 Authentication and Authorization

This API uses **JSON Web Tokens (JWTs)** for authentication. All endpoints for managing favourites are protected and require a valid JWT.

* **JWT Endpoint**: You can obtain a JWT by making a `POST` request to the `/jwt` endpoint with a `userID` in the request body. This endpoint is responsible for issuing tokens with different roles based on the provided user ID.

    ```sh
    curl -X POST http://localhost:8080/jwt \
    -H "Content-Type: application/json" \
    -d '{"userID": "user_id_here"}'
    ```

* **Role-Based Access**: The API currently implements a simple role-based access control based on the `userID` provided during token creation at the `/jwt` endpoint.

    * **Admin User**: A `userID` of `"admin_user_id"` grants **administrative access**, allowing the token holder to access and manage resources for **all users** on the API.
    * **Casual User**: Any other `userID` (e.g., `"user_id_here"`) is treated as a **casual user**. Their tokens are restricted to managing only **their own resources** and cannot access or modify data belonging to other users.

* **`JWT_SECRET` Environment Variable**: The application requires a `JWT_SECRET` to sign and verify tokens. This **must be provided** as an environment variable when running the Docker container or configured via `config.json`. The application will **fail to start** if this secret is missing, as it is crucial for cryptographic security.

---

### 🏛️ Application Architecture

The application is built with a clean, layered architecture to promote testability and maintainability.

* **Handlers**: Located in the `handlers` package, these components use the **Gin framework** to manage HTTP requests and responses. They parse request data and delegate business logic to the service layer.
* **Service Layer**: The `favourite_service` package contains the core business logic. It's designed to be independent of the handlers and relies on an interface (`FavouriteServiceIntf`) for dependency injection.
* **Repository**: The `repository` package defines the `Repository` interface, which abstracts the data storage layer. The `InMemoryRepository` is a concrete, concurrency-safe implementation that uses mutexes to prevent race conditions during concurrent access.
* **Middleware**: The application uses a custom `ErrorHandler` middleware to catch and handle API errors uniformly, ensuring a consistent JSON response format across all endpoints. Authentication is handled by **JWT middleware**, which protects all favourite-related endpoints.
* **Configuration**: The project uses the **Viper library** for structured configuration. It loads settings from a `config.json` file, environment variables (`JWT_SECRET`), or a set of default values.
* **Seeding**: The `seed.go` file provides a practical way to populate the in-memory repository with dummy data, making it easy to set up a development environment with a pre-configured state.

---

### 🧪 Unit Testing

The project's unit tests are designed to be thorough and reliable.

* **Mocking**: Tests for handlers and services use **mock objects** (e.g., `MockFavouriteService` and `MockRepository`) to simulate dependencies. This approach ensures tests are fast, isolated, and reliable, as they don't depend on a real database or external service.
* **Test-Driven Approach**: The test suite covers a wide range of scenarios, including:
    * Successful API calls and the validation of JSON responses.
    * Handling invalid JSON payloads and missing required parameters.
    * Specific service-level errors, such as a repository failing a database call.
    * Ensuring paging logic correctly sorts and returns paginated results.
* **Table-Driven Tests**: This common Go practice is used to efficiently test a function with multiple inputs and expected outputs, which reduces code duplication and improves readability.
* **Assertions**: The `github.com/stretchr/testify/assert` package is used for writing clear and expressive test assertions, such as `assert.NoError` and `assert.Equal`.

---

### 🔗 API Endpoints

The API documentation is automatically generated from the code's annotations and is available via the Swagger UI. You can view all available endpoints, their parameters, and their expected responses by navigating to:

**http://localhost:8080/swagger/index.html**

The main endpoints for managing favourites are:

* `GET /users/{userID}/favourites`: Get a user's paginated list of favourite assets.
* `POST /users/{userID}/favourites`: Add a new asset to a user's favourites.
* `PATCH /users/{userID}/favourites/{assetID}`: Update a specific favourite's description.
* `DELETE /users/{userID}/favourites/{assetID}`: Remove a specific favourite asset.