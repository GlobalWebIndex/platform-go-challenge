# Assets Favourites API – Context

This take-home implements the GWI Engineering Challenge.

Goal:

- A small Go web service is built that enables users to manage a personal list of favourite assets.
- Focus is placed on clean design (layers, tests, separation of concerns), not on a production-grade UI.

Constraints / assumptions:

- **API style:** A REST API is assumed for simplicity.
- A real database is **not** used; in-memory storage is sufficient for the exercise. The design should allow easy replacement with a persistent database in the future.
- The service should be fast enough to handle a large number of favourites per user.
- **Authentication:** No real authentication is required for the challenge; `user_id` is treated as a trusted input. However, authentication middleware **must be included** in the design (even if it's a mock/placeholder implementation) to demonstrate how the architecture supports easy replacement with real authentication later. This is a key design requirement for scalability.
- Assets can be of different types with different structures:
  - Chart
  - Insight
  - Audience

Non-functional goals for the implementation:

- Clear separation of layers (e.g. handler / service / repository / domain).
- Testable design with dependency injection.
- **Middleware architecture:** Authentication middleware should be designed as a pluggable component that can be easily swapped (e.g. mock JWT → real JWT validation) without changing business logic.
- Easy to extend later (e.g. add real DB, auth, pagination, search).

## Observability & Health

The service should expose basic observability endpoints and include structured logging:

### Logging

- **Structured logging:** Structured logging is used throughout the service for consistency and machine-readability
- **Implementation:** A structured logging library is used
- **Log levels:** Standard log levels are used (DEBUG, INFO, WARN, ERROR)
- **Key information logged:**
  - HTTP request/response details (method, path, status code, duration)
  - Error conditions with context
  - Authentication events
  - Business operations (favourites added, removed, updated)
- **Format:** JSON format is preferred for structured logging to enable easy parsing and integration with log aggregation systems
- **Keep it simple:** Focus is placed on essential logging; advanced features (sampling, correlation IDs) can be added later if needed

### Health Check

- **Health endpoint:** `GET /health` returns service health status (JSON)
- **Implementation:** The implementation is kept simple for this take-home assignment - `{"status": "healthy"}` is returned when the service is running
- **No authentication required:** Health endpoints are typically public for load balancer and orchestration system health checks
- **Note:** More advanced health probes (Liveness, Readiness, and Startup Probes) could be implemented later as separate endpoints, but for now they are out of scope for this exercise.

### Metrics

- **Metrics endpoint:** `GET /metrics` returns Prometheus-formatted metrics (text/plain)
- **Key metrics to expose:**
  - HTTP request counts by endpoint, method, and status code
  - HTTP request duration (latency) by endpoint and method
  - Error rates (4xx, 5xx responses)
  - Active operations (favourites added, removed, updated)
- **Implementation:** A metrics client library is used for standard metrics
- **Format:** Prometheus exposition format (text/plain)
- **Authentication:** For this take-home assignment, the `/metrics` endpoint uses the same authentication mechanism (`mockAuth`) as the rest of the API. **Note:** In production environments, metrics endpoints should be protected with stronger security measures (e.g., network-level restrictions, service mesh policies, or dedicated monitoring credentials) to prevent exposure of sensitive operational data. Metrics endpoints are often placed on separate ports or behind internal network boundaries.
- **Keep it simple:** Focus is placed on essential metrics; advanced metrics (percentiles, histograms) can be added later if needed
