package repository

import (
	"errors"

	"platform-go-challenge/models"
	"platform-go-challenge/utils"

	"gorm.io/gorm"
)

type Insightrepository struct{}

func (a Insightrepository) GetInsights(user_id int) ([]models.Insight, error) {
	result := []models.Insight{}
	err := utils.DB.Where("user_id = ?", user_id).Where("favourite = true").Find(&result).Find(&result).Error
	errors.Is(err, gorm.ErrRecordNotFound)

	if err != nil {
		return []models.Insight{}, err
	}

	return result, nil
}

func (a Insightrepository) GetInsightsPagination(user_id, limit, offset int) ([]models.Insight, error) {
	result := []models.Insight{}
	err := utils.DB.Where("user_id = ?", user_id).Where("favourite = true").Find(&result).Limit(limit).Offset(offset).Find(&result).Error

	errors.Is(err, gorm.ErrRecordNotFound)

	if err != nil {
		return []models.Insight{}, err
	}

	return result, nil
}

func (a Insightrepository) EditInsight(id int, favourite bool) (string, error) {
	//Find and update
	err := utils.DB.Model(&models.Insight{}).Where("id = ?", id).Update("favourite", favourite).Error
	errors.Is(err, gorm.ErrRecordNotFound)

	if err != nil {
		return "Error", err
	}

	return "Insight has been updated", err
}

func (a Insightrepository) AddInsight(insight models.Insight) (int, error) {
	err := utils.DB.Save(&insight).Error
	errors.Is(err, gorm.ErrRecordNotFound)

	if err != nil {
		return 0, err
	}

	return int(insight.ID), err
}
