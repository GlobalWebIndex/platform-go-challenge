package api

import (
	"log"
	"net/http"

	"platform-go-challenge/api/model"
	"platform-go-challenge/db"
	"platform-go-challenge/models"

	"github.com/gin-gonic/gin"
)

func CreateChart(c *gin.Context) {
	var chart models.Chart

	//bind the request body
	if err := c.ShouldBindJSON(&chart); err != nil {
		model.ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}
	db.GormDB.Create(&chart)
	model.ResponseJSON(c, http.StatusCreated, "Chart created successfully", chart)
}

func GetCharts(c *gin.Context) {
	if db.GormDB == nil {
		log.Fatal("DB pointer is nil")
	}

	var charts []models.Chart
	db.GormDB.Find(&charts)
	model.ResponseJSON(c, http.StatusOK, "Charts retrieved successfully", charts)
}

func GetChart(c *gin.Context) {
	if db.GormDB == nil {
		log.Fatal("DB pointer is nil")
	}

	var chart models.Chart
	if err := db.GormDB.First(&chart, c.Param("id")).Error; err != nil {
		model.ResponseJSON(c, http.StatusNotFound, "Chart not found", nil)
		return
	}
	model.ResponseJSON(c, http.StatusOK, "Chart retrieved successfully", chart)
}

func UpdateChart(c *gin.Context) {
	if db.GormDB == nil {
		log.Fatal("DB pointer is nil")
	}

	var chart models.Chart
	if err := db.GormDB.First(&chart, c.Param("id")).Error; err != nil {
		model.ResponseJSON(c, http.StatusNotFound, "Chart not found", nil)
		return
	}

	// bind the request body
	if err := c.ShouldBindJSON(&chart); err != nil {
		model.ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	db.GormDB.Save(&chart)
	model.ResponseJSON(c, http.StatusOK, "Chart updated successfully", chart)
}

func DeleteChart(c *gin.Context) {
	if db.GormDB == nil {
		log.Fatal("DB pointer is nil")
	}

	var chart models.Chart
	if err := db.GormDB.Delete(&chart, c.Param("id")).Error; err != nil {
		model.ResponseJSON(c, http.StatusNotFound, "Chart not found", nil)
		return
	}
	model.ResponseJSON(c, http.StatusOK, "Chart deleted successfully", nil)
}
