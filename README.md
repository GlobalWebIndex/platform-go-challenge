# GlobalWebIndex Engineering Challenge

## Disclaimer
Most of the code was written by Junie Pro making prompts to GPT-5 model. As my application was for the role of Engineering Manager, my main focus was to provide the correct prompts to the tool and make sure that the delivery will be up to high standards.
Whenever it was necessary, I edited the code to make it more readable and understandable.

## Running

- Go 1.20+
- Start the server (in-memory store):
  - `go run .`
- Start with file persistence:
  - `FAV_STORE=file FAV_PATH=./data/favourites.json go run .`
  - By default, file store path is `data/favourites.json` if FAV_PATH is not set.

### Containerized (Docker/Podman)

- Build with Docker:
  - `docker build -t favourites:latest .`
- Run (in-memory store):
  - `docker run --rm -p 8080:8080 favourites:latest`
- Run with file persistence (mount host dir for data):
  - `docker run --rm -p 8080:8080 -e FAV_STORE=file -e FAV_PATH=/data/favourites.json -v $(pwd)/data:/data favourites:latest`

- Build with Podman:
  - `podman build -t favourites:latest -f Containerfile .`
- Run (in-memory store):
  - `podman run --rm -p 8080:8080 favourites:latest`
- Run with file persistence:
  - `podman run --rm -p 8080:8080 -e FAV_STORE=file -e FAV_PATH=/data/favourites.json -v $(pwd)/data:/data favourites:latest`

Note: On Windows PowerShell, replace $(pwd) with ${PWD} or a full path.

## API

Base path: `/users/{userID}/favourites`

- GET `/users/{userID}/favourites`
  - Query params:
    - `limit` (int, default 100, max 1000)
    - `offset` (int, default 0)
    - `sort` one of `created_at`, `type` (default `created_at`)
    - `order` `asc` or `desc` (default `asc`)
- POST `/users/{userID}/favourites`
  - Body: Add a single asset
  - Example:
    - `{"type":"chart","description":"My chart","payload":{"title":"Sales","xAxisTitle":"Month","yAxisTitle":"USD","data":[1,2,3]}}`
- POST `/users/{userID}/favourites/bulk`
  - Body: JSON array of the same objects as above to add in bulk
- PATCH `/users/{userID}/favourites/{assetID}`
  - Body: `{"description":"new description"}`
- DELETE `/users/{userID}/favourites/{assetID}`

Asset response objects contain common fields `id`, `type`, `description`, `createdAt` plus type-specific fields.

## Storage

Two store implementations:
- In-memory (default): fast, volatile.
- File-backed (JSON): persists to a single file. Configure via env vars:
  - `FAV_STORE=file`
  - `FAV_PATH=./data/favourites.json`

## Tests and Benchmarks

- Run tests: `go test ./...`
- Run benchmarks: `go test -bench=. -benchmem`

## API Docs (OpenAPI)

- OpenAPI (YAML): `docs/openapi.yaml`
- OpenAPI (JSON): `docs/openapi.json`
- You can view them with Swagger UI by pointing it to the JSON URL/file. For local dev, run the server and open Swagger UI with the file:// path or host the JSON via any static file server.

## Submission

Just create a fork from the current repo and send it to us!

Good luck, potential colleague!

## Benchmarks Explained

This repository includes four Go benchmark tests in gwiexercise_test/bench_test.go. Below is a description of how each benchmark works and the practical benefits of running them.

1) BenchmarkInMemoryAddList
- What it measures: The throughput and memory allocations of the in-memory store when repeatedly adding favourites for a single user and then listing them.
- How it works:
  - Creates a fresh in-memory store: NewInMemoryStore().
  - Uses the testing.B loop (for i := 0; i < b.N; i++) to Add a simple Insight asset for the same user on every iteration. b.N is dynamically chosen by the Go test runner to reach stable timing.
  - Calls b.ReportAllocs() so the benchmark also reports memory allocations/op.
  - Calls store.List(user) once after the add loop to exercise the read path after many writes.
- Why it’s useful:
  - Gives a quick signal about write performance and per-operation heap allocations for the core in-memory data structure.
  - Helps catch performance regressions when changing the store implementation (e.g., map/list handling, locking, copy patterns).
  - The final List() call helps ensure that the read path remains efficient as the number of stored items grows.

2) BenchmarkInFileAddList
- What it measures: The throughput and allocations when using the JSON file-backed store (filestore.go) to persist each write to disk, compared with the in-memory version.
- How it works:
  - Creates a temporary JSON file path and initializes the file store with NewFileStore(tempPath).
  - In the benchmark loop (for i := 0; i < b.N; i++), adds the same simple Insight asset for a single user. The file store serializes the asset set and writes to disk as needed.
  - Calls b.ReportAllocs() and then performs a single List() at the end to exercise the read path.
- Why it’s useful:
  - Quantifies the overhead of persistence (JSON encoding/decoding, file system I/O, OS buffering) relative to the in-memory store.
  - Helps you decide when to use the file-backed store in development or small deployments, and highlights optimizations in filestore.go.
  - Detects regressions related to disk writes, JSON structure, or lock contention under persistent storage.

3) BenchmarkInMemoryBulkEndpoint
- What it measures: End-to-end performance of the HTTP bulk insertion endpoint (/users/{userID}/favourites/bulk) when handling batches of 100 items using the in-memory store.
- How it works:
  - Creates an in-memory store and wires it into a new HTTP server with NewServer(store).
  - Builds a payload slice of 100 assets (type "insight").
  - Calls b.ResetTimer() to exclude the one-time setup cost from the timing.
  - In the benchmark loop, sends an HTTP POST request to the bulk endpoint for each iteration and asserts the response status is 201 Created.
- Why it’s useful:
  - Measures realistic, end-to-end cost (routing, JSON decoding, validation, store writes) rather than just the store layer.
  - Highlights the performance characteristics and allocations of the bulk ingestion code path under a predictable batch size (100 items).
  - Useful for sizing throughput, spotting regressions in the HTTP layer or JSON handling, and guiding optimizations such as request batching or reduced allocations.

4) BenchmarkInFileBulkEndpoint
- What it measures: End-to-end performance of the same HTTP bulk insertion endpoint when the server is backed by the JSON file store. This shows the additional overhead introduced by persistence (encoding and disk I/O) during bulk writes.
- How it works:
  - Creates a temporary JSON file path and initializes a file-backed store with NewFileStore(tempPath).
  - Wires the store into a new HTTP server with NewServer(store).
  - Builds the same 100-item payload as the in-memory version.
  - Calls b.ResetTimer() to exclude setup time.
  - In the benchmark loop, posts the batch to /users/{userID}/favourites/bulk and asserts a 201 Created response.
- Why it’s useful:
  - Lets you compare in-memory vs file-backed end-to-end throughput and allocations for identical workloads.
  - Surfaces the cost of JSON marshaling and file writes under bulk ingestion.
  - Helps identify regressions tied to persistence (e.g., JSON shape changes, file flush frequency, locking around save operations).

How to run and read the results
- Run all benchmarks with memory stats: go test -bench=. -benchmem
- Run a specific benchmark (e.g., file-backed bulk): go test -bench=InFileBulkEndpoint -benchmem
- Typical output line items:
  - ns/op: Average time per benchmark iteration (lower is better).
  - B/op and allocs/op: Bytes allocated and number of allocations per iteration (lower is better).
- Compare results across commits/branches to detect regressions or validate performance improvements.
- Compare InMemory vs InFile ns/op and allocs/op to understand the relative overhead of persistence on your machine and filesystem for both single-add and bulk endpoints.
