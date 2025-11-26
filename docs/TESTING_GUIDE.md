# Testing Guide (AI Generated)

Comprehensive testing documentation for the platform-go-challenge project.

### Test Coverage

| Test Type | Location | Coverage |
|-----------|----------|----------|
| Unit | `tests/unit/` | AssetType enum validation, model logic |
| E2E | `tests/e2e/` | GraphQL queries with database integration |
| Benchmarks | `tests/performance/` | Query performance under different loads |

---

## Test Structure

```
tests/
├── unit/                         # Unit tests
│   └── userstar_test.go          # AssetType enum validation tests
├── e2e/                          # End-to-end integration tests
│   ├── setup_test.go             # Test infrastructure and helpers
│   └── userstared_test.go        # UserStared query functional tests
└── performance/                  # Performance benchmarks
    └── userstared_bench_test.go  # Benchmark tests for userstared query
```

---

## Setup

### Prerequisites

1. **Go 1.21+** installed
2. **PostgreSQL** database running
3. **Test database** created (optional but recommended)

### Environment Configuration

1. Create or verify `.env` file in project root:

```bash
# .env
DB_URL=postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable
TEST_DB_URL=postgres://myuser:mypassword@localhost:5432/mydb_test?sslmode=disable
```

**Note:** Tests automatically load `.env` from the project root. If `TEST_DB_URL` is not set, tests fall back to `DB_URL`.

2. Create test database (optional, for isolation):

```bash
createdb mydb_test
```

3. Install dependencies:

```bash
go mod download
```

---

## Running Tests

### Quick Commands

```bash
# Run all tests
go test ./...

# Run specific test suites
go test ./tests/unit/...          # Unit tests only
go test ./tests/e2e/...           # E2E tests only
go test ./tests/performance/...   # Benchmarks only

# Verbose output
go test ./tests/e2e/... -v

# Run specific test
go test ./tests/e2e/... -v -run TestUserStared_WithFavourites

# Run benchmarks with memory stats
go test ./tests/performance/... -bench=. -benchmem

# Run benchmarks with custom duration
go test ./tests/performance/... -bench=. -benchtime=10s
```

### Example Output

**Unit Tests:**
```
=== RUN   TestAssetType_IsValid
=== RUN   TestAssetType_IsValid/Valid_Audience
--- PASS: TestAssetType_IsValid (0.00s)
    --- PASS: TestAssetType_IsValid/Valid_Audience (0.00s)
PASS
ok      platform-go-challenge/tests/unit    0.002s
```

**E2E Tests:**
```
=== RUN   TestUserStared_EmptyFavourites
--- PASS: TestUserStared_EmptyFavourites (0.01s)
=== RUN   TestUserStared_WithFavourites
--- PASS: TestUserStared_WithFavourites (0.02s)
PASS
ok      platform-go-challenge/tests/e2e     0.162s
```

**Benchmarks:**
```
BenchmarkUserStared_SmallDataset-12      195    599114 ns/op
BenchmarkUserStared_MediumDataset-12      58   1835762 ns/op
BenchmarkUserStared_LargeDataset-12       20   7476081 ns/op
PASS
```

---

## Test Types

### 1. Unit Tests

**Location:** `tests/unit/`

**Purpose:** Fast, isolated validation of model logic and enums

**Current Coverage:**
- ✅ AssetType enum validation (`IsValid()` method)
- ✅ AssetType string conversion
- ✅ Database Value/Scan interface implementation
- ✅ Error handling for invalid types

**Run:**
```bash
go test ./tests/unit/... -v
```

**Test Count:** 24 tests covering all enum validation scenarios

---

### 2. E2E Tests

**Location:** `tests/e2e/`

**Purpose:** Integration testing with real database and GraphQL queries

**Current Coverage:**

| Test | Description |
|------|-------------|
| `TestUserStared_EmptyFavourites` | User with no starred items |
| `TestUserStared_WithFavourites` | User with one starred item per type |
| `TestUserStared_MultipleFavouritesOfSameType` | Multiple starred items of same type |
| `TestUserStared_OnlySpecificUser` | User isolation (only fetches correct user's data) |
| `TestUserStared_InvalidUserID` | Error handling for invalid user IDs |

**Run:**
```bash
go test ./tests/e2e/... -v
```

**Helper Functions Available:**
- `ExecuteGraphQL(t, query, variables)` - Execute GraphQL queries
- `CleanupTestData(testDB)` - Clean database before test
- `SeedTestData(t, testDB)` - Create sample test data

**Adding New E2E Tests:**

```go
func TestUserStared_YourNewScenario(t *testing.T) {
    // Arrange
    CleanupTestData(testDB)
    audienceID, chartID, insightID := SeedTestData(t, testDB)

    // Create test data
    testDB.Create(&models.UserStar{
        UserID:  1,
        Type:    models.AssetTypeAudience,
        AssetID: audienceID,
    })

    // Act
    query := `query { userstared(userID: "1") { userid audience { id } } }`
    resp := ExecuteGraphQL(t, query, nil)

    // Assert
    if len(resp.Errors) > 0 {
        t.Fatalf("expected no errors, got: %v", resp.Errors)
    }
    // ... more assertions
}
```

---

### 3. Benchmark Tests

**Location:** `tests/performance/`

**Purpose:** Measure query performance under different data loads

**Available Benchmarks:**

| Benchmark | Description | Data Size |
|-----------|-------------|-----------|
| `BenchmarkUserStared_SmallDataset` | Typical user load | 5 items per user |
| `BenchmarkUserStared_MediumDataset` | Average user load | 50 items per user |
| `BenchmarkUserStared_LargeDataset` | Heavy user load | 200 items per user |
| `BenchmarkUserStared_MinimalFields` | Minimal GraphQL query | 50 items, IDs only |
| `BenchmarkUserStared_AllFields` | Full GraphQL query | 50 items, all fields |
| `BenchmarkUserStared_Parallel` | Concurrent requests | 50 items/user, 100 users |

**Run Specific Benchmark:**
```bash
go test ./tests/performance/... -bench=BenchmarkUserStared_SmallDataset -benchmem
```

**Understanding Benchmark Output:**
```
BenchmarkUserStared_SmallDataset-12      195    599114 ns/op    123456 B/op    1234 allocs/op
```
- `195` - Number of iterations run
- `599114 ns/op` - Nanoseconds per operation (0.6ms)
- `123456 B/op` - Bytes allocated per operation (120 KB)
- `1234 allocs/op` - Number of allocations per operation

**Comparing Performance:**
```bash
# Before changes
go test ./tests/performance/... -bench=. -benchmem > old.txt

# After changes
go test ./tests/performance/... -bench=. -benchmem > new.txt

# Compare with benchstat
go install golang.org/x/perf/cmd/benchstat@latest
benchstat old.txt new.txt
```
