# GWI Favorites Service API - cURL Commands

Ready-to-use curl commands for all endpoints. Replace sample UUIDs with actual values in testing.

**Base URL:** `http://localhost:8080/api/v1`

---

## Health Check

### 1. GET - Health Check
Check if service is running and database is accessible.

```bash
curl -X GET "http://localhost:8080/health"
```

---

## User Management

### 2. GET - List All Users
Retrieve a paginated list of all users.

**Query Parameters:**
- `page` (default: 1): Page number
- `limit` (default: 20, max: 100): Results per page

```bash
curl -X GET "http://localhost:8080/api/v1/users?page=1&limit=20"
```

### 3. POST - Create New User
Create a new user in the system. User ID is auto-generated.

```bash
curl -X POST "http://localhost:8080/api/v1/users" -H "Content-Type: application/json" -d '{}'
```

### 4. DELETE - Delete User
Delete a user and all their associated data.

**Path Parameters:**
- `userID`: UUID of the user to delete

Example (replace cac2b0c4-241d-4109-a39e-ec60cd62c0e9 with actual user ID):

```bash
curl -X DELETE "http://localhost:8080/api/v1/users/cac2b0c4-241d-4109-a39e-ec60cd62c0e9"
```

---

## Asset Management

### 5. GET - List All Assets
Retrieve a paginated list of all assets.

**Query Parameters:**
- `page` (default: 1): Page number
- `limit` (default: 20, max: 100): Results per page
- `type` (optional): Filter by type [chart, insight, audience]

```bash
curl -X GET "http://localhost:8080/api/v1/assets?page=1&limit=20"
```

With type filter:
```bash
curl -X GET "http://localhost:8080/api/v1/assets?page=1&limit=20&type=chart"
```

### 6. POST - Create Chart Asset
Create a new chart asset with visual data representation.

**Request Body:**
- `type`: "chart" (required)
- `data.title`: Chart title
- `data.x_axis`: X-axis label
- `data.y_axis`: Y-axis label
- `data.data`: Array of numeric values

```bash
curl -X POST "http://localhost:8080/api/v1/assets" \
  -H "Content-Type: application/json" \
  -d '{"type": "chart", "data": {"title": "Monthly Sales", "x_axis": "Month", "y_axis": "Revenue ($)", "data": [1200, 1500, 1800, 2100, 2400, 2800]}}'
```

### 7. POST - Create Insight Asset
Create a new insight asset with textual findings.

**Request Body:**
- `type`: "insight" (required)
- `data.text`: The insight text
- `data.topic`: Topic classification

```bash
curl -X POST "http://localhost:8080/api/v1/assets" \
  -H "Content-Type: application/json" \
  -d '{"type": "insight", "data": {"text": "Users aged 25-34 show 45% higher engagement on mobile devices", "topic": "Mobile Engagement"}}'
```

### 8. POST - Create Audience Asset
Create a new audience asset with demographic segment definition.

**Request Body:**
- `type`: "audience" (required)
- `data.gender`: "Male" or "Female"
- `data.birth_country`: Country name
- `data.age_groups`: Array of age ranges (e.g., ["25-34", "35-44"])
- `data.social_media_hours_daily`: Time range (e.g., "3-5")
- `data.purchases_last_month`: Integer count

```bash
curl -X POST "http://localhost:8080/api/v1/assets" \
  -H "Content-Type: application/json" \
  -d '{"type": "audience", "data": {"gender": "Female", "birth_country": "United States", "age_groups": ["25-34", "35-44"], "social_media_hours_daily": "3-5", "purchases_last_month": 5}}'
```

### 9. DELETE - Delete Asset
Delete an asset from the system.

**Path Parameters:**
- `assetID`: UUID of the asset to delete

Example (replace 125a5ed8-cdb3-45c7-b82a-94ae32f95031 with actual asset ID):

```bash
curl -X DELETE "http://localhost:8080/api/v1/assets/125a5ed8-cdb3-45c7-b82a-94ae32f95031"
```

---

## Favorite Management

### 10. GET - Get User's Favorites
Retrieve a paginated list of assets favorited by a user. Sorted by newest first.

**Path Parameters:**
- `userID`: UUID of the user

**Query Parameters:**
- `page` (default: 1): Page number
- `limit` (default: 20, max: 100): Results per page
- `type` (optional): Filter by asset type [chart, insight, audience]

```bash
curl -X GET "http://localhost:8080/api/v1/users/cac2b0c4-241d-4109-a39e-ec60cd62c0e9/favorites?page=1&limit=20"
```

With type filter:
```bash
curl -X GET "http://localhost:8080/api/v1/users/cac2b0c4-241d-4109-a39e-ec60cd62c0e9/favorites?page=1&limit=20&type=chart"
```

### 11. POST - Add Asset to Favorites
Add an asset to a user's favorites list.

**Path Parameters:**
- `userID`: UUID of the user

**Request Body:**
- `asset_id`: UUID of the asset to favorite (required)
- `description`: Optional custom description

Example with custom description:

```bash
curl -X POST "http://localhost:8080/api/v1/users/cac2b0c4-241d-4109-a39e-ec60cd62c0e9/favorites" \
  -H "Content-Type: application/json" \
  -d '{"asset_id": "125a5ed8-cdb3-45c7-b82a-94ae32f95031", "description": "Important quarterly metrics"}'
```

Without custom description:
```bash
curl -X POST "http://localhost:8080/api/v1/users/cac2b0c4-241d-4109-a39e-ec60cd62c0e9/favorites" \
  -H "Content-Type: application/json" \
  -d '{"asset_id": "447c5b65-dd2d-4458-92fd-05a72ffb6c74"}'
```

### 12. PUT - Update Favorite Description
Update the custom description for a favorited asset.

**Path Parameters:**
- `userID`: UUID of the user
- `assetID`: UUID of the asset

**Request Body:**
- `description`: Updated description text (required)

Example (replace cac2b0c4-241d-4109-a39e-ec60cd62c0e9 and 125a5ed8-cdb3-45c7-b82a-94ae32f95031 with actual IDs):

```bash
curl -X PUT "http://localhost:8080/api/v1/users/cac2b0c4-241d-4109-a39e-ec60cd62c0e9/favorites/125a5ed8-cdb3-45c7-b82a-94ae32f95031" \
  -H "Content-Type: application/json" \
  -d '{"description": "Updated: Key performance indicators for Q4"}'
```

### 13. DELETE - Remove Asset from Favorites
Remove an asset from a user's favorites list (soft delete).

**Path Parameters:**
- `userID`: UUID of the user
- `assetID`: UUID of the asset

Example (replace cac2b0c4-241d-4109-a39e-ec60cd62c0e9 and 125a5ed8-cdb3-45c7-b82a-94ae32f95031 with actual IDs):

```bash
curl -X DELETE "http://localhost:8080/api/v1/users/cac2b0c4-241d-4109-a39e-ec60cd62c0e9/favorites/125a5ed8-cdb3-45c7-b82a-94ae32f95031"
```

---