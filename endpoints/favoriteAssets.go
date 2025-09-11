package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"platform-go-challenge/data"
)

// Request structure (JSON from client)
type FavoriteRequest struct {
	UserId int `json:"user_id"`
}

type FavoriteResponse struct {
	UserId    int         `json:"user_id"`
	AssetData []AssetData `json:"assetData"`
	Message   string      `json:"msg"`
	Code      int         `json:"code"`
}

type AssetData struct {
	FavId           int    `json:"fav_id"`
	Description     string `json:"description"`
	Characteristics any    `json:"characteristics"`
}

func FavoriteAssets(w http.ResponseWriter, r *http.Request) {

	var req FavoriteRequest
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		resp := FavoriteResponse{
			Message: "Invalid JSON",
			Code:    200,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	fmt.Println("Selected user:", req.UserId)

	// Filter assets for this user
	var isUserExists bool = false
	var userAssets []AssetData

	for _, fav := range data.FavAsset {
		var characteristics any
		if fav.UserId == req.UserId {
			isUserExists = true
			for _, allAssets := range data.Assets {
				if allAssets.ID == fav.AssetId {
					if allAssets.Type == "Chart" {
						for _, ca := range data.ChartAssets {
							if ca.AssetId == fav.AssetId {
								characteristics = ca
								break
							}
						}
					} else if allAssets.Type == "Insight" {
						for _, ia := range data.InsightAssets {
							if ia.AssetId == fav.AssetId {
								characteristics = ia
								break
							}
						}
					} else {
						for _, aa := range data.AudienceAssets {
							if aa.AssetId == fav.AssetId {
								characteristics = aa
								break
							}
						}
					}
					userAssets = append(userAssets, AssetData{
						FavId:           fav.ID,
						Description:     allAssets.Description,
						Characteristics: characteristics,
					})
				}
			}
		}
	}

	resp := FavoriteResponse{
		UserId:    req.UserId,
		AssetData: userAssets,
	}

	if !isUserExists {
		resp.Message = "User has not any favorite asset"
		resp.Code = http.StatusOK
	} else {
		resp.Message = "Success"
		resp.Code = http.StatusOK
	}

	json.NewEncoder(w).Encode(resp)

}
