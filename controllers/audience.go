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

type AudienceController struct{}

func (a AudienceController) UpdateAudience(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var audience models.Audience
		var error models.Error

		json.NewDecoder(r.Body).Decode(&audience)

		if reflect.ValueOf(audience.ID).IsZero() || reflect.ValueOf(audience.UserId).IsZero() {
			error.Message = "All fields are required."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}
		audienceRepo := Repository.AudienceRepository{}
		rowsUpdated, err := audienceRepo.EditAudience(int(audience.ID), bool(audience.Favourite))

		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error) //500
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		utils.SendSuccess(w, rowsUpdated)
	}
}

func (a AudienceController) AddAudience(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var audience models.Audience
		var error models.Error

		json.NewDecoder(r.Body).Decode(&audience)

		if reflect.ValueOf(audience.UserId).IsZero() || reflect.ValueOf(audience.Country).IsZero() || reflect.ValueOf(audience.Gender).IsZero() || reflect.ValueOf(audience.Purchases).IsZero() || reflect.ValueOf(audience.AgeFrom).IsZero() || reflect.ValueOf(audience.AgeTo).IsZero() || reflect.ValueOf(audience.SocialHours).IsZero() {
			error.Message = "All fields are required."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}
		audienceRepo := Repository.AudienceRepository{}
		rowsUpdated, err := audienceRepo.AddAudience(audience)

		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error) //500
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		utils.SendSuccess(w, rowsUpdated)
	}
}
