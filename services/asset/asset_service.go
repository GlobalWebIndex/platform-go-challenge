package asset_service

import (
	"encoding/json"
	"gwi/database"
	asset_model "gwi/models/asset"
	"gwi/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetAssets handles the HTTP request to retrieve assets.
func GetAssets(w http.ResponseWriter, r *http.Request) {
	// This function currently does not implement any logic.
}

// AddAsset handles the HTTP request to add a new asset.
func AddAsset(w http.ResponseWriter, r *http.Request) {
	// This function currently does not implement any logic.
}

// DeleteAsset handles the HTTP request to delete an asset.
func DeleteAsset(w http.ResponseWriter, r *http.Request) {
	// This function currently does not implement any logic.
}

// EditAsset handles the HTTP request to edit an existing asset.
// It retrieves the asset ID from the URL parameters, decodes the
// JSON body for update data, checks user authorization, and
// merges the new data into the existing asset description.
// If successful, it updates the asset and sends a success response;
// otherwise, it sends an appropriate error response.
func EditAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	assetID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, utils.ErrBadRequest("Invalid asset ID"))
		return
	}

	authorizationHeader := r.Header.Get("Authorization")
	_, err = utils.ExtractUserID(authorizationHeader)
	if err != nil {
		utils.SendError(w, utils.ErrUnauthenticated("Invalid user"))
		return
	}

	// Decode request body into a generic map
	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		utils.SendError(w, utils.ErrBadRequest("Invalid JSON"))
		return
	}

	// Find asset in database
	var asset *asset_model.Asset
	for i := range database.AssetsDB {
		if database.AssetsDB[i].ID == assetID {
			// Check if user is authorized to edit this asset (owns it)
			if !(utils.IsUserAuthorized(database.AssetsDB[i].UserID, r.Header.Get("Authorization"))) {
				utils.SendError(w, utils.ErrUnauthorized("User is not authorized for this action"))
				return
			}
			asset = &database.AssetsDB[i]
			break
		}
	}

	if asset == nil {
		utils.SendError(w, utils.ErrNotFound("Asset not found"))
		return
	}

	// Unmarshal existing Description
	var existing map[string]interface{}
	if len(asset.Description) > 0 {
		if err := json.Unmarshal(asset.Description, &existing); err != nil {
			existing = make(map[string]interface{})
		}
	} else {
		existing = make(map[string]interface{})
	}

	// Merge new keys into existing JSON
	for k, v := range updateData {
		existing[k] = v
	}

	// Marshal back to json.RawMessage
	merged, err := json.Marshal(existing)
	if err != nil {
		utils.SendError(w, utils.ErrInternalServer("failed to update asset"))
		return
	}

	asset.Description = merged

	utils.SendSuccess(w, asset)
}
