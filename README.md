# GWI Favorites Asset Service

### Service Overview

The GWI Favorites Asset Service is a web service that allows users to manage a personalized list of favorite assets
within the GWI platform. Each user can add, view, edit, or remove various types of assets—including Charts, Insights,
and Audiences—from their favorite list.

The GWI Favorites Asset Service is a web service that implements a **user favorites management system** for the GWI
platform. It allows users to `add, view, update, and remove` their personalized list of favorites.

This service focuses **exclusively** on the functionality related to managing user favorites. To simplify the
implementation
and highlight this core logic, the `supporting data—users and assets—are pre-populated in memory and assumed to exist`,
it would be a good practice to read the [Sample Data Overview](infra/README.md#sample-data-overview) before consuming
the service.
This allows the service to operate independently without requiring full authentication, asset creation, or user
management flows.

_Below is a summary of the key API endpoints available in the Favorites Asset Service, which allow users to manage their
favorites effectively._

#### API Endpoints at a Glance

For detailed documentation of the API endpoints, refer to
the [API Documentation](./presentation/README.md#api-documentation) file.

`NOTE: All endpoints require the X-User-ID header to identify the user making the request.`

- **GET** `{base_url}/favorite/` ----> Retrieves all favorite assets for the specified user.

- **POST** `{base_url}/favorite/{aid}` ----> Adds the asset with the given asset ID to the user's favorites.

- **PATCH** `{base_url}/favorite/{id}` ----> Updates the description of a favorite asset.

- **DELETE** `{base_url}/favorite/{id}` ----> Removes an asset from the user's favorites.

## 🏗️ Architecture

This service follows **Clean Architecture** principles with the following layers:

```
┌─────────────────┐
│   Presentation  │ ← HTTP handlers, routing, serialization
├─────────────────┤
│   Application   │ ← Business logic, use cases, orchestration
├─────────────────┤
│     Domain      │ ← Core entities, interfaces, business rules
├─────────────────┤
│   Persistence   │ ← Database access, repository implementations
└─────────────────┘
```

### 📁 Project Structure

```
├── app/                    # Application layer (use cases, business logic)
├── cmd/                    # Application entry points
├── domain/                 # Core domain entities and interfaces
├── infra/                  # Infrastructure setup (Docker, migrations)
├── internal/               # Internal packages (config, server setup)
├── persistence/            # Data access layer (repositories)
└── presentation/           # HTTP handlers and routing
```

## Infrastructure

This service uses a PostgreSQL database for the persistence layer via a dockerized setup.
For details on the database schema, migrations, and Docker setup, refer to
the [infra/README.md documentation](infra/README.md) file.

## 🚀 HowTo Run

### Prerequisites

* `REQUIRED` Docker and Docker Compose installed | [Docker installation guide](https://docs.docker.com/get-docker/)
* `[NECCESSARY FOR LOCAL DEVELOPMENT] ` Go 1.24+ | [Go installation guide](https://golang.org/doc/install)

#### Release Version

##### Running with Docker Compose

1. Start the application stack using Docker Compose:
   ```bash
   docker-compose up
   ```

##### Local Development Setup

For local development:

1. Make sure you have Go 1.24+ installed on your machine. If not, follow the installation instructions for your OS:
   ```bash
   # Follow installation instructions for your OS from https://go.dev/doc/install
   ```

2. Install required tools:
   ```bash
   # Install Mockgen for mock generation
   go install go.uber.org/mock/mockgen@latest
   ```

3. Start the database service:
   ```bash
   cd infra/docker && docker-compose up -d
   ```
   > NOTE: You can skip this step if you have your own database setup, just adjust the configuration files accordingly.
   > For more details, see [infra/README.md](infra/README.md) documentation.

4. Environment Variables

   The application uses environment variables for configuration. **NO REQUIRED environment variables are needed** since
   they
   are set to **default** values, but you can override them.

   ```bash
   ADDR=:8080
   DB_NETWORK=tcp
   DB_ADDRESS=localhost:5432
   DB_USER=gwi
   DB_PASSWORD=gwi
   DB_NAME=gwi
   DB_APP_NAME=gwi
   ```

   > NOTE: The default values are set to match the Docker Compose setup.

5. Run the application
   ```bash
   go run ./cmd/main.go
   ```

## CI/CD Pipeline

The project includes GitHub Actions for continuous integration.

### GitHub Actions

The GitHub Actions workflow at `.github/workflows/ci.yml` runs:

- Unit tests.
- Database integration tests using **Testcontainers**.

## Testing Strategy

- **Unit Tests**: For individual components, and with the advantages of clean architecture, each component can be tested
  in isolation, you can see that each package(directory) has its own test file.
- **Database Integration Tests**: Integration tests for the persistence layer, these tests verify that our
  components work correctly together with real dependencies (database, cache, etc.) rather than
  mocks using **Testcontainers**.

### Running Tests

```bash
go test -cover ./... -v -timeout 5m
```

### Generating Mocks

The project uses Go's mockgen for generating mock implementations of interfaces for testing.
The domain interfaces have //go:generate directives that specify how to generate mock implementations.

To generate or update all mocks:

```bash
# From the project root directory
go generate ./...
```

This will create or update mock files in the appropriate directories:
