# API Documentation

This API allows clients to manage a list of favorite entries using a RESTful interface.

## Authentication

All endpoints require the `X-User-ID` header to identify the user making the request.
The user ID should be a positive integer representing the user's unique identifier in the system.


**Reference the sample data overview** in the [infra/README.md](../infra/README.md#sample-data-overview) for pre-populated user and asset data.

## Common Response Formats

### Success Response
- **Status Code**: Varies by endpoint
- **Content-Type**: `application/json`

### Error Response
```json
{
  "error": "Error message describing the issue"
}
```

## Endpoints

### 1. Get All User Favorites

Retrieve all favorite assets for a specific user.

- **URL**: `/favorite/`
- **Method**: `GET`
- **Headers**:
    - `X-User-ID` (required): User identifier as integer

#### Success Response
- **Status Code**: `200 OK`
- **Content-Type**: `application/json`
- **Body**:
```json
[
  {
    "ID": 1,
    "UserID": 1,
    "AssetID": 1,
    "Description": "      "
  },
  {
    "ID": 2,
    "UserID": 1,
    "AssetID": 4,
    "Description": "Key Q4 insight"
  }
]
```

#### Error Responses
- **400 Bad Request**: Missing or invalid user ID
- **404 Not Found**: No favorites found for the user
- **500 Internal Server Error**: Server error

#### Example Request
```bash
curl -X GET "http://localhost:8080/favorite/" \
  -H "X-User-ID: 3"
```

#### Example Response
```json
[
  {
    "ID": 5,
    "UserID": 3,
    "AssetID": 6,
    "Description": "Friday peak analysis"
  }
]
```

---

### 2. Add Asset to User Favorites

Adds a new asset to the user's favorite list.

- **URL**: `/favorite/{aid}`
- **Method**: `POST`
- **Headers**:
    - `X-User-ID` (required): User identifier as integer
    - `Content-Type`: `application/json`
- **URL Parameters**:
    - `aid` (required): Asset ID to add to favorites

#### Request Body
```json
{
  "description": "Optional description of the favorite asset"
}
```

#### Success Response
- **Status Code**: `201 Created`
- **Content-Type**: `application/json`
- **Body**:
```json
{
  "id": 1,
  "user_id": 1,
  "asset_id": 123,
  "description": "My favorite asset"
}
```

#### Error Responses
- **400 Bad Request**: Missing/invalid user ID, missing/invalid asset ID, or invalid JSON body
- **409 Conflict**: Asset is already in user's favorites
- **500 Internal Server Error**: Server error

#### Example Request
```bash
curl -X POST "http://localhost:8080/favorite/5" \
  -H "X-User-ID: 1" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "This is my favorite trading asset"
  }'
```

#### Example Response
```json
{
  "ID": 7,
  "UserID": 1,
  "AssetID": 5,
  "Description": "This is my favorite trading asset"
}
```

---

### 3. Update User Favorite Asset

Updates the description of an existing favorite asset.

- **URL**: `/favorite/{id}`
- **Method**: `PATCH`
- **Headers**:
    - `X-User-ID` (required): User identifier as integer
    - `Content-Type`: `application/json`
- **URL Parameters**:
    - `id` (required): Favorite ID to update

#### Request Body
```json
{
  "description": "Updated description"
}
```

#### Success Response
- **Status Code**: `200 OK`
- **Content-Type**: `application/json`
- **Body**:
```json
{
  "id": 1,
  "user_id": 1,
  "asset_id": 5,
  "description": "Updated description"
}
```

#### Error Responses
- **400 Bad Request**: Missing/invalid user ID, missing/invalid favorite ID, or invalid JSON body
- **404 Not Found**: Favorite not found
- **500 Internal Server Error**: Server error

#### Example Request
```bash
curl -X PATCH "http://localhost:8080/favorite/1" \
  -H "X-User-ID: 1" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Updated description for my favorite asset"
  }'
```

#### Example Response
```json
{
  "ID": 1,
  "UserID": 1,
  "AssetID": 1,
  "Description": "Updated description for my favorite asset"
}
```

---

### 4. Remove Asset from User Favorites

Removes an asset from the user's favorites list.

- **URL**: `/favorite/{id}`
- **Method**: `DELETE`
- **Headers**:
    - `X-User-ID` (required): User identifier as integer
- **URL Parameters**:
    - `id` (required): Favorite ID to remove

#### Success Response
- **Status Code**: `204 No Content`
- **Body**: Empty

#### Error Responses
- **400 Bad Request**: Missing/invalid user ID or missing/invalid favorite ID
- **404 Not Found**: Favorite not found
- **500 Internal Server Error**: Server error

#### Example Request
```bash
curl -X DELETE "http://localhost:8080/favorite/15" \
  -H "X-User-ID: 1"
```

#### Example Response
No content returned on successful deletion.

---

## Error Codes and Messages

### Common Error Scenarios

| Status Code | Error Message | Description |
|-------------|---------------|-------------|
| 400 | "X-User-ID header is required" | Missing user ID header |
| 400 | "invalid user ID format" | User ID is not a valid integer |
| 400 | "asset ID is required in URL path" | Missing asset/favorite ID in URL |
| 400 | "invalid asset ID format" | Asset/favorite ID is not a valid integer |
| 400 | "invalid json body" | Request body contains invalid JSON |
| 404 | No favorites found | User has no favorite assets |
| 404 | Favorite not found | Specified favorite doesn't exist |
| 409 | Duplicate favorite | Asset is already in user's favorites |
| 500 | Internal server error | Unexpected server error |

## Notes

- All IDs (user_id, asset_id, favorite_id) must be positive integers
- The `description` field is **not** optional but can be empty
- User authentication is handled via the `X-User-ID` header
- All successful responses (except DELETE) return JSON content
- DELETE operations return no content on success (204 status code)