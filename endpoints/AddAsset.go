package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"platform-go-challenge/data"
)

// Request structure (JSON from client)
type AddAssetRequest struct {
	UserId  int `json:"user_id"`
	AssetId int `json:"asset_id"`
}

type AddAssetResponse struct {
	Success bool   `json:"success"`
	Message string `json:"msg"`
	Code    int    `json:"code"`
}

func AddAsset(w http.ResponseWriter, r *http.Request) {

	var req AddAssetRequest
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		resp := AddAssetResponse{
			Message: "Invalid JSON",
			Code:    200,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Check if user exists and it's first asset that included, not new record required then.
	var lastItem int
	if len(data.FavAsset) > 0 {
		lastItem = data.FavAsset[len(data.FavAsset)-1].ID
	} else {
		lastItem = 0
	}
	fmt.Println("Favorite Assets are :", data.FavAsset[len(data.FavAsset)-1].ID)

	newAsset := data.FavAssets{
		ID:      lastItem + 1,
		UserId:  req.UserId,
		AssetId: req.AssetId,
	}

	data.FavAsset = append(data.FavAsset, newAsset)

	fmt.Println("Favorite Assets are :", data.Assets)
	resp := AddAssetResponse{
		Success: true,
		Message: "New asset added succesfully",
		Code:    http.StatusOK,
	}

	// Encode response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
