package controllers

import (
	"encoding/json"
	"fmt"
	"gwiEng/internal/api/models"
	"gwiEng/internal/database"
	"gwiEng/pkg/helpers"
	"gwiEng/pkg/logs"
	"log"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

func AddFavouriteAsset(w http.ResponseWriter, r *http.Request, logger *zap.Logger, client *mongo.Client) error {
	var requestID string
	if requestID = r.Header.Get("X-Request-ID"); requestID == "" {
		requestID = uuid.New().String()
	}
	log.Println("Received a POST request")

	
	var requestBody models.AddFavouriteBody
	if err := helpers.DecodeJSON(r, &requestBody); err != nil {
		return helpers.RespondAndLogError(
			"POST user's favourite",
			"favourites",
			"",
			requestID,
			err,
			http.StatusBadRequest,
			w,
			logger)
	}
	

	userID := mux.Vars(r)["userId"]
	assetID := mux.Vars(r)["assetId"]

	// Check that the asset ID provided exists in the assets collection
	filter := bson.D{
		{Key: "_id", Value: assetID},
	}

	var asset bson.M
	if err := database.FetchUsingFilter(
		client, 
		"GwiDB",
		"assets",
		filter,
		&asset,
		"asset"); err != nil {
		if strings.Contains(err.Error(), "no documents") {
			log.Printf("Err checking asset: %v exists --> %v", assetID, err)
			return helpers.RespondAndLogError(
				"POST", 
				"user favourite", 
				assetID, 
				requestID, 
				fmt.Errorf("%s:%v", assetID, err), 
				http.StatusNotFound, 
				w, logger)
		}

		log.Printf("Err checking asset: %v exists --> %v", assetID, err)
		http.Error(w, "Error fetching asset Id = "+fmt.Sprintf("%s: ", assetID)+err.Error(), http.StatusInternalServerError)
		errMsg := logs.ErrLogMsg{
			Action:     "PUT",
			Collection: "food",
			ResourceID: assetID,
			RequestID:  requestID,
			Error:      fmt.Errorf("food with id = %v not found", assetID),
		}
		logs.LogErrorAction(errMsg, logger)
		return nil
	}

	// Check that the user ID exists in the users collection
	usersFilter := bson.D{
		{Key: "userID", Value: userID},
	}
	var user bson.M
	if err := database.FetchUsingFilter(
		client,
		"GwiDB",
	    "users",
		usersFilter,
		&user,
		"user"); err != nil {

		if strings.Contains(err.Error(), "no documents") {
			log.Printf("Err checking user: %v exists --> %v", userID, err)
			http.Error(w, "Error fetching userID = "+fmt.Sprintf("%s: ", userID)+err.Error(), http.StatusNotFound)
			errMsg := logs.ErrLogMsg{
				Action:     "POST",
				Collection: "favourite asset",
				ResourceID: assetID,
				RequestID:  requestID,
				Error:      fmt.Errorf("asset with id = %v not found", userID),
				
			}
			logs.LogErrorAction(errMsg, logger)
			return nil
		}
		
		log.Printf("Err checking user: %v exists --> %v", userID, err)
		http.Error(w, "Error fetching user Id = "+fmt.Sprintf("%s: ", userID)+err.Error(), http.StatusInternalServerError)
		errMsg := logs.ErrLogMsg{
			Action:     "PUT",
			Collection: "users",
			ResourceID: userID,
			RequestID:  requestID,
			Error:      fmt.Errorf("user with id = %v not found", userID),
		}
		logs.LogErrorAction(errMsg, logger)
		return nil

	}

	// Check that the combination does not already exists in the db
	
	filter = bson.D{
		{Key: "userID", Value: userID},
		{Key: "assetID", Value: assetID},
	}

	var object bson.M
	if err := database.FetchUsingFilter(
		client,
		"GwiDB",
		"favourites",
		filter,
		&object,
		"asset"); err != nil {
		if !strings.Contains(err.Error(), "no documents") {
			log.Printf("Favourite asset: %v exists --> %v", assetID, err)
			return helpers.RespondAndLogError(
				"POST", 
				"user favourite", 
				assetID, 
				requestID, 
				fmt.Errorf("%s:%v", assetID, err), 
				http.StatusBadRequest, 
				w, logger)
		}
	}

	spew.Dump(object)
	if len(object) != 0 {
			return helpers.RespondAndLogError(
				"POST", 
				"user favourite", 
				assetID, 
				requestID, 
				fmt.Errorf("%s:%s already exists", assetID, userID), 
				http.StatusBadRequest, 
				w, logger)
	}




	var favourite  = models.Favourites{
		UserID: userID,
		AssetID: assetID,
		AssetType: requestBody.AssetType,
		Description: requestBody.Description,
	}

	// Insert the new favourite asset for the user
	if err := database.InsertOne(
		client,
		"GwiDB",
		"favourites",
		favourite,
		"favourite",
	); err != nil {
		log.Printf("Error inserting favourite asset with id = %v: %v\n", assetID, err)
		return helpers.RespondAndLogError("POST dashboard", "recipients", assetID, requestID, fmt.Errorf("asset = %v insertion error: %v", assetID, err), http.StatusInternalServerError, w, logger)
	}

	response := helpers.MsgResponse{
		Msg: fmt.Sprintf("Asset with id = %v registered successfully as favourite ", assetID),
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Error occured: %v\n", err)
		return nil
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)

	return nil

}



