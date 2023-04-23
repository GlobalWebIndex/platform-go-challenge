package controllers

import (
	"encoding/json"
	"net/http"

	"reflect"

	"gorm.io/gorm"
	"platform2.0-go-challenge/models"
	Repository "platform2.0-go-challenge/repository"
	"platform2.0-go-challenge/utils"
)

type InsightController struct{}

func (i InsightController) UpdateInsight(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var insight models.Insight
		var error models.Error

		json.NewDecoder(r.Body).Decode(&insight)

		if reflect.ValueOf(insight.ID).IsZero() || reflect.ValueOf(insight.UserId).IsZero() {
			error.Message = "All fields are required."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}
		insightRepo := Repository.Insightrepository{}
		rowsUpdated, err := insightRepo.EditInsight(int(insight.ID), bool(insight.Favourite))

		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error) //500
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		utils.SendSuccess(w, rowsUpdated)
	}
}

func (i InsightController) AddInsight(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var insight models.Insight
		var error models.Error

		json.NewDecoder(r.Body).Decode(&insight)

		if reflect.ValueOf(insight.Text).IsZero() || reflect.ValueOf(insight.UserId).IsZero() {
			error.Message = "All fields are required."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}
		insightRepo := Repository.Insightrepository{}
		rowsUpdated, err := insightRepo.AddInsight(insight)

		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error) //500
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		utils.SendSuccess(w, rowsUpdated)
	}
}
