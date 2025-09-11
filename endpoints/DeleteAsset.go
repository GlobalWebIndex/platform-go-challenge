package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"platform-go-challenge/data"
)

// Request structure (JSON from client)
type DeleteAssetRequest struct {
	UserId  int `json:"user_id"`
	AssetId int `json:"asset_id"`
}

type DeleteAssetResponse struct {
	Success bool   `json:"success"`
	Message string `json:"msg"`
	Code    int    `json:"code"`
}

func DeleteAsset(w http.ResponseWriter, r *http.Request) {

	var req DeleteAssetRequest
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		resp := DeleteAssetResponse{
			Message: "Invalid JSON",
			Code:    200,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	var success_deletion bool = false

	for i := len(data.FavAsset) - 1; i >= 0; i-- {
		if data.FavAsset[i].AssetId == req.AssetId && data.FavAsset[i].UserId == req.UserId {
			data.FavAsset = append(data.FavAsset[:i], data.FavAsset[i+1:]...)
			success_deletion = true
			break
		}
	}

	fmt.Println("Favorite Assets after deletion are :", data.FavAsset)
	resp := DeleteAssetResponse{
		Success: success_deletion,
	}

	if success_deletion {
		resp.Message = "Success deletion"
		resp.Code = http.StatusOK
	} else {
		resp.Message = "No favorite asset found for the given attributes"
		resp.Code = http.StatusOK
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
