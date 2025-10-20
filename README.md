# GlobalWebIndex Engineering Challenge

# Favourites API

A small Go web service that demonstrates a favourites API:

- list a user's favourites
- add/remove an asset
- update an asset's description.

The project uses a concurrency-safe in-memory store and a simple JWT-based auth middleware for protected routes.

This README provides a concise developer guide:

- How to build, run and test locally
- API endpoints and example curl commands
- Notes on the in-memory store and where to add persistence

## Build

From the repository root:

```bash
go build ./...
```

## Run (dev)

Start the server (binds to :3002):

```bash
go run ./cmd/app
```

The server seeds an in-memory store on startup and prints two seeded user UUIDs you can use in examples.

## Tests

Run the unit tests:

```bash
go test ./...
```

Current tests cover the handler layer and in-memory store.

## Docker

A `Dockerfile` is provided that builds the Go binary inside a golang builder image and produces a runnable image that exposes port 3002. The current Dockerfile uses a single-stage build based on the official golang image and outputs a static binary named `favourites-app`.

Build the image:

```bash
docker build -t platform-go-challenge:latest .
```

Run the container (maps local 3002 to container 3002):

```bash
docker run --rm -p 3002:3002 platform-go-challenge:latest
```

## API

All protected endpoints require a Bearer JWT token. Use the token generation endpoint to get a token for a given user UUID. Note: the auth middleware requires the token's user (subject) to match the {userId} path parameter — otherwise requests will be rejected with 403 Forbidden.

1. Generate token

- POST /token
- Body: { "user_id": "<user-uuid>" }
- Response: { "token": "<jwt>" }

Example:

```bash
curl -s -X POST -H "Content-Type: application/json" \
	-d '{"user_id":"431ff79b-4e65-4625-9f7f-c8125e2be057"}' \
	http://localhost:3002/token
```

Use the returned token in the `Authorization: Bearer <token>` header for subsequent requests. The token subject must equal the {userId} path parameter on protected routes.

2. Add favourite

- POST /users/{userId}/favourites
- Header: Authorization: Bearer <token>
- Body: JSON representation of `model.Asset`. For convenience the tests and seed helpers use `model.NewAsset(...)` shapes. Minimal example:

```json
{
  "type": "insight",
  "description": "Top insight",
  "data": { "text": "40% of..." }
}
```

Example:

```bash
curl -X POST -H "Content-Type: application/json" \
	-H "Authorization: Bearer <token>" \
	-d '{"type":"insight","description":"Test desc","data":{"text":"hello"}}' \
	http://localhost:3002/users/<user-uuid>/favourites
```

Successful response: 201 Created with JSON { "message": "favourite added", "favourite": { ... } }

3. List favourites

- GET /users/{userId}/favourites
- Header: Authorization: Bearer <token>
- Query params (optional): offset, limit, sortBy (description|type), order (asc|desc)

Example:

```bash
curl -H "Authorization: Bearer <token>" \
	http://localhost:3002/users/<user-uuid>/favourites
```

Response: 200 OK { "favourites": [ ... ] }

4. Patch description

- PATCH /users/{userId}/favourites/{assetId}
- Header: Authorization: Bearer <token>
- Body: { "description": "New description" }
- Response: 204 No Content on success (current implementation sends a JSON message with the 204 status)

5. Delete favourite

- DELETE /users/{userId}/favourites/{assetId}
- Header: Authorization: Bearer <token>
- Response: 200 OK { "message": "favourite removed", "favouriteId": "<asset-id>" }

## Seeded users

The seed helper pre-populates two demo users. When you run the server you should see printed UUIDs similar to:

- 431ff79b-4e65-4625-9f7f-c8125e2be057
- d61733f6-910e-4e26-bb31-0cac3a7d4279

Use these UUIDs to generate tokens and try the endpoints.

## Extending persistence

The storage abstraction is defined by `internal/db.FavouriteRepository`. To add persistence (Postgres, Redis, etc.) implement that interface and replace the `db.NewMemoryStore()` call in `cmd/app/main.go` with your implementation.

Key methods to implement:

- GetAllFavourites(userID uuid.UUID, q PageQuery) ([]model.Asset, error)
- AddFavourite(userID uuid.UUID, a model.Asset) error
- RemoveFavourite(userID uuid.UUID, assetID uuid.UUID) error
- PatchFavourite(userID uuid.UUID, assetID uuid.UUID, description string) error

## Notes and caveats

- The store is in-memory and ephemeral. Restarting the server clears all runtime state.
- The JWT secret is hard-coded for the challenge and should be replaced by a secure secret management strategy for any real deployment.
- The `model.Asset.Data` field uses `any` to keep the example compact. For a production service a well-defined schema per asset type with validation is recommended.
