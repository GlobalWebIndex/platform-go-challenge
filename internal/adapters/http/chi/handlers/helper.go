package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// parseUUIDParam extracts a URL path parameter by key and validates it as a UUID.
// On error, it writes a 400 Bad Request response and returns (uuid.Nil, false).
func parseUUIDParam(writer http.ResponseWriter, req *http.Request, key string) (uuid.UUID, bool) {
	val := chi.URLParam(req, key)
	id, err := uuid.Parse(val)
	if err != nil {
		WriteJsonError(writer, "invalid `+key+`", http.StatusBadRequest)
		return uuid.Nil, false
	}
	return id, true
}

// parsePagination reads "limit" and "offset" query parameters, validates them,
// applies defaults (limit=20, max 50), and returns them.
// On invalid values, it writes a 400 Bad Request response and returns ok=false.
func parsePagination(writer http.ResponseWriter, req *http.Request) (limit, offset int, ok bool) {
	q := req.URL.Query()

	if val := q.Get("limit"); val != "" {
		n, err := strconv.Atoi(val)
		if err != nil || n < 0 {
			WriteJsonError(writer, "invalid limit", http.StatusBadRequest)
			return 0, 0, false
		}
		limit = n
	}
	if val := q.Get("offset"); val != "" {
		n, err := strconv.Atoi(val)
		if err != nil || n < 0 {
			WriteJsonError(writer, "invalid offset", http.StatusBadRequest)
			return 0, 0, false
		}
		offset = n
	}

	// Defaults and max bounds for limit.
	if limit == 0 {
		limit = 20
	} else if limit > 50 {
		limit = 50
	}

	return limit, offset, true
}

func WriteJsonError(writer http.ResponseWriter, msg string, status int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(map[string]string{"error": msg})
}
