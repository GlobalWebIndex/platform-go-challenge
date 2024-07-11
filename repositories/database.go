package repositories

import (
	"fmt"
	"gwi-test/domain"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Db struct {
	*gorm.DB
}

func NewDbRepo() *Db {
	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	return &Db{
		db,
	}
}

func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to database")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running Mingrations")

	//Create the tables
	db.AutoMigrate(
		&domain.Audience{},
		&domain.Insight{},
		&domain.Chart{},
		&domain.ChartData{},
		&domain.Favourite{},
		&domain.Asset{},
	)

	return db, nil
}

func (db *Db) GetAssetsById(userId uint) ([]domain.Asset, error) {
	var assets []domain.Asset
	var favourites []domain.Favourite
	err := db.Model(&domain.Favourite{}).Where("user_id = ?", userId).Find(&favourites).Error
	if err != nil {
		return nil, err
	}
	for _, favourite := range favourites {
		fmt.Println("[LOG] looking for asset with id: ", favourite.AssetId)
		var asset domain.Asset
		err := db.First(&asset, favourite.AssetId).Error
		if err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}
	return assets, nil
}

func (db *Db) DeleteAsset(assetId uint) error {
	return db.Delete(&domain.Asset{}, assetId).Error
}

func (db *Db) AddAsset(asset *domain.Asset) error {
	return db.Create(asset).Error
}

func (db *Db) UpdateAsset(asset *domain.Asset) error {
	return db.Save(asset).Error
}

//-----Getters for Chart/Audience/Inisight

func (db *Db) GetChart(id uint) (*domain.Chart, error) {
	var chart domain.Chart
	result := db.First(&chart, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &chart, nil
}

func (db *Db) GetAudience(id uint) (*domain.Audience, error) {
	var audience domain.Audience
	result := db.First(&audience, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &audience, nil
}

func (db *Db) GetInsight(id uint) (*domain.Insight, error) {
	var insight domain.Insight
	result := db.First(&insight, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &insight, nil
}
