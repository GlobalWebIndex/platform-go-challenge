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
	path := request.URL.Path
	pathWithoutPrefix := strings.TrimPrefix(path, "/users/")
	parts := strings.Split(pathWithoutPrefix, "/")

	if len(parts) != 2 {
		http.Error(writer, "invalid path", http.StatusBadRequest)
		return
	}

	userId := parts[0]

	switch request.Method {
	case "GET":
		handleGetFavorites(writer, request, userId)
	case "POST":
		handlePostFavorite(writer, request, userId)
	case "PUT":
		if len(parts) < 2 {
			http.Error(writer, "missing asset ID for PUT", http.StatusBadRequest)
			return
		}
		assetId := parts[1]
		handleUpdateFavorite(writer, request, userId, assetId)
	case "DELETE":
		if len(parts) < 2 {
			http.Error(writer, "missing asset ID for DELETE", http.StatusBadRequest)
			return
		}
		assetId := parts[1]
		handleRemoveFavorite(writer, request, userId, assetId)
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

func handleRemoveFavorite(writer http.ResponseWriter, request *http.Request, userId string, assetId string) {
	err := database.RemoveFavorite(userId, assetId)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
