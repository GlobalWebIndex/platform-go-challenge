package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"platform-go-challenge/db"
	"platform-go-challenge/models"
)

// Global var to be used (like a static instance in Java)
var database = db.NewDB()

func FavoritesAPIs(writer http.ResponseWriter, request *http.Request) {
	urlParts := strings.Split(strings.Trim(request.URL.Path, "/"), "/")

	if len(urlParts) < 2 || urlParts[0] != "users" {
		http.Error(writer, "invalid path", http.StatusBadRequest)
		return
	}

	userId := urlParts[1]
	var assetId string
	if len(urlParts) > 2 {
		assetId = urlParts[2]
	}

	switch request.Method {
	case "GET":
		if assetId == "" {
			handleGetFavorites(writer, request, userId)
		} else {
			handleGetSingleFavorite(writer, request, userId, assetId)
		}
	case "POST":
		handlePostFavorite(writer, request, userId)
	case "PUT":
		if assetId == "" {
			http.Error(writer, "missing assetId for PUT", http.StatusBadRequest)
			return
		}
		handleUpdateFavorite(writer, request, userId, assetId)
	case "DELETE":
		if assetId == "" {
			http.Error(writer, "missing assetId for DELETE", http.StatusBadRequest)
			return
		}
		handleRemoveFavorite(writer, userId, assetId)
	default:
		http.Error(writer, fmt.Sprintf("Method %s not allowed", request.Method), http.StatusMethodNotAllowed)

	}
}

func handleGetFavorites(writer http.ResponseWriter, request *http.Request, userId string) {
	userFavorites, err := database.GetFavorites(userId)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(writer).Encode(userFavorites)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func handleGetSingleFavorite(writer http.ResponseWriter, request *http.Request, userId string, assetId string) {
	assets, err := database.GetFavorites(userId)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, asset := range assets {
		if asset.ID == assetId {
			writer.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(writer).Encode(asset)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}
	http.Error(writer, fmt.Sprintf("Asset %s not found", assetId), http.StatusNotFound)
}

func handlePostFavorite(writer http.ResponseWriter, request *http.Request, userId string) {
	var asset models.Asset

	err := json.NewDecoder(request.Body).Decode(&asset)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	asset.CreatedAt = time.Now()

	err = database.AddFavorite(userId, asset)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(writer).Encode(asset)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

}

func handleUpdateFavorite(writer http.ResponseWriter, request *http.Request, userId string, assetId string) {
	var asset models.Asset

	err := json.NewDecoder(request.Body).Decode(&asset)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	asset.ID = assetId

	err = database.UpdateFavorites(userId, asset)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	err = json.NewEncoder(writer).Encode(asset)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func handleRemoveFavorite(writer http.ResponseWriter, userId string, assetId string) {
	err := database.RemoveFavorite(userId, assetId)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
