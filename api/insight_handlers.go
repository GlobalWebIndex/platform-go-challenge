package api

import (
	"log"
	"net/http"

	"platform-go-challenge/api/model"
	"platform-go-challenge/db"
	"platform-go-challenge/models"

	"github.com/gin-gonic/gin"
)

func CreateInsight(c *gin.Context) {
	var insight models.Insight

	//bind the request body
	if err := c.ShouldBindJSON(&insight); err != nil {
		model.ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}
	db.GormDB.Create(&insight)
	model.ResponseJSON(c, http.StatusCreated, "Insight created successfully", insight)
}

func GetInsights(c *gin.Context) {
	if db.GormDB == nil {
		log.Fatal("DB pointer is nil")
	}

	var insights []models.Insight
	db.GormDB.Find(&insights)
	model.ResponseJSON(c, http.StatusOK, "Insights retrieved successfully", insights)
}

func GetInsight(c *gin.Context) {
	if db.GormDB == nil {
		log.Fatal("DB pointer is nil")
	}

	var insight models.Insight
	if err := db.GormDB.First(&insight, c.Param("id")).Error; err != nil {
		model.ResponseJSON(c, http.StatusNotFound, "Insight not found", nil)
		return
	}
	model.ResponseJSON(c, http.StatusOK, "Insight retrieved successfully", insight)
}

func UpdateInsight(c *gin.Context) {
	if db.GormDB == nil {
		log.Fatal("DB pointer is nil")
	}

	var insight models.Insight
	if err := db.GormDB.First(&insight, c.Param("id")).Error; err != nil {
		model.ResponseJSON(c, http.StatusNotFound, "Insight not found", nil)
		return
	}

	// bind the request body
	if err := c.ShouldBindJSON(&insight); err != nil {
		model.ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	db.GormDB.Save(&insight)
	model.ResponseJSON(c, http.StatusOK, "Insight updated successfully", insight)
}

func DeleteInsight(c *gin.Context) {
	if db.GormDB == nil {
		log.Fatal("DB pointer is nil")
	}

	var insight models.Insight
	if err := db.GormDB.Delete(&insight, c.Param("id")).Error; err != nil {
		model.ResponseJSON(c, http.StatusNotFound, "Insight not found", nil)
		return
	}
	model.ResponseJSON(c, http.StatusOK, "Insight deleted successfully", nil)
}
