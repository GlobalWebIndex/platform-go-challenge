package gwiexercise

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	MaxBodyBytes = 1 << 20 // 1MB
	MaxBulkItems = 1000
)

type BadRequestError struct{ Msg string }

func (e BadRequestError) Error() string { return e.Msg }

type ConflictError struct{ Msg string }

func (e ConflictError) Error() string { return e.Msg }

type PayloadTooLargeError struct{ Msg string }

func (e PayloadTooLargeError) Error() string { return e.Msg }

type ValidationErrorDetail struct {
	Index   int    `json:"index"`
	Field   string `json:"field,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

type BulkValidationError struct {
	Details []ValidationErrorDetail `json:"details"`
}

func (e BulkValidationError) Error() string { return "validation failed" }

func mapErrorToHTTPStatus(err error) int {
	switch err.(type) {
	case PayloadTooLargeError:
		return http.StatusRequestEntityTooLarge
	case BadRequestError, BulkValidationError:
		return http.StatusBadRequest
	case ConflictError:
		return http.StatusConflict
	default:
		return http.StatusBadRequest
	}
}

type Server struct {
	store Store
}

func NewServer(store Store) *Server {
	return &Server{store: store}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	path := strings.Trim(r.URL.Path, "/")
	segs := strings.Split(path, "/")
	if len(segs) < 3 || segs[0] != "users" || segs[2] != "favourites" {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	userID := segs[1]

	if len(segs) == 3 {
		// /users/{userID}/favourites
		s.handleUserFavourites(w, r, userID)
		return
	}
	if len(segs) == 4 {
		if segs[3] == "bulk" {
			s.handleUserFavouritesBulk(w, r, userID)
			return
		}
		assetID := segs[3]
		s.handleUserFavouriteByID(w, r, userID, assetID)
		return
	}
	writeError(w, http.StatusNotFound, "not found")
}

func (s *Server) handleUserFavourites(w http.ResponseWriter, r *http.Request, userID string) {
	switch r.Method {
	case http.MethodGet:
		assets := s.store.List(userID)
		sortField, desc := parseSortParams(r)
		sortAssets(assets, sortField, desc)
		limit, offset := parseLimitOffset(r)
		assets = applyPagination(assets, limit, offset)
		writeJSON(w, http.StatusOK, assets)
	case http.MethodPost:
		body, err := readRequestBody(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		asset, err := decodeAddAsset(body)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		created, err := s.store.Add(userID, asset)
		if err != nil {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, created)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Server) handleUserFavouritesBulk(w http.ResponseWriter, r *http.Request, userID string) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	body, err := readRequestBody(r)
	if err != nil {
		writeAppError(w, err)
		return
	}
	assets, err := decodeBulkAddAssets(body)
	if err != nil {
		writeAppError(w, err)
		return
	}
	if err := s.preflightBatch(userID, assets); err != nil {
		writeAppError(w, err)
		return
	}
	created, err := s.addAssetsBatch(userID, assets)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func parseIntDefault(s string, def int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

func createdAtOf(a Asset) time.Time {
	switch v := a.(type) {
	case *ChartAsset:
		return v.CreatedAt
	case *InsightAsset:
		return v.CreatedAt
	case *AudienceAsset:
		return v.CreatedAt
	default:
		return time.Time{}
	}
}

func (s *Server) handleUserFavouriteByID(w http.ResponseWriter, r *http.Request, userID, assetID string) {
	switch r.Method {
	case http.MethodDelete:
		if err := s.store.Remove(userID, assetID); err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case http.MethodPatch:
		var req struct {
			Description string `json:"description"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json: "+err.Error())
			return
		}
		if strings.TrimSpace(req.Description) == "" {
			writeError(w, http.StatusBadRequest, "description is required")
			return
		}
		updated, err := s.store.UpdateDescription(userID, assetID, req.Description)
		if err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, updated)
	default:
		w.Header().Set("Allow", "PATCH, DELETE")
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		log.Printf("write json error: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func writeAppError(w http.ResponseWriter, err error) {
	status := mapErrorToHTTPStatus(err)
	if bve, ok := err.(BulkValidationError); ok {
		resp := map[string]any{
			"error":   bve.Error(),
			"details": bve.Details,
		}
		writeJSON(w, status, resp)
		return
	}
	writeError(w, status, err.Error())
}

func parseSortParams(r *http.Request) (field string, desc bool) {
	q := r.URL.Query()
	field = q.Get("sort")
	if field == "" {
		field = "created_at"
	}
	order := strings.ToLower(q.Get("order"))
	desc = order == "desc"
	return
}

