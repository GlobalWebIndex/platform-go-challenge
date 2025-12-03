# Go Challenge - Favourites API

A REST API built with Go.

## Prerequisites

- Go 1.25.4
- Git

## Installation & Setup

### 1. Clone the repository

```bash
git clone <repository-url>
cd go-challenge
```

### 2. Download dependencies

```bash
go mod download
```

### 3. Run the application

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

Once the server is running, you can use the following endpoints:

### Get User's Favourites

```bash
GET http://localhost:8080/api/users/{userID}/favourites
```

### Add Asset to Favourites

```bash
POST http://localhost:8080/api/users/{userID}/favourites
Content-Type: application/json

{
  "assetID": "your-asset-id"
}
```

### Remove Favourite

```bash
DELETE http://localhost:8080/api/users/{userID}/favourites/{assetID}
```

### Update Asset Description

```bash
PATCH http://localhost:8080/api/assets/{assetID}/description
Content-Type: application/json

{
  "description": "New description"
}
```

## Building for Production

TODO...

## Project Structure

```
go-challenge/
├── handlers/          # HTTP request handlers
├── repository/        # Data access layer (in-memory storage)
├── main.go           # Application entry point
├── go.mod            # Go module dependencies
└── go.sum            # Dependency checksums
```

## Dependencies

- `github.com/google/uuid` - UUID generation

# Testing the API

# 🚀 API Test Scenarios

This document outlines the API calls for testing the core functionality of user favourites and asset updates.

---

## 1. Get Favourites (Seeded Data)

Use these commands to retrieve a user's list of favourited assets, demonstrating different pagination options.

---

### 1.1. Get First Page (Default 10 items)

This command uses the default page and limit settings.

```bash
curl "http://localhost:8080/api/users/user-001/favourites"
```

### 1.2. Get Page 2 (Default 10 items)

This explicitly requests the second page of results.

```bash
curl "http://localhost:8080/api/users/user-001/favourites?page=2"
```

### 1.3. Get Page 2 with Custom Limit (5 items)

This requests the second page and sets the page size to 5 items.

```bash
curl "http://localhost:8080/api/users/user-001/favourites?page=2&limit=5"
```

### 1.4. Get Max Page Size (First 20 items)

This requests the first page with the maximum allowed limit of 20 items.

```bash
curl "http://localhost:8080/api/users/user-001/favourites?page=1&limit=20"
```

## 2. Add Asset to Favourites

Add a specific asset (`asset-040`) to the user's favourites list.

### Bash Command

```bash
curl -X POST http://localhost:8080/api/users/user-001/favourites \
  -H "Content-Type: application/json" \
  -d '{"asset_id":"asset-040"}'
```

### PowerShell Command

```PowerShell
$uri = "http://localhost:8080/api/users/user-001/favourites"
$body = '{"asset_id": "asset-040"}'
$headers = @{"Content-Type" = "application/json"}
Invoke-WebRequest -Uri $uri -Method POST -Headers $headers -Body $body
```

3. 🔍 Verify New Favourite Added
   Retrieve the favourites list again to confirm that the new asset (asset-101) is present.

### Bash Command

```bash
curl  "http://localhost:8080/api/users/user-001/favourites?page=3"
```

4. 📝 Update Asset's Description
   Update the description of the newly favourited asset (asset-040).
   This is an asset-level operation.

### Bash Command

```bash
curl -X PATCH http://localhost:8080/api/assets/asset-040/description \
  -H "Content-Type: application/json" \
  -d '{"description": "Updated: Comprehensive study on millennial social media habits based on 2024 research"}'
```

### Powershell

```PowerShell
$uri = "http://localhost:8080/api/assets/asset-040/description"
$body = '{"description": "Updated: Comprehensive study on millennial social media habits based on 2024 research"}'
$headers = @{"Content-Type" = "application/json"}
Invoke-RestMethod -Uri $uri -Method PATCH -Headers $headers -Body $body
```

5. 🔄 Verify Description Change
   Retrieve the favourites list again. The asset entry for asset-101 should now reflect the updated description.

### Bash Command

```bash
curl  "http://localhost:8080/api/users/user-001/favourites?page=3"
```

6. Remove a Favourite
   Remove the asset (asset-101s) from the user's favourites list.

### Bash Command

```bash
curl -X DELETE http://localhost:8080/api/users/user-001/favourites/asset-040
```

7.  Verify Removal
    Retrieve the favourites list one last time to confirm that asset-003 is no longer included.

### Bash Command

```bash
curl  http://localhost:8080/api/users/user-001/favourites
```
