# GlobalWebIndex Engineering Challenge                                        

##### Table of Contents  
- [Introduction](#introduction)
- [Getting Started](#getting-started)  
  - [Prerequisites](#prerequisites)  
- [Quick Start](#quick-start)
- [Testing](#testing)



## Introduction
A tiny, clean HTTP API in Go that lets users “star” assets (charts, insights, audiences) and fetch their favourites for quick access when enter in thee platform's dashboard. This server has two endpoints implementing the following:

- an endpoint that enables a user to ADD/EDIT(shared attributes)/DELETE favorites
- an endpoint that gets the user id and returns back a list of favorites

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
For the purpose of testing execute the following before initiating the testing:
- `docker compose up`,
-  source the environmental variables containing the BASE_URL and
- `DATABASE_URL='postgres://platform-go-challenge:pass@localhost:5432/platform-go-challenge' \ go run ./src/scripts/seed`

Finally after those, test the endpoints using curl:

1) GET user favorites:
- GET user favorites first element to simulate that you get the first page and get *NEXT_T* and *NEXT_ID*:
```
curl -s "$BASE_URL/v1/users/u1/favorites?limit=1" | jq
```
- GET user favorites next page:
```
curl -s "$BASE_URL/v1/users/u1/favorites?limit=2&after_t=NEXT_T&after_id=NEXT_ID" | jq
```

2) GET user favorites by type:
```
curl -s "$BASE_URL/v1/users/u1/favorites?limit=20&type=chart" | jq
```

3) POST user favorite from assets:
```
curl -i -X POST "$BASE_URL/v1/users/u1/favorites" \
  -H "Content-Type: application/json" \
  -d "{\"assetId\":\"$ASSET_ID\"}"
```

4) DELETE user favorite:
```
curl -i -X DELETE "$BASE_URL/v1/users/u1/favorites/$ASSET_ID"
```

5) PATCH user favorite description:
```
curl -i -X PATCH "$BASE_URL/v1/users/u1/favorites/$ASSET_ID" \
  -H 'Content-Type: application/json' \
  --data '{"description":"My private description"}'
```
