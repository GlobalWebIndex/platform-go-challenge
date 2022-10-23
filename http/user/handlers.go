package user

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"platform-go-challenge/internal/app/assets"
	"platform-go-challenge/internal/app/users"
	"platform-go-challenge/internal/pagination"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type Handler struct {
	service       *users.Service
	authorization string
}

func NewUserHandler(userService *users.Service, authorization string) Handler {
	return Handler{
		service:       userService,
		authorization: authorization,
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

func (h *Handler) GetToken(w http.ResponseWriter, r *http.Request) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["admin"] = true
	claims["user_id"] = mux.Vars(r)["user_id"]
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token.Claims = claims
	tokenString, _ := token.SignedString([]byte(h.authorization))

	_ = json.NewEncoder(w).Encode(tokenString)
}

func (h *Handler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := mux.Vars(r)["user_id"]
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			fmt.Println("Malformed token")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Malformed Token"))
		} else {
			jwtToken := authHeader[1]
			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(h.authorization), nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && claims["user_id"] == userID {
				ctx := context.WithValue(r.Context(), "props", claims)
				// Access context values in handlers like this
				// props, _ := r.Context().Value("props").(jwt.MapClaims)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				fmt.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("Unauthorized"))
			}
		}
	})
}
