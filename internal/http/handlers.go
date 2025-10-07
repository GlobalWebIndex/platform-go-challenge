package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/SpBalis/platform-go-challenge/internal/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Favs  *service.FavouritesService
	Users *service.UsersService
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID := mustInt64(chi.URLParam(r, "id"))
	limit := atoiDefault(r.URL.Query().Get("limit"), 50)
	offset := atoiDefault(r.URL.Query().Get("offset"), 0)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	items, total, err := h.Favs.List(ctx, userID, limit, offset)
	if err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	if items == nil {
		items = []service.Favourite{}
	}

	//Pagination headers
	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	base := r.URL
	links := make([]string, 0, 3)

	//Current page
	q := base.Query()
	q.Set("limit", strconv.Itoa(limit))
	q.Set("offset", strconv.Itoa(offset))
	base.RawQuery = q.Encode()
	links = append(links, `<`+base.String()+`>; rel="self"`)

	//Next page
	nextOffset := offset + limit
	if nextOffset < total {
		q := base.Query()
		q.Set("limit", strconv.Itoa(limit))
		q.Set("offset", strconv.Itoa(nextOffset))
		base.RawQuery = q.Encode()
		links = append(links, `<`+base.String()+`>; rel="next"`)
	}

	//Previous page
	prevOffset := offset - limit
	if prevOffset < 0 {
		prevOffset = 0
	}
	if offset > 0 {
		q := base.Query()
		q.Set("limit", strconv.Itoa(limit))
		q.Set("offset", strconv.Itoa(prevOffset))
		base.RawQuery = q.Encode()
		links = append(links, `<`+base.String()+`>; rel="prev"`)
	}
	w.Header().Set("Link", strings.Join(links, ", "))

	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	userID := mustInt64(chi.URLParam(r, "id"))
	var in struct {
		Type        service.AssetType `json:"type"`
		Description string            `json:"description"`
		Data        any               `json:"data"`
		CustomDesc  *string           `json:"custom_description,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	a := service.Asset{Type: in.Type, Description: in.Description, Data: in.Data}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	fav, err := h.Favs.Add(ctx, userID, a, in.CustomDesc)
	if err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, fav)
}

func (h *Handler) Remove(w http.ResponseWriter, r *http.Request) {
	userID := mustInt64(chi.URLParam(r, "id"))
	assetID := mustInt64(chi.URLParam(r, "assetId"))
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := h.Favs.Remove(ctx, userID, assetID); err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) EditDesc(w http.ResponseWriter, r *http.Request) {
	userID := mustInt64(chi.URLParam(r, "id"))
	assetID := mustInt64(chi.URLParam(r, "assetId"))
	var in struct {
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := h.Favs.EditDescription(ctx, userID, assetID, in.Description); err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (h *Handler) ClearAll(w http.ResponseWriter, r *http.Request) {
	userID := mustInt64(chi.URLParam(r, "id"))
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := h.Favs.ClearAll(ctx, userID); err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Email == "" {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	id, err := h.Users.Create(ctx, in.Email)
	if err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"id": id, "email": in.Email})
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit := atoiDefault(r.URL.Query().Get("limit"), 50)
	offset := atoiDefault(r.URL.Query().Get("offset"), 0)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	out, err := h.Users.List(ctx, limit, offset)
	if err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	if out == nil {
		out = []service.User{}
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := mustInt64(chi.URLParam(r, "id"))
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	u, err := h.Users.Get(ctx, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, u)
}

func atoiDefault(s string, def int) int {
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return def
}

func mustInt64(s string) int64 { n, _ := strconv.ParseInt(s, 10, 64); return n }
