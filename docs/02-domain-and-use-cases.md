# Domain & Use Cases – Assets Favourites API

## Domain

### Asset (polymorphic)

Common fields:

- `id` (string)
- `type` (enum: `chart`, `insight`, `audience`)
- `description` (string)

Type-specific data:

- Chart:

  - `title` (string)
  - `xAxis` (string)
  - `yAxis` (string)
  - `data` (simplified as `[]float64`)

- Insight:

  - `text` (string)

- Audience (simplified for the exercise):

  - `gender` (string: `male` | `female` | `other`)
  - `birthCountry` (string)
  - `ageGroup` (string, e.g. `"18-24"`)
  - `hoursOnSocialMediaPerDay` (int)
  - `purchasesLastMonth` (int)

### Favourite

A favourite represents a link between a user and an asset:

- `userId` (string)
- `assetId` (string)
- `createdAt` (timestamp)

For the challenge it can be assumed that:

- Assets are stored in an in-memory store.
- Favourites are stored in-memory keyed by `userId`.

**Asset ID Generation:**

- When creating assets inline (via POST /favourites), asset IDs are auto-generated.
- Format: UUID v4 (e.g., "550e8400-e29b-41d4-a716-446655440000").
- Asset IDs are immutable once assigned.
- The ID generation must be thread-safe to handle concurrent requests.

**Concurrency Considerations:**

- In-memory storage must be thread-safe to handle concurrent requests.
- Multiple requests for the same user should be handled safely (e.g., using `sync.RWMutex` or channels).
- Race conditions when adding/removing favourites must be prevented.
- Asset creation and favourite operations should be atomic where possible.

---

## Core Use Cases

1. **List user favourites**

   - Input: `user_id`
   - Output: list of favourites with resolved asset details.

2. **Add favourite**

   - Input: `user_id`, asset payload
   - Behavior:
     - The asset is created inline with a new auto-generated UUID.
     - It is added to the user's favourites.
     - Each call creates a new asset with a unique ID, so duplicates are impossible by design.

3. **Remove favourite**

   - Input: `user_id`, `asset_id`
   - Behavior:
     - The favourite entry for that user/asset pair is removed if it exists.
     - Success is returned even if already absent (idempotent).

4. **Edit favourite description**
   - Input: `user_id`, `asset_id`, new `description`
   - Behavior:
     - The asset's description for this favourite is updated.
     - If the favourite does not exist, 404 Not Found is returned.

---

## Out of Scope (for this challenge)

- **Real authentication and authorization** (but authentication middleware architecture is **in scope** - see design notes below).
- Real database, migrations, or persistence.
- Full search and filtering.
- Pagination (can be added later if needed).

## Design Considerations

### Authentication Middleware Architecture

While real authentication is not required, the design must include:

- **Authentication middleware interface/abstraction** that can be easily swapped
- **Mock JWT implementation** where JWT tokens are parsed from `Authorization: Bearer <token>` header, `user_id` is extracted from token claims, and it is validated that it matches the path parameter
- **Clear separation** between middleware and business logic
- **Easy replacement path** for real JWT validation without refactoring business logic

The middleware should:

- Parse JWT token from `Authorization: Bearer <token>` header
- Extract `user_id` from JWT claims (mock validation - no signature verification needed)
- Validate that JWT `user_id` matches the `userId` path parameter
- Inject user context into request context
- Be easily replaceable with minimal code changes
- Follow Go middleware patterns (e.g., `func(http.Handler) http.Handler`)
