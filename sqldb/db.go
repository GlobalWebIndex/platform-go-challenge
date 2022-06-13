package sqldb

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

func NewDB(username, pass, host, dbname string) (*DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, pass, host, dbname)
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return &DB{db: gormDB}, err
}

func (db *DB) createTables() {
	db.db.AutoMigrate(&User{})
	db.db.AutoMigrate(&Insight{})
	db.db.AutoMigrate(&Audience{})
	db.db.AutoMigrate(&Chart{})
	db.db.AutoMigrate(&FavouriteInsight{})
	db.db.AutoMigrate(&FavouriteChart{})
	db.db.AutoMigrate(&FavouriteAudience{})
}

func (db *DB) dropTablesIfExist() {
	mgt := db.db.Migrator()
	if mgt.HasTable(&User{}) {
		err := mgt.DropTable(&User{})
		if err != nil {
			log.Println("Error DB: ", err)
		}
	}
	if mgt.HasTable(&Insight{}) {
		err := mgt.DropTable(&Insight{})
		if err != nil {
			log.Println("Error DB: ", err)
		}
	}
	if mgt.HasTable(&Audience{}) {
		err := mgt.DropTable(&Audience{})
		if err != nil {
			log.Println("Error DB: ", err)
		}
	}
	if mgt.HasTable(&Chart{}) {
		err := mgt.DropTable(&Chart{})
		if err != nil {
			log.Println("Error DB: ", err)
		}
	}
	if mgt.HasTable(&FavouriteInsight{}) {
		err := mgt.DropTable(&FavouriteInsight{})
		if err != nil {
			log.Println("Error DB: ", err)
		}
	}
	if mgt.HasTable(&FavouriteChart{}) {
		err := mgt.DropTable(&FavouriteChart{})
		if err != nil {
			log.Println("Error DB: ", err)
		}
	}
	if mgt.HasTable(&FavouriteAudience{}) {
		err := mgt.DropTable(&FavouriteAudience{})
		if err != nil {
			log.Println("Error DB: ", err)
		}
	}
}
