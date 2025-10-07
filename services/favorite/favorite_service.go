package favorite_service

import (
	"encoding/json"
	asset_model "gwi/models/asset"
	favorite_model "gwi/models/favorite"
	"net/http"
	"strconv"

	"gwi/database"
	"gwi/utils"
)

// GetFavorites retrieves the list of favorite assets for the authenticated user.
// It extracts the user ID from the JWT token provided in the Authorization header,
// fetches the user's favorites from the database, and populates the asset details
// for each favorite. Finally, it sends the list of favorites back to the client
// as a JSON response. If the user is not authenticated, it sends an error response.
func GetFavorites(w http.ResponseWriter, r *http.Request) {
	// fetch user_id from JWT
	authorizationHeader := r.Header.Get("Authorization")
	userID, err := utils.ExtractUserID(authorizationHeader)
	if err != nil {
		utils.SendError(w, utils.ErrUnauthenticated("Invalid user"))
		return
	}

	var userFavorites []favorite_model.Favorite
	for _, fav := range database.FavoritesDB {
		if fav.UserID == userID {
			for _, a := range database.AssetsDB {
				if a.ID == fav.AssetID {
					fav.Asset = &a
					break
				}
			}
			userFavorites = append(userFavorites, fav)
		}
	}

	utils.SendSuccess(w, userFavorites)
}

// AddFavorite handles the HTTP request to add a favorite asset for a user.
// It decodes the JSON request body to extract the asset ID, checks if the asset exists,
// verifies the user's authentication, and ensures that the favorite does not already exist.
// If all checks pass, it creates a new favorite entry and adds it to the database.
func AddFavorite(w http.ResponseWriter, r *http.Request) {
	var input struct {
		AssetID int `json:"asset_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.SendError(w, utils.ErrBadRequest("Invalid JSON"))
		return
	}

	// Check if asset exists
	var assetFound *asset_model.Asset
	for _, a := range database.AssetsDB {
		if a.ID == input.AssetID {
			assetFound = &a
			break
		}
	}
	if assetFound == nil {
		utils.SendError(w, utils.ErrNotFound("Asset not found"))
		return
	}

	authorizationHeader := r.Header.Get("Authorization")
	userID, err := utils.ExtractUserID(authorizationHeader)
	if err != nil {
		utils.SendError(w, utils.ErrUnauthenticated("Invalid user"))
		return
	}
	newFavorite := favorite_model.Favorite{
		UserID:  userID,
		AssetID: input.AssetID,
		Asset:   assetFound,
	}
	// Check if favorite already exists
	for _, f := range database.FavoritesDB {
		if f.UserID == userID && f.AssetID == input.AssetID {
			utils.SendError(w, utils.ErrBadRequest("Favorite already exists"))
			return
		}
	}
	newFavorite.ID = len(database.FavoritesDB) + 1
	database.FavoritesDB = append(database.FavoritesDB, newFavorite)

	utils.SendSuccess(w, newFavorite)
}

// DeleteFavorite handles the HTTP request to delete a favorite item.
// It extracts the favorite ID from the URL, checks if the ID is valid,
// verifies if the user is authorized to delete the favorite,
// and removes the favorite from the database if found.
// If the ID is invalid, the user is unauthorized, or the favorite is not found,
// it sends an appropriate error response.
func DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/favorites/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, utils.ErrBadRequest("Invalid ID"))
		return
	}
	authorizationHeader := r.Header.Get("Authorization")
	_, err = utils.ExtractUserID(authorizationHeader)
	if err != nil {
		utils.SendError(w, utils.ErrUnauthenticated("Invalid user"))
		return
	}

	for i, f := range database.FavoritesDB {
		if f.ID == id {
			// Check if user is authorized to delete this favorite (owns it)
			if !(utils.IsUserAuthorized(f.UserID, r.Header.Get("Authorization"))) {
				utils.SendError(w, utils.ErrUnauthorized("User is not authorized for this action"))
				return
			}
			database.FavoritesDB = append(database.FavoritesDB[:i], database.FavoritesDB[i+1:]...)
			utils.SendSuccess(w, nil)
			return
		}
	}

	utils.SendError(w, utils.ErrNotFound("Favorite not found"))
}
