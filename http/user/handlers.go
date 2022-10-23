package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"platform-go-challenge/internal/app/assets"
	"platform-go-challenge/internal/app/users"

	"github.com/gorilla/mux"
)

type Handler struct {
	service *users.Service
}

func NewUserHandler(userService *users.Service) Handler {
	return Handler{
		service: userService,
	}
}

type starAction struct {
	AssetID     uint32           `json:"asset_id"`
	AssetType   assets.AssetType `json:"asset_type"`
	Description string           `json:"description"`
}

func (s *starAction) DecodeAndValidate(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		return err
	}
	if _, found := assets.AssetTypes[s.AssetType]; !found {
		return fmt.Errorf("invalid asset type: %v ", s.AssetType)
	}
	return nil
}

// GetDashboard
func (h *Handler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		http.Error(w, "wrong user id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	dashboard, err := h.service.GetDashboard(ctx, uint32(userID))
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	if dashboard.ID == 0 {
		http.Error(w, "entity not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(dashboard)
}

func (h *Handler) AddToDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		http.Error(w, "wrong user id", http.StatusBadRequest)
		return
	}
	action := starAction{}

	err = action.DecodeAndValidate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.service.AddToDashboard(ctx, uint32(userID), action.AssetID, action.AssetType, action.Description)
	if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RemoveFromDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		http.Error(w, "wrong user id", http.StatusBadRequest)
		return
	}
	action := starAction{}
	err = action.DecodeAndValidate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.service.RemoveFromDashboard(ctx, uint32(userID), action.AssetID, action.AssetType)
	if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) EditDescription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		http.Error(w, "wrong user id", http.StatusBadRequest)
		return
	}
	action := starAction{}
	err = action.DecodeAndValidate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.service.EditDescription(ctx, uint32(userID), action.AssetID, action.AssetType, action.Description)
	if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
