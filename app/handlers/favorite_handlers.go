package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"gwi-platform/database"
	"gwi-platform/models"
	"gwi-platform/utils"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func AddFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.UserFavorites
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	utils.InfoLogger.Println("AddFavoriteHandler:AddFavorite", req)

	if req.UserID == 0 || req.Favorites == nil || len(*req.Favorites) == 0 {
		http.Error(w, "User ID and at least one favorite are required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	db := database.GetDB()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to begin transaction: %v", err)
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "CALL InsertAndUpdateShard(?, ?, ?)")
	if err != nil {
		utils.ErrorLogger.Printf("Failed to prepare statement: %v", err)
		http.Error(w, "Failed to prepare statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	for _, fav := range *req.Favorites {
		_, err = stmt.ExecContext(ctx, req.UserID, fav.AssetID, string(fav.AssetType))
		if err != nil {
			utils.ErrorLogger.Printf("Failed to add favorite: %v", err)
			http.Error(w, fmt.Sprintf("Failed to add favorite: %v", err), http.StatusInternalServerError)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		utils.ErrorLogger.Printf("Failed to commit transaction: %v", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Favorites added successfully"})
}

func GetUserFavoritesDetailedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	utils.InfoLogger.Println("GetUserFavoritesDetailedHandler:Getting favorites for user_id", userIDStr)
	if userIDStr == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	/* adding paging */
	pageNumStr := r.URL.Query().Get("page_num")
	var shardID string
	if pageNumStr != "" {
		pageNum, err := strconv.ParseInt(pageNumStr, 10, 64)
		if err != nil {
			utils.ErrorLogger.Printf("Error parsing page_num: %v", err)
			http.Error(w, "Invalid page_num parameter", http.StatusBadRequest)
			return
		}
		shardID = fmt.Sprintf("%d_%d", userID, pageNum)
		utils.InfoLogger.Println("GetUserFavoritesDetailedHandler: Setting shardId", shardID)
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	db := database.GetDB()
	var shards []string

	if shardID != "" {
		shards = append(shards, shardID)
	} else {
		rows, err := db.QueryContext(ctx, sqlGetShards, userID)
		if err != nil {
			utils.ErrorLogger.Printf("Failed to retrieve user shards: %v", err)
			http.Error(w, "Failed to retrieve user shards", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var shard string
			if err := rows.Scan(&shard); err != nil {
				utils.ErrorLogger.Printf("Error scanning shards: %v", err)
				http.Error(w, "Error scanning shards", http.StatusInternalServerError)
				return
			}
			shards = append(shards, shard)
		}

		if err = rows.Err(); err != nil {
			utils.ErrorLogger.Printf("Error iterating shards: %v", err)
			http.Error(w, "Error iterating shards", http.StatusInternalServerError)
			return
		}
	}

	resultChan := make(chan []models.FavoriteAssetDetailed, len(shards))
	var wg sync.WaitGroup

	for _, shard := range shards {
		wg.Add(1)
		go func(shard string) {
			defer wg.Done()
			favorites, err := getFavoritesForShard(ctx, db, shard)
			if err != nil {
				utils.ErrorLogger.Printf("Error getting favorites for shard %s: %v", shard, err)
				return
			}
			resultChan <- favorites
		}(shard)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var allFavorites []models.FavoriteAssetDetailed
	for favorites := range resultChan {
		allFavorites = append(allFavorites, favorites...)
	}

	response := models.UserFavoritesDetailed{
		UserID:    models.ID(userID),
		Favorites: &allFavorites,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getFavoritesForShard(ctx context.Context, db database.DB, shard string) ([]models.FavoriteAssetDetailed, error) {

	utils.InfoLogger.Println("getFavoritesForShard:Getting favorites for shard", shard)

	rows, err := db.QueryContext(ctx, sqlGetFavorites, shard)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve favorites: %w", err)
	}
	defer rows.Close()

	var favorites []models.FavoriteAssetDetailed
	for rows.Next() {
		var fav models.FavoriteAssetDetailed
		var assetType models.AssetType
		var assetID models.ID
		var chartTitle, chartXAxis, chartYAxis, insightText, gender, birthCountry, ageGroup sql.NullString
		var dailyHours, purchasesLastMonth sql.NullFloat64

		err := rows.Scan(
			&fav.ID, &assetType, &assetID,
			&chartTitle, &chartXAxis, &chartYAxis,
			&insightText,
			&gender, &birthCountry, &ageGroup, &dailyHours, &purchasesLastMonth,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning favorites: %w", err)
		}

		switch assetType {
		case models.ChartAsset:
			chart := models.Chart{
				BaseAsset:  models.BaseAsset{AssetID: assetID, AssetType: models.ChartAsset},
				Title:      chartTitle.String,
				XAxisTitle: chartXAxis.String,
				YAxisTitle: chartYAxis.String,
			}
			dataPoints, err := getChartDataPoints(ctx, db, assetID)
			if err != nil {
				utils.ErrorLogger.Printf("Error getting data points for chart %d: %v", assetID, err)
			}
			chart.DataPoints = &dataPoints
			fav.Asset = chart
		case models.InsightAsset:
			fav.Asset = models.Insight{
				BaseAsset: models.BaseAsset{AssetID: assetID, AssetType: models.InsightAsset},
				Text:      insightText.String,
			}
		case models.AudienceAsset:
			fav.Asset = models.Audience{
				BaseAsset:               models.BaseAsset{AssetID: assetID, AssetType: models.AudienceAsset},
				Gender:                  gender.String,
				BirthCountry:            birthCountry.String,
				AgeGroup:                ageGroup.String,
				DailyHoursOnSocialMedia: dailyHours.Float64,
				PurchasesLastMonth:      int(purchasesLastMonth.Float64),
			}
		}
		favorites = append(favorites, fav)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return favorites, nil
}

func getChartDataPoints(ctx context.Context, db database.DB, chartID models.ID) ([]models.DataPoint, error) {

	utils.InfoLogger.Println("getChartDataPoints: Getting dataPoints for chart", chartID)

	rows, err := db.QueryContext(ctx, sqlGetDataPoints, chartID)
	if err != nil {
		return nil, fmt.Errorf("error querying data points: %w", err)
	}
	defer rows.Close()

	var dataPoints []models.DataPoint
	for rows.Next() {
		var dp models.DataPoint
		err := rows.Scan(&dp.X, &dp.Y)
		if err != nil {
			return nil, fmt.Errorf("error scanning data point: %w", err)
		}
		dataPoints = append(dataPoints, dp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	return dataPoints, nil
}
