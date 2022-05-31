package sqldb

import (
	"database/sql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

func NewDB() (*DB, error) {
	sqlDB, err := sql.Open("mysql", "mydb_dsn")
	if err != nil {
		return nil, err
	}
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	return &DB{db: gormDB}, err
}
