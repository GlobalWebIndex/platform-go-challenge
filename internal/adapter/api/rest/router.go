package rest

import (
	"net/http"
)

// NewRouter initializes the HTTP router and registers routes.
func NewRouter(h *Handler, authH *AuthHandler, jwtSecret string, mws ...Middleware) http.Handler {
	mux := http.NewServeMux()

	// Auth Routes (Public)
	mux.HandleFunc("POST /signup", authH.SignUp)
	mux.HandleFunc("POST /login", authH.Login)

	// Public Routes
	// mux.HandleFunc("GET /favorites", h.List)  // Moved to protected
	// mux.HandleFunc("GET /favorites/{id}", h.Get) // Moved to protected

	// Protected Routes
	auth := AuthMiddleware(jwtSecret)

	mux.Handle("GET /favorites", auth(http.HandlerFunc(h.List)))
	mux.Handle("GET /favorites/{id}", auth(http.HandlerFunc(h.Get)))
	mux.Handle("POST /favorites", auth(http.HandlerFunc(h.Create)))
	// mux.Handle("GET /favorites/mine", auth(http.HandlerFunc(h.ListMine))) // Removed, redundant
	mux.Handle("DELETE /favorites/{id}", auth(http.HandlerFunc(h.Delete)))
	mux.Handle("PATCH /favorites/{id}", auth(http.HandlerFunc(h.UpdateDescription)))

	// Documentation
	mux.HandleFunc("GET /openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/openapi.yaml")
	})

	mux.HandleFunc("GET /api-docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := `<!DOCTYPE html>
				<html lang="en">
				<head>
					<meta charset="utf-8" />
					<meta name="viewport" content="width=device-width, initial-scale=1" />
					<meta name="description" content="SwaggerUI" />
					<title>SwaggerUI</title>
					<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
				</head>
				<body>
				<div id="swagger-ui"></div>
				<script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js" crossorigin></script>
				<script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-standalone-preset.js" crossorigin></script>
				<script>
					window.onload = () => {
						window.ui = SwaggerUIBundle({
							url: '/openapi.yaml',
							dom_id: '#swagger-ui',
							presets: [
								SwaggerUIBundle.presets.apis,
								SwaggerUIStandalonePreset
							],
							layout: "StandaloneLayout",
						});
					};
				</script>
				</body>
				</html>`
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(html))
	})

	// Wrap with middleware
	return Chain(mux, mws...)
}
