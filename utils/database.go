package utils

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func ConnectPostgresDB() {
	DB, err = gorm.Open(postgres.Open(os.Getenv("ELEPHANTSQL_URL")), &gorm.Config{})
	// 	Logger: logger.Default.LogMode(logger.Info),
	// })
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to DB")
	}

}
