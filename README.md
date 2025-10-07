# GWI challenge

A dashboard for assets. 


<i>This project is the [technical assignment for GWI](https://github.com/GlobalWebIndex/platform-go-challenge).</i>


## Table of Contents

- Technical Implementation
- Software versions
- Installation & Deployment
- Features
- Data population
- Authentication
- Potential enhancements
- Bugs
- License
- Contact

## Technical Implementation
- The backend has been implemented by using Go.

The main goals during development were:
- Separation of concerns: each entity has a route, a service and a model.
- Security:
    - JWT
    - Unified response
- Error handling
- Testing
- Clear naming and documentation

## Software versions
- Python: 1.25.1
- Docker: 28.4.0

## Installation & Deployment

```bash
# Clone the repository
git clone git@github.com:stevenoddity/platform-go-challenge.git

# Navigate to the project directory
cd platform-go-challenge

# Build docker image
docker build -t gwi-challenge .

# Run docker image 
docker run -p 8080:8080 gwi-challenge

```

```bash
Application can be accessed at http://127.0.0.1:8080/
```
Command for running tests:
```bash
go test -v ./...
```

# Features
- ## List Favorites
```bash
curl -i -X GET -H "Authorization: Bearer $JWT" "http://127.0.0.1:8080/favorites
```
Example:
```bash
curl -i -X GET -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.fKCJWNXwhs7ukzI7vpAN2v1z5PBFmiqLlAEhoxbuDB4" "http://127.0.0.1:8080/favorites"
```
- ## Add a new Favorite
```bash
curl -i -X POST -H "Authorization: Bearer $JWT" "http://127.0.0.1:8080/favorites" -d '{
    "asset_id": {asset_id}
  }'
```
Example:
```bash
curl -i -X POST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.fKCJWNXwhs7ukzI7vpAN2v1z5PBFmiqLlAEhoxbuDB4" "http://127.0.0.1:8080/favorites" -d '{
    "asset_id": 3
  }'
```
- ## Delete a Favorite
```bash
curl -i -X DELETE -H "Authorization: Bearer $JWT" "http://127.0.0.1:8080/favorites/{favorite_id}"
```
Example:
```bash
curl -i -X DELETE -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.fKCJWNXwhs7ukzI7vpAN2v1z5PBFmiqLlAEhoxbuDB4" "http://127.0.0.1:8080/favorites/1"
```
- ## Edit Asset's description
```bash
curl -i -X PUT -H "Authorization: Bearer $JWT" "http://127.0.0.1:8080/assets/{asset_id}" -d '{
    "data": {
      "new_data_field": "new_value",
      "new_data_field_2": 64000
    },
    "new_field": "example"
  }'
```
Example:
```bash
curl -i -X PUT -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.fKCJWNXwhs7ukzI7vpAN2v1z5PBFmiqLlAEhoxbuDB4" "http://127.0.0.1:8080/assets/1" -d '{
    "data": {
      "category": "crypto",
      "price": 64000
    },
    "new_field": "example"
  }'
```

# Data population

Initial dataset can be located at:
```bash
database/database.go
```

# Authentication
## JWT token for user_id = 1:
```bash
 eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.fKCJWNXwhs7ukzI7vpAN2v1z5PBFmiqLlAEhoxbuDB4
```
## Generate JWT token ad hoc:
Use secret-key at https://www.jwt.io/
```bash
secret_key="gwi-jwt-secret"
```

# Potential Enhancements
- Enhance testing
- Use database: NoSQL (document-based) would be my suggestion due to the unstructured format of data (i.e. mongoDB or ElasticSearch)
- Logging 
- Pagination
- Lazy loading (minimize json loadind)
- Concurrency (use locks)
- Authorization on authenticated users
- Expiration in JWT
- Blacklist old JWT tokens
- Rate limiter on endpoints
- Enhance input validators
- Use Maps instead of loops

# Bugs

- Concurency issues 
- HTTP Response Status of creation should be 201

# Licences

# Contact

Don't hesitate to ask me for clarifications at steveofsam@gmail.com
