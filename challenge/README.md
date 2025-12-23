# GWI Favorites Service

A backend API for managing user favorite assets. Users can save and organize Charts, Insights, and Audiences for quick access.

## What It Does

The service lets users build personal collections of assets they care about. Multiple users can favorite the same asset independently. When you favorite something, you can add a custom note to remember why it matters to you.

Think of it like bookmarks, but for business data.

## Key Features

- **Three asset types**: Charts (visual data), Insights (findings), Audiences (demographic segments)
- **Personal favorites**: Each user maintains their own list
- **Custom descriptions**: Add notes to explain why you saved something
- **Organized listing**: Filter by type, paginate through results
- **Safe deletion**: Removed favorites are archived, not erased
- **Built on proven tech**: Go + PostgreSQL

## Tech Stack

- **Language**: Go 1.21
- **Database**: PostgreSQL
- **Framework**: Gorilla Mux (HTTP routing)
- **Deployment**: Docker + Docker Compose

## Files Overview

| File                       | Purpose |
|----------------------------|---------|
| `go_impl.go`               | Main application code (handlers, service logic, database queries) |
| `implementation_test.go`   | Unit tests for all endpoints and error cases |
| `schema.sql`               | Database schema and indexes |
| `Dockerfile`               | Container build (includes tests in build process) |
| `docker-compose.yml`       | Orchestrates app + PostgreSQL |
| `swagger-api.yaml`         | API specification |
| `gwi_api_curl_commands.md` | Example API calls |
| `go_mod`                   | Declares project name and dependencies |

## Running the Service

### With Docker (Recommended)

```bash
docker-compose up --build
```

This will:
1. Run all tests during build (fails if any test fails)
2. Start PostgreSQL
3. Start the API on http://localhost:8080

### Locally

You need Go 1.21+ and PostgreSQL:

```bash
# Download dependencies
go mod download

# Create database
createdb gwe_challenge

# Initialize schema
psql gwe_challenge < schema.sql

# Run the server
go run go_impl.go
```

## Testing the API

### Health check
```bash
curl http://localhost:8080/health
```

### Create a user
```bash
curl -X POST http://localhost:8080/api/v1/users
```

### Create an asset
```bash
curl -X POST http://localhost:8080/api/v1/assets \
  -H "Content-Type: application/json" \
  -d '{
    "type": "chart",
    "data": {
      "title": "Daily Users",
      "x_axis": "Date",
      "y_axis": "Count",
      "data": [100, 150, 120, 180]
    }
  }'
```

### Add to favorites
```bash
curl -X POST http://localhost:8080/api/v1/users/{userID}/favorites \
  -H "Content-Type: application/json" \
  -d '{
    "asset_id": "{assetID}",
    "description": "Important metric for Q4"
  }'
```

See `gwi_api_curl_commands.md` for complete examples.

## API Endpoints

### Users
- `GET /api/v1/users` - List all users
- `POST /api/v1/users` - Create user
- `DELETE /api/v1/users/{userID}` - Delete user

### Assets
- `GET /api/v1/assets` - List assets (filter by type)
- `POST /api/v1/assets` - Create asset
- `DELETE /api/v1/assets/{assetID}` - Delete asset

### Favorites
- `GET /api/v1/users/{userID}/favorites` - Get user's favorites (supports pagination and type filtering)
- `POST /api/v1/users/{userID}/favorites` - Add to favorites
- `PUT /api/v1/users/{userID}/favorites/{assetID}` - Update description
- `DELETE /api/v1/users/{userID}/favorites/{assetID}` - Remove from favorites

### System
- `GET /health` - Health check

Full API spec in `swagger-api.yaml`.

## Database Design

Three tables handle the data:

**users** - Stores user profiles
- ID (UUID)
- Created timestamp

**assets** - Stores charts, insights, audiences
- ID (UUID)
- Type (chart/insight/audience)
- Data (JSON - flexible structure per type)
- Created timestamp

**favorites** - Links users to assets
- ID (UUID)
- User ID, Asset ID (foreign keys)
- Optional custom description
- Added timestamp
- Deleted timestamp (soft delete)

The `favorites` table has an index on `(user_id, added_at DESC)` to speed up the most common query: getting a user's recent favorites.

## How It Works

1. **Create a user**: Returns a UUID for that user
2. **Create assets**: Add charts, insights, or audiences to the system (any user can create)
3. **Favorite an asset**: User marks something as a favorite with optional notes
4. **Retrieve favorites**: User gets their personalized list (paginated, filterable)
5. **Update or remove**: Modify descriptions or delete from favorites

Multiple users can favorite the same asset. Deletion is soft (archived) not hard (erased).

## Testing

Tests run automatically during Docker build. If any test fails, the build stops.

Run tests locally:
```bash
go test -v
```

Tests cover:
- All CRUD operations
- Error cases (404, validation)
- Pagination and filtering
- Edge cases (empty lists, duplicates)
- Health check

## Deployment Notes

The service is stateless. You can run multiple instances behind a load balancer, all pointing to the same PostgreSQL database.

For production:
- Use environment variables for database connection
- Set appropriate pool size based on load
- Monitor query performance
- Keep indexes maintained

## Error Handling

The API returns standard HTTP status codes:
- `200` - Success
- `201` - Created
- `204` - Deleted
- `400` - Bad request (validation error)
- `404` - Not found (user or asset missing)
- `409` - Conflict (already favorited)
- `500` - Server error

Error responses include a message explaining what went wrong.

## Performance

**Get favorites**: ~5ms (indexed)
**Add favorite**: ~10ms
**Update description**: ~10ms

The service handles pagination efficiently. Even with thousands of favorites per user, results load instantly because only the requested page is retrieved.

## Code Quality

- 15+ unit tests, all passing
- Proper error handling throughout
- Separation of concerns (handlers, service, storage layers)
- Connection pooling for database efficiency
- Input validation on all endpoints

## What's Next

Future enhancements could include:
- Authentication and authorization
- Rate limiting
- Search and full-text indexing
- Batch operations (add/remove multiple at once)
- Export favorites (CSV, JSON)
- Analytics (most favorited assets)

---

For detailed API examples, see `gwi_api_curl_commands.md`.
For API specification, see `swagger-api.yaml`.