package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"platform-go-challenge/data"
)

// Request structure (JSON from client)
type EditAssetRequest struct {
	AssetId        int    `json:"asset_id"`
	NewDescription string `json:"new_description"`
}

type EditAssetResponse struct {
	Success bool   `json:"success"`
	Message string `json:"msg"`
	Code    int    `json:"code"`
}

func EditAssets(w http.ResponseWriter, r *http.Request) {

	var req EditAssetRequest
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		resp := EditAssetResponse{
			Message: "Invalid JSON",
			Code:    200,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	var success_change bool = false

	for i := range data.Assets {
		if data.Assets[i].ID == req.AssetId {
			data.Assets[i].Description = req.NewDescription
			success_change = true
			break
		}
	}

	fmt.Println("Assets are  :", data.Assets)
	resp := EditAssetResponse{
		Success: success_change,
		Code:    http.StatusOK,
	}

	if success_change {
		resp.Message = "Success edit of an asset"
	} else {
		resp.Message = "Failed to edit asset"
	}

	// Encode response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
