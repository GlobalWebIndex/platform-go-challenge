package api

import (
	"log"
	"net/http"

	"platform-go-challenge/api/model"
	"platform-go-challenge/db"
	"platform-go-challenge/models"

	"github.com/gin-gonic/gin"
)

func CreateUserStar(c *gin.Context) {
	var userstar models.UserStar

	//bind the request body
	if err := c.ShouldBindJSON(&userstar); err != nil {
		model.ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}
	db.GormDB.Create(&userstar)
	model.ResponseJSON(c, http.StatusCreated, "UserStar created successfully", userstar)
}

func GetUserStars(c *gin.Context) {
	if db.GormDB == nil {
		log.Fatal("DB pointer is nil")
	}

	var userstars []models.UserStar
	db.GormDB.Find(&userstars)
	model.ResponseJSON(c, http.StatusOK, "UserStars retrieved successfully", userstars)
}

func GetUserStar(c *gin.Context) {
	if db.GormDB == nil {
		log.Fatal("DB pointer is nil")
	}

	var userstar models.UserStar
	if err := db.GormDB.First(&userstar, c.Param("id")).Error; err != nil {
		model.ResponseJSON(c, http.StatusNotFound, "UserStar not found", nil)
		return
	}
	model.ResponseJSON(c, http.StatusOK, "UserStar retrieved successfully", userstar)
}

func UpdateUserStar(c *gin.Context) {
	if db.GormDB == nil {
		log.Fatal("DB pointer is nil")
	}

	var userstar models.UserStar
	if err := db.GormDB.First(&userstar, c.Param("id")).Error; err != nil {
		model.ResponseJSON(c, http.StatusNotFound, "UserStar not found", nil)
		return
	}

	// bind the request body
	if err := c.ShouldBindJSON(&userstar); err != nil {
		model.ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	db.GormDB.Save(&userstar)
	model.ResponseJSON(c, http.StatusOK, "UserStar updated successfully", userstar)
}

func DeleteUserStar(c *gin.Context) {
	if db.GormDB == nil {
		log.Fatal("DB pointer is nil")
	}

	var userstar models.UserStar
	if err := db.GormDB.Delete(&userstar, c.Param("id")).Error; err != nil {
		model.ResponseJSON(c, http.StatusNotFound, "UserStar not found", nil)
		return
	}
	model.ResponseJSON(c, http.StatusOK, "UserStar deleted successfully", nil)
}