func sortAssets(assets []Asset, field string, desc bool) {
	if len(assets) <= 1 {
		return
	}
	sort.Slice(assets, func(i, j int) bool {
		a, b := assets[i], assets[j]
		switch field {
		case "type":
			if desc {
				return string(a.GetType()) > string(b.GetType())
			}
			return string(a.GetType()) < string(b.GetType())
		case "created_at":
			fallthrough
		default:
			ai, aj := createdAtOf(a), createdAtOf(b)
			if desc {
				return ai.After(aj)
			}
			return ai.Before(aj)
		}
	})
}

func parseLimitOffset(r *http.Request) (limit, offset int) {
	q := r.URL.Query()
	limit = parseIntDefault(q.Get("limit"), 100)
	if limit <= 0 || limit > 1000 {
		limit = 100
	}
	offset = parseIntDefault(q.Get("offset"), 0)
	if offset < 0 {
		offset = 0
	}
	return
}

func applyPagination(assets []Asset, limit, offset int) []Asset {
	end := offset + limit
	if offset >= len(assets) {
		return []Asset{}
	} else if end < len(assets) {
		return assets[offset:end]
	} else if offset < len(assets) {
		return assets[offset:]
	}
	return assets
}

func readRequestBody(r *http.Request) ([]byte, error) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, BadRequestError{Msg: "invalid body"}
	}
	if len(b) > MaxBodyBytes {
		return nil, PayloadTooLargeError{Msg: "request body too large"}
	}
	return b, nil
}

func decodeAddAsset(body []byte) (Asset, error) {
	var req AddAssetRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, BadRequestError{Msg: "invalid json: " + err.Error()}
	}
	if err := req.ValidateBasic(); err != nil {
		return nil, BadRequestError{Msg: err.Error()}
	}
	asset, err := DecodeAssetFromAddRequest(req)
	if err != nil {
		return nil, BadRequestError{Msg: err.Error()}
	}
	return asset, nil
}

func decodeBulkAddAssets(body []byte) ([]Asset, error) {
	var reqs []AddAssetRequest
	if err := json.Unmarshal(body, &reqs); err != nil {
		return nil, BadRequestError{Msg: "invalid json array: " + err.Error()}
	}
	if len(reqs) > MaxBulkItems {
		return nil, PayloadTooLargeError{Msg: "too many items in bulk request: max " + strconv.Itoa(MaxBulkItems)}
	}
	assets := make([]Asset, 0, len(reqs))
	for i, r := range reqs {
		if err := r.ValidateBasic(); err != nil {
			return nil, BulkValidationError{Details: []ValidationErrorDetail{{Index: i, Message: err.Error(), Code: "invalid"}}}
		}
		a, err := DecodeAssetFromAddRequest(r)
		if err != nil {
			return nil, BulkValidationError{Details: []ValidationErrorDetail{{Index: i, Message: err.Error(), Code: "invalid"}}}
		}
		assets = append(assets, a)
	}
	return assets, nil
}

func (s *Server) preflightBatch(userID string, assets []Asset) error {
	existing := s.store.List(userID)
	existingIDs := make(map[string]struct{}, len(existing))
	for _, e := range existing {
		existingIDs[e.GetID()] = struct{}{}
	}
	batchIDs := make(map[string]struct{}, len(assets))
	for i, a := range assets {
		id := a.GetID()
		if id == "" {
			for {
				candidate := newID()
				if _, ok := existingIDs[candidate]; ok {
					continue
				}
				if _, ok := batchIDs[candidate]; ok {
					continue
				}
				setID(a, candidate)
				id = candidate
				break
			}
		} else {
			if _, ok := existingIDs[id]; ok {
				return ConflictError{Msg: "item " + strconv.Itoa(i) + ": asset with same id already exists for user"}
			}
			if _, ok := batchIDs[id]; ok {
				return ConflictError{Msg: "item " + strconv.Itoa(i) + ": asset with same id already exists for user"}
			}
		}
		batchIDs[id] = struct{}{}
	}
	return nil
}

func (s *Server) addAssetsBatch(userID string, assets []Asset) ([]Asset, error) {
	created := make([]Asset, 0, len(assets))
	for i, a := range assets {
		c, err := s.store.Add(userID, a)
		if err != nil {
			return nil, ConflictError{Msg: "item " + strconv.Itoa(i) + ": " + err.Error()}
		}
		created = append(created, c)
	}
	return created, nil
}
