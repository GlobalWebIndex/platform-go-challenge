## Favourites API

OpenAPI-first Go service for listing, adding, updating, and deleting user favourites across insights, audiences, and charts. Uses JWT auth (RS256), generated server/types via `oapi-codegen`, and in-memory adapters.

### Prerequisites
- Go 1.25+
- `oapi-codegen` installed (`go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.2.0`)
- `public.pem` / `private.pem` RSA keypair (see JWT section)

### Project layout
- `swagger.yaml` — OpenAPI contract (source of truth)
- `http/openapi.gen.go` — generated server/types (do not edit)
- `http/` — handlers implementing generated interfaces
- `favourites/application|adapters|domain` — favourites hexagon
- `assets/application|adapters|domain` — assets hexagon
- `cmd/favourites-api` — main server entrypoint

### Generate code from OpenAPI
```bash
oapi-codegen -config oapi-codegen.yaml swagger.yaml
```
This regenerates `http/openapi.gen.go`. Re-run after modifying `swagger.yaml`.

### Run locally
```bash
go run ./cmd/favourites-api          # uses configs/local.yaml (docs enabled)
# or
go run ./cmd/favourites-api --env local            # explicit env selection
# or
go test ./...   # runs unit + e2e tests (uses local private.pem/public.pem)
```
Config:
- `configs/<env>.yaml` merged with env vars (defaults -> file -> env). Default `APP_ENV=local`.
- Flags: `--env` overrides `APP_ENV`.
- Docs: set `enableDocs: true` (e.g., local) and open `http://localhost:8080/docs`.
- Key paths/issuer/audience can be overridden via env (`PUBLIC_KEY_PATH`, `PRIVATE_KEY_PATH`, `JWT_ISSUER`, `JWT_AUDIENCE`).

### Docker
```bash
docker build -t favourites-api .
docker run -p 8080:8080 -e APP_ENV=prod favourites-api
```
`Dockerfile` includes `public.pem` and the prod config only. Mount your own public key if you replace it. Never bake `private.pem` into the image.

### JWT keys and token generation
1. Generate RSA keypair (for local/dev):
   ```bash
   openssl genrsa -out private.pem 2048
   openssl rsa -in private.pem -pubout -out public.pem
   ```
2. Use `cmd/token-gen` to generate a token:
   ```bash
   go run ./cmd/token-gen --env local -user user-1   # flags override config
   ```
3. Use the token as `Authorization: Bearer <token>`.

### API usage
Base URL: `http://localhost:8080`
- `GET /health` — health check
- `GET /v1/favourites?page=&limit=` — list favourites (paginated)
- `POST /v1/favourites` — add favourite (`type`: INSIGHT|AUDIENCE|CHART, `assetId`, optional `description`)
- `PATCH /v1/favourites/{id}` — update description
- `DELETE /v1/favourites/{id}` — delete favourite
- `GET /v1/assets/{type}/{id}` — fetch owned asset (`type`: insights|audiences|charts)

All endpoints except `/health` & `/docs` and require a valid JWT signed with `private.pem` and matching issuer/audience from config.

### Tests
`go test ./...` runs unit tests plus E2E against in-memory adapters and real JWT signing/verification with local keys.
