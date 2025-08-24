// Package controllers contains the logic behind the handlers 
package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"gwiEng/internal/api/models"
	"gwiEng/internal/database"
	"gwiEng/pkg/helpers"
	"gwiEng/pkg/logs"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

func GetUserFavourites(w http.ResponseWriter, r *http.Request, logger *zap.Logger, client *mongo.Client) error {
	var requestID string
	if requestID = r.Header.Get("X-Request-ID"); requestID == "" {
		requestID = uuid.New().String()
	}

	queryParams := r.URL.Query()
	pageNumber := queryParams.Get("page")
	pageSize := queryParams.Get("pageSize")

	log.Printf("Controller GetUserFavourites: page number = %v and page size = %v\n", pageNumber, pageSize)

	pageNumInt, pageSizeInt, err := helpers.CheckPageParams(pageNumber, pageSize)
	if err != nil {
		http.Error(w, "Error in check Page Params: "+err.Error(), http.StatusBadRequest)
		errMsg := logs.ErrLogMsg{
			Action:     "GET",
			Collection: "food",
			ResourceID: "",
			RequestID:  requestID,
			Error:      err,
		}
		logs.LogErrorAction(errMsg, logger)
		return nil
	}

	offset, limit := helpers.CalculateOffsetAndLimit(pageNumInt, pageSizeInt)

	log.Println("offset = ", offset)
	println("limit = ", limit)

	// Parse user's ID from path parameter
	userID := mux.Vars(r)["userId"]
	log.Println("user's ID = ", userID)


	 filter := bson.D{
	 	{Key: "userID", Value: userID},
	 }

	totalFavs, err := database.CountDocuments(client, "GwiDB", "favourites", filter)
	if totalFavs == 0 {
		log.Println("err: ", errors.New("no favourites for user found"))
		err := errors.New("no favourites found")
		return helpers.RespondAndLogError(
			"GET user's favourites",
			"favourites",
			userID,
			requestID,
			err,
			http.StatusNotFound,
			w,
			logger)
	}

	if err != nil {
		log.Println("err: ", err)
		http.Error(w, "Internal Error: "+err.Error(), http.StatusInternalServerError)
		errMsg := logs.ErrLogMsg{
			Action:     "GET",
			Collection: "userFavourites",
			ResourceID: userID,
			RequestID:  requestID,
			Error:      err,
		}
		logs.LogErrorAction(errMsg, logger)
		return nil
	}
	log.Printf("total favourites = %v with user id = %v", totalFavs, userID)

	if offset >= int(totalFavs) {
		http.Error(w, fmt.Sprintf("Page requested = %v in not valid, total documents = %v", pageNumInt, totalFavs), http.StatusBadRequest)
		errMsg := logs.ErrLogMsg{
			Action:     "GET",
			Collection: "userFavourites",
			ResourceID: userID,
			RequestID:  requestID,
			Error:      fmt.Errorf("page offset = %v >= total Food = %v", offset, totalFavs),
		}
		logs.LogErrorAction(errMsg, logger)
		return nil
	}

	//var favs []bson.D
	favs := &[]models.Favourites{}
	if err := database.FetchSkipLimit(
		client,
		"GwiDB",
		"favourites",
		offset,
		limit,
		filter,
		favs,
	); err != nil {
		log.Printf("Err in Fetch Skip Limit = %v", err)
		http.Error(w, "Error fetching  =  "+err.Error(), http.StatusInternalServerError)
		errMsg := logs.ErrLogMsg{
			Action:     "GET",
			Collection: "favourites",
			ResourceID: userID,
			RequestID:  requestID,
			Error:      err,
		}
		logs.LogErrorAction(errMsg, logger)
		return nil
	}

	var assetIDs []string

	// create an array of assetID to retrieve them from the DB
	for _,favourite := range *favs {
		assetIDs = append(assetIDs, favourite.AssetID)
	}

	spew.Dump(favs)
	spew.Dump(assetIDs)

	var favAssets  []bson.M
	if err := database.FetchAssetsByIDs(
		client, 
		"GwiDB",
		"assets",
		assetIDs,
		&favAssets,
	); err != nil {
		log.Printf("Err in Fetching assets = %v", err)
		http.Error(w, "Error fetching assets  =  "+err.Error(), http.StatusInternalServerError)
		errMsg := logs.ErrLogMsg{
			Action:     "GET",
			Collection: "favourites",
			ResourceID: userID,
			RequestID:  requestID,
			Error:      err,
		}
		logs.LogErrorAction(errMsg, logger)
		return nil
	}

	log.Printf("----- Favourite Assets ------")
	spew.Dump(favAssets)

	// Add description to each asset
	for i, fav := range *favs {
		favAssets[i]["description"] = fav.Description
	}

	//jsonResp, err := json.Marshal(response)
	jsonResp, err := json.MarshalIndent(favAssets, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)

	return nil
	
}






















