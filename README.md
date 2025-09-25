# GlobalWebIndex Engineering Challenge                                        

##### Table of Contents  
- [Introduction](#introduction)
- [Getting Started](#getting-started)  
  - [Prerequisites](#prerequisites)  
- [Quick Start](#quick-start)


## Introduction
A tiny, clean HTTP API in Go that lets users “star” assets (charts, insights, audiences) and fetch their favourites for quick access when enter in thee platform's dashboard. This server has two endpoints implementing the following:

- GET /v1/users/:user_id/favorites?limit=&after_t=&after_id=&type=
- POST /v1/users/:user_id/favorites body: {"assetId":"..."}
- DELETE /v1/users/:user_id/favorites/:asset_id
- PATCH /v1/assets/:asset_id body: {"description":"..."}

## Getting Started

### Prerequisites
You’ll need the following to develop and run the API locally:

- Go: **1.23+**
- Docker Engine and Docker Compose v2 
- Postgres 16
- Gin HTTP router

## Quick start (local dev)

You’ll need Docker (for Postgres) and Go 1.22+ (to run the server):
#### Without docker-compose

1) Create your env file
```
cp .env.template .env
# uncomment DATABASE_URL run the server locally (outside docker-compose), be sure .env has:
# DATABASE_URL=postgres://platform-go-challenge:pass@localhost:5432/platform-go-challenge?sslmode=disable
```

2) Start Postgres
```
docker run --name platform-go-challenge-db -p 5432:5432 \
  -e POSTGRES_USER=platform-go-challenge -e POSTGRES_PASSWORD=pass -e POSTGRES_DB=platform-go-challenge \
  -d postgres:16
```

3) Export envs for the Go process
```
# (uncomment DATABASE_URL in .env first if running server locally (outside docker-compose))
source .env
```

4) Seed (required when dev locally only) and run the server
```
go run ./src/scripts/seed
go run ./src/scripts/server
```

#### With docker-compose

0) From the repo root
```
cp .env.template .env   
```

1) Spin-up docker-compose
```
docker compose up --build
```

2) (Optional) Seed sample assets into the compose Postgres
```
DATABASE_URL="postgres://platform-go-challenge:pass@localhost:5432/platform-go-challenge" \
go run ./src/scripts/seed
```

## Testing

You may test the endpoints using curl:

1) GET user favorites:
```
curl -s "$BASE_URL/v1/users/u1/favorites?limit=1" | jq
```

2) GET user favorites by type:
```
curl -s "$BASE_URL/v1/users/u1/favorites?limit=20&type=chart" | jq
```
3) 