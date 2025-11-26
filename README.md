# GlobalWebIndex Engineering Challenge

## Introduction

This challenge is designed to give you the opportunity to demonstrate your abilities as a software engineer and specifically your knowledge of the Go language.

On the surface the challenge is trivial to solve, however you should choose to add features or capabilities which you feel demonstrate your skills and knowledge the best. For example, you could choose to optimise for performance and concurrency, you could choose to add a robust security layer or ensure your application is highly available. Or all of these.

Of course, usually we would choose to solve any given requirement with the simplest possible solution, however that is not the spirit of this challenge.

## Challenge

Let's say that in GWI platform all of our users have access to a huge list of assets. We want our users to have a peronal list of favourites, meaning assets that favourite or “star” so that they have them in their frontpage dashboard for quick access. An asset can be one the following
* Chart (that has a small title, axes titles and data)
* Insight (a small piece of text that provides some insight into a topic, e.g. "40% of millenials spend more than 3hours on social media daily")
* Audience (which is a series of characteristics, for that exercise lets focus on gender (Male, Female), birth country, age groups, hours spent daily on social media, number of purchases last month)
e.g. Males from 24-35 that spent more than 3 hours on social media daily.

Build a web server which has some endpoint to receive a user id and return a list of all the user’s favourites. Also we want endpoints that would add an asset to favourites, remove it, or edit its description. Assets obviously can share some common attributes (like their description) but they also have completely different structure and data. It’s up to you to decide the structure and we are not looking for something overly complex here (especially for the cases of audiences). There is no need to have/deploy/create an actual database although we would like to discuss about storage options and data representations.

Note that users have no limit on how many assets they want on their favourites so your service will need to provide a reasonable response time.

A working server application with functional API is required, along with a clear readme.md. Useful and passing tests would be also be viewed favourably

It is appreciated, though not required, if a Dockerfile is included.

## Submission

Just create a fork from the current repo and send it to us!

Good luck, potential colleague!

---

## My project

In this repo, I created (with help of tutorials, golang documentation and Claude Code)
* a simple rest api that enables the user to create, update, delete and get for each of the 3 asset types, as well as set any assset as his/her favourite/stared
* a simple graphql api with the same functionality as the rest api (both query and mutation functionality is included) as well as a query to get all assets that a user has put a star on.

I used Go as requested and I have am storing the data in a very simple Postgresql db schema that is created by models in the [models](models/) directory.

I used for the 
* Rest api https://github.com/gin-gonic/gin
* Graphql https://github.com/99designs/gqlgen
* ORM https://github.com/go-gorm/gorm

An overview of the structure of the project can be found at [PROJECT_STRUCTURE.md](docs/PROJECT_STRUCTURE.md). All the information on the api and how to use it can be found at See [API_DOCUMENTATION.md](docs/API_DOCUMENTATION.md). Information on the types of tests and how to use them can be found at [TESTING_GUIDE](docs/TESTING_GUIDE.md).

## Running the Application

### Prerequisites 
1. PostgreSQL database running
2. `.env` file with database connection strings:
   ```bash
   DB_URL=postgres://user:password@localhost:5432/mydb?sslmode=disable
   TEST_DB_URL=postgres://user:password@localhost:5432/mydb?sslmode=disable
   ```
3. Docker with docker-compose 

### Development

```bash
# Install dependencies
go mod download

# Generate GraphQL code (after schema changes)
gqlgen generate

# Build
go build .

# Run
go run .

# Server starts on http://localhost:8080
```

### Testing

```bash
# Run all tests
go test ./...

# Run specific test suites
go test ./tests/unit/...          # Unit tests only
go test ./tests/e2e/...           # E2E tests only
go test ./tests/performance/...   # Benchmarks only

# Run benchmarks with custom time
go test ./tests/performance/... -bench=. -benchtime=1s

# Verbose output
go test ./tests/e2e/... -v
```

### Endpoints

- **REST API**: `http://localhost:8080/`
  - Audience: `/audience`, `/audiences`
  - Chart: `/chart`, `/charts`
  - Insight: `/insight`, `/insights`
  - UserStar: `/userstar`, `/userstars`

- **GraphQL API**: `http://localhost:8080/graphql` (POST)
- **GraphQL Playground**: `http://localhost:8080/graphql` (GET - Interactive IDE)

## Performance Optimization Ideas 

### 1. Database 

#### [Indexing (AI did this in a different branch)](https://github.com/bsakel/platform-go-challenge/commit/809e44b323f616ffc3b4ec691ecf605aed691714) 

Ensure proper indexes exist:

```sql
CREATE INDEX idx_user_favourites_user_id ON user_favourites(user_id);
CREATE INDEX idx_user_favourites_asset_id ON user_favourites(asset_id);
CREATE INDEX idx_user_favourites_type ON user_favourites(type);
```

#### Query Optimization

Use `EXPLAIN ANALYZE` to identify slow queries:

```sql
EXPLAIN ANALYZE
SELECT * FROM user_favourites WHERE user_id = 1;
```

#### [Connection Pooling (AI did this in a different branch)](https://github.com/bsakel/platform-go-challenge/commit/89bb1bd103fb024f0a0dddb2172042be9d9aaa6c) 

Configure GORM connection pool in `db/db.go`:

```go
sqlDB, err := GormDB.DB()
if err == nil {
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
}
```

### 2. Caching 

For frequently accessed data, consider caching:

**Option 1: In-Memory Cache (Simple)**
```go
// 5-minute TTL cache
type CacheEntry struct {
    Data      *model.UserStared
    ExpiresAt time.Time
}
```

**Option 2: Redis Cache (Production)**
```go
// Cache key: "userstared:{userID}"
// TTL: 5-15 minutes
// Invalidate on mutations
```

### 3. [Concurency in the userstars query resolver (AI did this in a different branch)](https://github.com/bsakel/platform-go-challenge/commit/066d2f31ac45e082a457f079c0e88f5dfb5d9a82) 

In [userstared.resolvers.go](graph/resolvers/userstared.resolvers.go) the basic structure of operations is
	
  1. Fetch all user stars for this user
  2. Fetch all audiences 
  3. Fetch all charts 
  4. Fetch all insights 
  5. Build and return the UserStared response

Steps 2,3 and 4 are independent to each other and could be executed as group of goroutines and synced by a waitgroup.

### 4. Adapt web logic and GraphQL queries to allow for pagination

Another way to make the ui perform better is to add add logic to paginate the results instead of quering for all results at once.
