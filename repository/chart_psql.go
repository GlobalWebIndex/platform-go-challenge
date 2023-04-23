package repository

import (
	"errors"

	"gorm.io/gorm"
	"platform2.0-go-challenge/models"
	"platform2.0-go-challenge/utils"
)

type ChartRepository struct{}

func (a ChartRepository) GetCharts(user_id int) ([]models.Chart, error) {
	result := []models.Chart{}

	err := utils.DB.Where("user_id = ?", user_id).Where("favourite = true").Find(&result).Error

	errors.Is(err, gorm.ErrRecordNotFound)

	if err != nil {
		return []models.Chart{}, err
	}

	return result, nil
}

func (a ChartRepository) GetChartsPagination(user_id, limit, offset int) ([]models.Chart, error) {
	//var chart models.Chart
	result := []models.Chart{}
	err := utils.DB.Where("user_id = ?", user_id).Where("favourite = true").Limit(limit).Offset(offset).Find(&result).Error
	errors.Is(err, gorm.ErrRecordNotFound)

	if err != nil {
		return []models.Chart{}, err
	}

	return result, nil
}

func (a ChartRepository) EditChart(id int, favourite bool) (string, error) {
	//Find and update
	err := utils.DB.Model(&models.Chart{}).Where("id = ?", id).Update("favourite", favourite).Error
	errors.Is(err, gorm.ErrRecordNotFound)

	if err != nil {
		return "Error", err
	}

	return "Chart has been updated", err
}

func (a ChartRepository) AddChart(chart models.Chart) (int, error) {
	err := utils.DB.Save(&chart).Error
	errors.Is(err, gorm.ErrRecordNotFound)

	if err != nil {
		return 0, err
	}

	return int(chart.ID), err
}
