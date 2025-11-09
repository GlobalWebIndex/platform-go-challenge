// internal/server/server.go
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"strconv"

	"github.com/KostasDasios/platform-go-challenge/internal/config"
	"github.com/KostasDasios/platform-go-challenge/internal/middleware"
	"github.com/KostasDasios/platform-go-challenge/internal/repo"
	"github.com/KostasDasios/platform-go-challenge/internal/service"
)

const (
    defaultLimit = 100
    maxLimit     = 1000
)

type Server struct {
	cfg     *config.Config
	svc     *service.Service
	mux     *http.ServeMux
	handler http.Handler // mux wrapped with middleware chain
}

// NewServer builds a Server with an in-memory repository.
// Swap NewInMemoryRepo with a persistent implementation without touching handlers.
func NewServer(cfg *config.Config) *Server {
	r := repo.NewInMemoryRepo()
	svc := service.NewService(r)

	mux := http.NewServeMux()
	s := &Server{cfg: cfg, svc: svc, mux: mux}
	s.routes()

	// allow Swagger UI on 8081 for local testing
    allowed := []string{"http://localhost:8081"}

	// Construct a lightweight rate limiter middleware based on environment config.
	// Default: ~20 requests/sec per user or IP (configurable via RATE_LIMIT_MS).
	rl := middleware.NewRateLimiter(time.Duration(cfg.RateLimitMillis) * time.Millisecond)

	// Middleware chain: security headers -> request id -> logger -> body limit -> rate limiter -> routes
	// MaxBody set to 1MB (configurable via env) for POST/PATCH payloads.
	s.handler = middleware.SecurityHeaders(
		middleware.CORS(allowed)(
			middleware.RequestID(
				middleware.Logger(
					middleware.MaxBody(cfg.MaxBodyBytes,
						rl.Middleware(
							middleware.APIKeyAuth(cfg.APIKey, s.mux),
						),
					),
				),
			),
		),
	)

	return s
}

// Handler exposes the fully wrapped HTTP handler (mux + middleware chain).
func (s *Server) Handler() http.Handler { return s.handler }

func (s *Server) routes() {
	// Liveness
	s.mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	// Readiness (for future external deps; always true for in-memory)
	s.mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"ready":true}`)
	})

	// REST endpoints:
	//   GET    /users/{userID}/favourites
	//   POST   /users/{userID}/favourites
	//   PATCH  /users/{userID}/favourites/{favID}
	//   DELETE /users/{userID}/favourites/{favID}
	s.mux.HandleFunc("/users/", s.routeUsers)
}

func (s *Server) routeUsers(w http.ResponseWriter, r *http.Request) {
	// Expected paths: /users/{uid}/favourites[/favID]
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 || parts[0] != "users" || parts[2] != "favourites" {
		http.NotFound(w, r)
		return
	}
	userID := parts[1]
	var favID string
	if len(parts) >= 4 {
		favID = parts[3]
	}

	switch r.Method {
	case http.MethodGet:
		if favID != "" {
			http.NotFound(w, r)
			return
		}
		s.handleList(w, r, userID)
	case http.MethodPost:
		if favID != "" {
			http.NotFound(w, r)
			return
		}
		s.handleCreate(w, r, userID)
	case http.MethodPatch:
		if favID == "" {
			http.NotFound(w, r)
			return
		}
		s.handlePatch(w, r, userID, favID)
	case http.MethodDelete:
		if favID == "" {
			http.NotFound(w, r)
			return
		}
		s.handleDelete(w, userID, favID)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (s *Server) handleList(w http.ResponseWriter, r *http.Request, userID string) {
    // Parse query params
    qs := r.URL.Query()

    // limit
    limit := defaultLimit
    if v := qs.Get("limit"); v != "" {
        if n, err := strconv.Atoi(v); err == nil && n > 0 {
            if n > maxLimit {
                n = maxLimit
            }
            limit = n
        }
    }

    // offset
    offset := 0
    if v := qs.Get("offset"); v != "" {
        if n, err := strconv.Atoi(v); err == nil && n >= 0 {
            offset = n
        }
    }

    list, err := s.svc.ListFavourites(userID)
    if err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
        return
    }

    // Safe slice
    start := offset
    if start > len(list) {
        start = len(list)
    }
    end := start + limit
    if end > len(list) {
        end = len(list)
    }
    page := list[start:end]

    writeJSON(w, http.StatusOK, map[string]any{
        "favourites": page,
        "total":      len(list),
        "limit":      limit,
        "offset":     offset,
    })
}


func (s *Server) handleCreate(w http.ResponseWriter, r *http.Request, userID string) {
	var payload struct {
		Asset json.RawMessage `json:"asset"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}
	f, err := s.svc.CreateFavourite(userID, payload.Asset)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, f)
}

func (s *Server) handlePatch(w http.ResponseWriter, r *http.Request, userID, favID string) {
	var payload struct {
		Description *string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Description == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "description is required"})
		return
	}
	upd, err := s.svc.UpdateFavouriteDescription(userID, favID, *payload.Description)
	if err != nil {
		status := http.StatusNotFound
		if err.Error() == "invalid path" {
			status = http.StatusBadRequest
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, upd)
}

func (s *Server) handleDelete(w http.ResponseWriter, userID, favID string) {
	if err := s.svc.DeleteFavourite(userID, favID); err != nil {
		status := http.StatusNotFound
		if err.Error() == "invalid path" {
			status = http.StatusBadRequest
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
