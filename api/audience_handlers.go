package api

import (
	"log"
	"net/http"

	"platform-go-challenge/api/model"
	"platform-go-challenge/db"
	"platform-go-challenge/models"

	"github.com/gin-gonic/gin"
)

func CreateAudience(c *gin.Context) {
	var audience models.Audience

	//bind the request body
	if err := c.ShouldBindJSON(&audience); err != nil {
		model.ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}
	db.GormDB.Create(&audience)
	model.ResponseJSON(c, http.StatusCreated, "Audience created successfully", audience)
}

func GetAudiences(c *gin.Context) {
	if db.GormDB == nil {
		log.Fatal("DB pointer is nil")
	}

	var audiences []models.Audience
	db.GormDB.Find(&audiences)
	model.ResponseJSON(c, http.StatusOK, "Audiences retrieved successfully", audiences)
}

func GetAudience(c *gin.Context) {
	if db.GormDB == nil {
		log.Fatal("DB pointer is nil")
	}

	var audience models.Audience
	if err := db.GormDB.First(&audience, c.Param("id")).Error; err != nil {
		model.ResponseJSON(c, http.StatusNotFound, "Audience not found", nil)
		return
	}
	model.ResponseJSON(c, http.StatusOK, "Audience retrieved successfully", audience)
}

func UpdateAudience(c *gin.Context) {
	if db.GormDB == nil {
		log.Fatal("DB pointer is nil")
	}

	var audience models.Audience
	if err := db.GormDB.First(&audience, c.Param("id")).Error; err != nil {
		model.ResponseJSON(c, http.StatusNotFound, "Audience not found", nil)
		return
	}

	// bind the request body
	if err := c.ShouldBindJSON(&audience); err != nil {
		model.ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	db.GormDB.Save(&audience)
	model.ResponseJSON(c, http.StatusOK, "Audience updated successfully", audience)
}

func DeleteAudience(c *gin.Context) {
	if db.GormDB == nil {
		log.Fatal("DB pointer is nil")
	}

	var audience models.Audience
	if err := db.GormDB.Delete(&audience, c.Param("id")).Error; err != nil {
		model.ResponseJSON(c, http.StatusNotFound, "Audience not found", nil)
		return
	}
	model.ResponseJSON(c, http.StatusOK, "Audience deleted successfully", nil)
}
