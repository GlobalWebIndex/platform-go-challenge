package presentation

import (
	"encoding/json"
	"io"
	"net/http"
	"platform-go-challenge/domain"

	"github.com/gorilla/mux"
)

type FavouriteHandler struct {
	fs domain.FavouriteService
}

func NewFavouriteHandler(fs domain.FavouriteService) *FavouriteHandler {
	return &FavouriteHandler{fs: fs}
}

func (fh FavouriteHandler) AddToFavourites(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	assetID, ok := vars["asset_id"]
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		result, _ := json.Marshal("Asset ID is missing")
		writer.Write(result)
		return
	}

	userID := request.Header.Get("user_id")

	favourite, err := fh.fs.AddToFavourites(userID, assetID)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		result, _ := json.Marshal(err)
		writer.Write(result)
		return
	}

	writer.WriteHeader(http.StatusOK)
	result, _ := json.Marshal(favourite)
	writer.Write(result)
}

func (fh FavouriteHandler) GetFavourites(writer http.ResponseWriter, request *http.Request) {
	userID := request.Header.Get("user_id")

	assets, err := fh.fs.GetFavourites(userID)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		result, _ := json.Marshal(err)
		writer.Write(result)
		return
	}

	writer.WriteHeader(http.StatusOK)
	result, _ := json.Marshal(assets)
	writer.Write(result)
}

func (fh FavouriteHandler) EditFavourites(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	assetID, ok := vars["asset_id"]
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		result, _ := json.Marshal("Asset ID is missing")
		writer.Write(result)
		return
	}

	userID := request.Header.Get("user_id")

	body, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(500)
		result, _ := json.Marshal(err)
		writer.Write(result)
		return
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		writer.WriteHeader(500)
		result, _ := json.Marshal(err)
		writer.Write(result)
		return
	}

	asset, err := fh.fs.EditFavourites(userID, assetID, data["description"].(string))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		result, _ := json.Marshal(err)
		writer.Write(result)
		return
	}

	writer.WriteHeader(http.StatusOK)
	result, _ := json.Marshal(asset)
	writer.Write(result)
}

func (fh FavouriteHandler) DeleteFavourite(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	assetID, ok := vars["asset_id"]
	if !ok {
		http.Error(writer, "Asset ID is missing", http.StatusBadRequest)
		return
	}

	userID := request.Header.Get("user_id")

	_, err := fh.fs.DeleteFavourite(userID, assetID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
