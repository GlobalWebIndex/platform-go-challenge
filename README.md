## Prerequisites

- Go 1.23 or higher
- Make (optional, for using Makefile commands)
- Docker (optional, for containerized deployment)

## Quick Start

1. Build and run: 
    ```bash
    make docker-run
    ```
    the service recreates the db on startup with two tables 
    - assets 
    - favorites (with data for three users (user1,user2,user3))
    

2. Select user:
   ```bash
   make use-user1
   ```

3. Get favorites:
    ```bash
    make api-list PAGE=2 LIMIT=10
    
    ```
4. Remove favorite
    ```bash
    make api-remove ASSET_ID=<id>
    
    ```
5. Update Desc
    ```bash
    make api-update ASSET_ID=<id> DESC="test"
    
    ```

6. You can access the db by running
    ```bash
    make exec-sql
    ```
    query examples:
    ```sql
    SELECT * FROM favorites;
    SELECT * FROM assets;
    ```



## API Endpoints

### Authentication
All endpoints except /swagger require a valid JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

### Available Endpoints
- **GET /swagger**

- **GET /favorites**
  - List user's favorites with pagination
  - Query parameters: 
    - page (default: 1)
    - limit (default: 20)

- **POST /favorites/{assetId}**
  - Add an asset to favorites

- **DELETE /favorites/{assetId}**
  - Remove an asset from favorites

- **PATCH /favorites/{assetId}/description**
  - Update asset description
  - Body: `{"description": "new description"}`

## Testing

Run the test suite:
```bash
go test ./...
```
Run all curl requests:
```bash
make api-test-all PAGE=2 LIMIT=10 ASSET_ID=<id> DESC="test"
```

## Project Structure

```
├── cmd
│   └── challenge
├── internal
│   ├── api
│   │   ├── docs
│   │   ├── handlers
│   │   ├── middleware
│   │   └── types
│   ├── app
│   ├── domain
│   ├── repository
│   │   └── testdata
│   └── service
│       └── types
└── tools
    └── jwt
```
