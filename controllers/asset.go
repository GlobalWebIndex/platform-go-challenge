package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"platform2.0-go-challenge/models"
	Repository "platform2.0-go-challenge/repository"
	"platform2.0-go-challenge/utils"
)

type Controller struct{}

var assets []models.Asset
var assetresp []models.AssetResponse

func (c Controller) GetUserAssets(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var error models.Error
		limit_param := r.URL.Query()["limit"]
		offset_param := r.URL.Query()["offset"]
		params := mux.Vars(r)

		assetresp = []models.AssetResponse{}
		assetRepo := Repository.AssetRepository{}

		id, _ := strconv.Atoi(params["user_id"])
		if len(limit_param) > 0 && len(offset_param) > 0 {
			limit, _ := strconv.Atoi(limit_param[0])
			offset, _ := strconv.Atoi(offset_param[0])
			row, err := assetRepo.GetUserAssetsPagination(id, limit, offset)

			if err != nil {
				if err == sql.ErrNoRows {
					error.Message = "Not Found"
					utils.SendError(w, http.StatusNotFound, error)
					return
				} else {
					error.Message = "Server error"
					utils.SendError(w, http.StatusInternalServerError, error)
					return
				}
			}
			w.Header().Set("Content-Type", "application/json")
			utils.SendSuccess(w, row)
		}
		row, err := assetRepo.GetUserAssets(id)

		if err != nil {
			if err == sql.ErrNoRows {
				error.Message = "Not Found"
				utils.SendError(w, http.StatusNotFound, error)
				return
			}
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, row)

	}
}

