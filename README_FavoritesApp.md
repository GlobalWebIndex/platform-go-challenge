# FavoritesApp

## Description
This is a backend application for managing user assets and favorites. It allows users to add, retrieve, remove, and edit various types of assets (charts, insights, audiences) and mark them as favorites.

## Features
- User management (implicit through asset ownership)
- CRUD operations for different asset types (Chart, Insight, Audience)
- Ability to mark assets as favorites
- PostgreSQL for data persistence

### Prerequisites

- Go (version 1.22 or higher)

- Docker (for running PostgreSQL and Redis)

- Redis

### Steps

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-username/FavoritesApp.git
    cd FavoritesApp
    ```

2.  **Set up PostgreSQL using Docker Compose:**
    ```bash
    docker-compose up -d
    ```
    This will start a PostgreSQL container. The database schema must be initialized on startup (refer to `InitQuery.sql`).

3.  **Install Go dependencies:**
    ```bash
    go mod tidy
    ```

4.  **Build the application:**
    ```bash
    go build -o main .
    ```

## Usage

### Running the application
```bash
./main
```
The application will start on `http://localhost:8080`. 

### API Endpoints

#### Assets

-   **GET /users/{userId}/assets**
    -   Get all assets for a specific user.
    -   Example: `GET /users/a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11/assets`

-   **POST /users/{userId}/assets**
    -   Add a new asset for a user.
    -   Request Body (example for a Chart):
        ```json
        {
            "type": "chart",
            "id": "chart-123",
            "title": "Sales Performance",
            "description": "Monthly sales data",
            "x_axis_title": "Month",
            "y_axis_title": "Revenue",
            "data": [
                {"datapoint_code": "JAN", "value": 100.5},
                {"datapoint_code": "FEB", "value": 120.0}
            ]
        }
        ```
    -   Request Body (example for an Insight):
        ```json
        {
            "type": "insight",
            "id": "insight-456",
            "description": "Key market trend analysis"
        }
        ```
    -   Request Body (example for an Audience):
        ```json
        {
            "type": "audience",
            "id": "audience-789",
            "gender": "Male",
            "country": "USA",
            "age_group": "25-34",
            "social_hours": 3,
            "purchases": 5,
            "description": "Target audience segment"
        }
        ```

-   **DELETE /users/{userId}/assets/{assetId}**
    -   Remove an asset for a user.
    -   Example: `DELETE /users/a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11/assets/chart-123`

-   **PUT /users/{userId}/assets/{assetId}**
    -   Edit an asset's description for a user.
    -   Request Body:
        ```json
        {
            "description": "Updated description for the asset"
        }
        ```
    -   Example: `PUT /users/a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11/assets/chart-123`

#### Favorites

-   **GET /users/{userId}/favourites**
    -   Get all favorite assets for a specific user.
    -   Example: `GET /users/a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11/favourites`

-   **POST /users/{userId}/favourites/{assetId}**
    -   Add an asset to a user's favorites.
    -   Example: `POST /users/a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11/favourites/chart-123`

-   **DELETE /users/{userId}/favourites/{assetId}**
    -   Remove an asset from a user's favorites.
    -   Example: `DELETE /users/a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11/favourites/chart-123`

## Technologies Used
- Go
- PostgreSQL
- Docker
- Redis
- Gorilla Mux (for routing)
- Testcontainers-go (for testing)


## License
[MIT License](LICENSE) - (Assuming MIT license, will add a LICENSE file if not present)
