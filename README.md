# GlobalWebIndex Engineering Challenge


## Quick Start

1. Run tests
```
docker build --target run-test-stage --progress plain --no-cache -t yours-tests-image-name .
```
2. Build Release image
```
docker build --target build-release-stage --no-cache -t yours-image-name .
```

3. Run Image

```
docker run -p 3333:3333 yours-image-name
```

## Overview

This Go web server allows users to manage a personal list of favorite assets. Assets can be of three types: Chart, Insight, or Audience. Users can add, remove their favorites and edit assets.

- **Web server**: Go’s http package
- **Routing**: Handled in `main.go` with `http.ServeMux`
- **Endpoints**:
    1. `/favAssets`: Fetch the user's favorite endpoints.
    2. `/editAssets`: Edit the description of any asset that exists.
    3. `/deleteAssets`: Remove asset from favorite users list.
    4. `/addAssets`: Add a new asset in user's favorites.
- **Storage**: For simplicity in memory structs (`data/dbData.go`) are used.

Diagram of the structs:

+-------------+         +--------------+            +---------------------------------+
|   Users     |         |  Assets      |            |    Favorites                    |      
+-------------+         +--------------+            +---------------------------------+
| UserId (PK) |         | ID (PK)      |            |   ID (PK)                       |
| Name        |         | Type         |            |   UserId (FK -> Users.UserId)   |
+-------------+         | Description  |            |   AssetId (FK -> Assets.ID)     |
                        +---------------+           +---------------------------------+

                             ^
         -----------------------------------------------------------------------------
         |                                      |                                    |
     +--------------------------+          +--------------------------+          +---------------------------+
     | Chart                    |          | Insight                  |          |Audience                   |
     +--------------------------+          +--------------------------+          +---------------------------+
     | ChartId (PK)             |          | InsightId (PK)           |          | AudienceId (PK)           |
     | AssetID (FK -> Assets.ID)|          | AssetID (FK -> Assets.ID)|          | AssetID (FK -> Assets.ID) |
     | XAxis                    |          +--------------------------+          | Gender                    |
     | YAxis                    |                                                | AgeGroup                  |
     | Data                     |                                                | BirthCountry              |
     +--------------------------+                                                | HoursOnSocial             |
                                                                                 | PurchasesLastM            |
                                                                                 +---------------------------+

## Example Requests

1. Fetch users favorites assets

### Request

```
curl --location 'http://localhost:3333/favAssets' \
--header 'Content-Type: application/json' \
--data '{
  "user_id": 1
}'
```

### Response
```
{
    "user_id": 1,
    "assetData": [
        {
            "fav_id": 1,
            "description": "Annual inflation rate",
            "characteristics": {
                "chartId": 1,
                "assetId": 2,
                "title": "Monthly Sales",
                "x_axis": "Month",
                "y_axis": "Revenue",
                "values": [
                    100,
                    200,
                    400
                ]
            }
        },
        {
            "fav_id": 2,
            "description": "Insights for youngs",
            "characteristics": {
                "insightId": 1,
                "asset_id": 1,
                "title": "Insights for youngs"
            }
        }
    ],
    "msg": "Success",
    "code": 200
}

```

2. Add asset into favorites

### Request

```
curl --location 'http://localhost:3333/addAssets' \
--header 'Content-Type: application/json' \
--data '{
    "user_id": 1,
    "asset_id": 5
}'
```

### Response

```
{
    "success": true,
    "msg": "New asset added succesfully",
    "code": 200
}
```

3. Delete asset from favorite

### Request

```
curl --location 'http://localhost:3333/deleteAssets' \
--header 'Content-Type: application/json' \
--data '{
    "user_id": 1,
    "asset_id": 2
}'

```

### Response

```
{
    "success": true,
    "msg": "Success deletion",
    "code": 200
}
```

4. Change the description of an asset

### Request

```
curl --location 'http://localhost:3333/editAssets' \
--header 'Content-Type: application/json' \
--data '{
    "asset_id": 1,
    "new_description": "Insights for youngs here"
}'

```
### Response

```
{
    "success": true,
    "msg": "Success edit of an asset",
    "code": 200
}

```

## Future Improvements

1. Add a real database for persistence

2. Add authentication mechanism

3. Optimize large asset lists (pagination, caching)