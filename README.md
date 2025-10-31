# Preferred Assets API

A Go-based REST API for managing assets, users, and favourites with JWT authentication and role-based access control.

## Overview

This API allows users to manage three types of assets (audiences, charts, and insights) and maintain personal favourite lists. The system supports user management, asset operations, and favourites tracking with secure authentication.

## Features

- **User Management**: Create, read, update, and delete users
- **Asset Management**: Create and delete assets of three types (audience, chart, insight)
- **Favourites System**: Add and remove assets from user favourites
- **JWT Authentication**: Secure endpoints with Keycloak integration
- **Role-Based Access Control**: Admin and user roles with different permissions
- **RESTful API**: Clean, well-documented endpoints following OpenAPI specification

## Technology Stack

- **Backend**: Go framework
- **Authentication**: Keycloak with OAuth 2.0
- **Documentation**: Swagger/OpenAPI 2.0
- **Containerization**: Docker with Docker Compose
- **Data Storage**: In-memory storage (easily replaceable with persistent storage)

## API Endpoints

### Authentication
All endpoints require JWT Bearer token authentication. Obtain a token using:
```bash
curl --location 'http://localhost:8090/realms/preferred-assets-realm/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'grant_type=password' \
--data-urlencode 'client_id=preferred-assets-api' \
--data-urlencode 'client_secret=AIz24JAXB4uQNsQ8uoANOFVobtFVvzlu' \
--data-urlencode 'username=admin' \
--data-urlencode 'password=admin'
```

### Users
- `GET /api/v1/users` - Get all users (Admin only)
- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user
- `GET /api/v1/users/{id}/favourites` - Get user favourites

### Assets
- `POST /api/v1/assets` - Create a new asset
- `DELETE /api/v1/assets/{assetId}` - Delete an asset

### Favourites
- `POST /api/v1/favourites` - Add asset to favourites
- `DELETE /api/v1/favourites/{userId}/{assetId}` - Remove asset from favourites

## Asset Types

### 1. Audience
Represents demographic segments with attributes:
- Gender (Male/Female)
- Birth country
- Age group
- Hours spent on social media daily
- Number of purchases last month

### 2. Chart
Visual data representations with:
- Axes titles
- Data points (2D arrays)
- Chart configuration

### 3. Insight
Text-based insights with:
- Descriptive text
- Key findings and observations

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.19+ (if running locally)

### Running with Docker (Recommended)
```bash
# Clone the repository
git clone <repository-url>
cd preferred-assets-api

# Start all services
docker compose up --build

# The API will be available at http://localhost:8081
# Keycloak admin console at http://localhost:8090
```

### Running Locally (Alternative)
If you encounter connectivity issues between the API and Keycloak in Docker:
Stop the preffered_assets_api from docker and keep keycloak runnning
```bash
# Run the API locally
cd .{saved folder}/preferred_assets_api/cmd/api/main
go run main.go

# The API will be available at http://localhost:8081
```

## Authentication

### Default Users
Two pre-configured users are available:

**Admin User:**
- Username: `admin`
- Password: `admin`
- Role: Full access to all endpoints

**Regular User:**
- Username: `user` 
- Password: `user`
- Role: Standard user permissions

### Obtaining Access Tokens
```bash
# For admin user
curl --location 'http://localhost:8090/realms/preferred-assets-realm/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'grant_type=password' \
--data-urlencode 'client_id=preferred-assets-api' \
--data-urlencode 'client_secret=AIz24JAXB4uQNsQ8uoANOFVobtFVvzlu' \
--data-urlencode 'username=admin' \
--data-urlencode 'password=admin'

# For regular user
curl --location 'http://localhost:8090/realms/preferred-assets-realm/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'grant_type=password' \
--data-urlencode 'client_id=preferred-assets-api' \
--data-urlencode 'client_secret=AIz24JAXB4uQNsQ8uoANOFVobtFVvzlu' \
--data-urlencode 'username=user' \
--data-urlencode 'password=user'
```

## Usage Examples

### Create a User
```bash
curl -X POST "http://localhost:8081/api/v1/users" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_JWT_TOKEN" \
-d '{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePass123"
}'
```

### Create an Asset (Audience)
```bash
curl -X POST "http://localhost:8081/api/v1/assets" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_JWT_TOKEN" \
-d '{
  "id": "audience_001",
  "title": "Young Social Media Users",
  "type": "audience",
  "description": "Audience segment of young adults active on social media",
  "gender": "female",
  "birth_country": "US",
  "age_group": "18-24",
  "hours_social": 15,
  "purchases_last_month": 3
}'
```

### Add to Favourites
```bash
curl -X POST "http://localhost:8081/api/v1/favourites" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_JWT_TOKEN" \
-d '{
  "_id": "user_123",
  "asset_id": "audience_001"
}'
```

### Get User Favourites
```bash
curl -X GET "http://localhost:8081/api/v1/users/user_123/favourites" \
-H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## API Documentation

Full API documentation is available via Swagger UI when the application is running:
- Swagger JSON: `http://localhost:8081/swagger/doc.json`
- Interactive documentation can be viewed using Swagger UI tools

## Development

### Running Tests
```bash
cd .{saved folder}/preferred_assets_api/cmd/api/main
go test ./...
```

### Building Locally
```bash
cd .{saved folder}/preferred_assets_api/cmd/api/main
go build -o preferred-assets-api .
```

## Configuration

The application can be configured through environment variables:

- `API_PORT`: Server port (default: 8081)
- `KEYCLOAK_URL`: Keycloak server URL
- `KEYCLOAK_REALM`: Keycloak realm name
- `KEYCLOAK_CLIENT_ID`: OAuth client ID

## Storage Notes

The current implementation uses in-memory storage for demonstration purposes. For production use, consider implementing persistent storage solutions such as:

- PostgreSQL for relational data
- MongoDB for document storage
- Redis for caching

## License

MIT License - see LICENSE file for details.

## Support

For API support, contact: support@yourapp.com

API Documentation: http://localhost:8081/swagger/doc.json