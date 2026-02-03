package rest

import (
	"encoding/json"
	"errors"
	"iter"
	"log/slog"
	"net/http"

	"go-favorites-app/internal/core/domain/favorites"
	"go-favorites-app/internal/core/ports"
)

type Handler struct {
	service ports.FavoriteService
	logger  *slog.Logger
}

func NewHandler(service ports.FavoriteService, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// Create handles POST /favorites
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok || userID == "" {
		h.respondError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	var req createAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	asset, err := parseAsset(req.Raw, req.Type, userID)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.service.Save(r.Context(), asset); err != nil {
		// Differentiate validation vs internal error
		h.logger.Error("failed to save asset", "error", err)
		if errors.Is(err, favorites.ErrValidation) {
			h.respondError(w, http.StatusBadRequest, err)
			return
		}
		h.respondError(w, http.StatusInternalServerError, err) // Or 400 if validation error
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(asset); err != nil {
		h.logger.Error("failed to write response", "error", err)
	}
}

// Get handles GET /favorites/{id}
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.respondError(w, http.StatusBadRequest, errors.New("missing id"))
		return
	}

	asset, err := h.service.FindByID(r.Context(), id)
	if err != nil {
		// Should check if not found
		h.logger.Error("failed to find asset", "id", id, "error", err)
		h.respondError(w, http.StatusNotFound, err)
		return
	}

	_ = json.NewEncoder(w).Encode(asset)
}

// List handles GET /favorites with streaming
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok || userID == "" {
		h.respondError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	// Use NewPagination helper
	p := NewPagination(r)

	iter, err := h.service.FindAllByUser(ctx, userID, p.Limit, p.Offset)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/x-ndjson")
	w.WriteHeader(http.StatusOK)

	// Stream response using NDJSON (Newline Delimited JSON)
	h.streamResponse(w, iter)
}

func (h *Handler) streamResponse(w http.ResponseWriter, iter iter.Seq2[favorites.Asset, error]) {
	enc := json.NewEncoder(w)
	for item, err := range iter {
		if err != nil {
			h.logger.Error("stream error", "err", err)
			return
		}
		if err := enc.Encode(item); err != nil {
			h.logger.Error("encode error", "err", err)
			return
		}
	}
}

// Delete handles DELETE /favorites/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok || userID == "" {
		h.respondError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	id := r.PathValue("id")
	if err := h.service.Delete(r.Context(), id, userID); err != nil {
		if err.Error() == "forbidden: you do not own this asset" {
			h.respondError(w, http.StatusForbidden, err)
			return
		}
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UpdateDescription handles PATCH /favorites/{id}
// Payload: {"description": "..."}
func (h *Handler) UpdateDescription(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok || userID == "" {
		h.respondError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	id := r.PathValue("id")
	var req struct {
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	asset, err := h.service.UpdateDescription(r.Context(), id, req.Description, userID)
	if err != nil {
		if err.Error() == "forbidden: you do not own this asset" {
			h.respondError(w, http.StatusForbidden, err)
			return
		}
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(asset); err != nil {
		h.logger.Error("failed to write response", "error", err)
	}
}

func (h *Handler) respondError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); err != nil {
		h.logger.Error("failed to write response", "error", err)
	}
}
