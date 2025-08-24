package controllers

import (
	"encoding/json"
	"fmt"
	"gwiEng/internal/api/models"
	"gwiEng/internal/database"
	"gwiEng/pkg/helpers"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)


func PatchFavouriteAsset(w http.ResponseWriter, r *http.Request, logger *zap.Logger, client *mongo.Client) error {
	var requestID string
	if requestID = r.Header.Get("X-Request-ID"); requestID == "" {
		requestID = uuid.New().String()
	}
	log.Println("Received a DEL request")

	userID := mux.Vars(r)["userId"]
	assetID := mux.Vars(r)["assetId"]

	// The filter to be used in the Patch functionality
	filter := bson.D{
		{Key: "userID", Value: userID},
		{Key: "assetID", Value: assetID},
	}

	var requestBody models.PatchFavouriteBody
	if err := helpers.DecodeJSON(r, &requestBody); err != nil {
		return helpers.RespondAndLogError(
			"PATCH user's favourite",
			"favourites",
			"",
			requestID,
			err,
			http.StatusInternalServerError,
			w,
			logger)
	}

	update := bson.M{
		"$set": bson.M{
			"description": requestBody.Desription,
		},
	}

	if err := database.UpdateOne(
		client, 
		"GwiDB",
		"favourites",
		filter,
		update); err != nil {
		log.Printf("error updating assetID = %v for userID = %v\n", assetID, userID)
		return helpers.RespondAndLogError("PATCH favourite", "favourites", assetID, requestID, err, http.StatusInternalServerError, w, logger)
	}

	response := helpers.MsgResponse{
		Msg: fmt.Sprintf("Asset with id = %v updated successfully for user = %v", assetID, userID),
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error occured marshalling response: %v\n", err)
		return nil
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)

	return nil
}








