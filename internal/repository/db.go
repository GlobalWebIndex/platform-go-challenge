package repository

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/spf13/viper"
)

type DBHandler interface {
	NewUserQuery() UserQuery
	NewBusinessQuery() BusinessQuery
	NewProductQuery() ProductQuery
	NewAdminQuery() AdminQuery
}

type dbHandler struct {
	db *sql.DB
}

var DB *sql.DB

func NewDBHandler(db *sql.DB) DBHandler {
	return &dbHandler{db}
}

func NewDB() (*sql.DB, error) {
	viper.AddConfigPath("../config")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("cannot read from a config")
	}
	host := viper.Get("database.host").(string)
	port := viper.Get("database.port").(string)
	user := viper.Get("database.user").(string)
	dbname := viper.Get("database.dbname").(string)
	password := viper.Get("database.password").(string)

	// Starting a database
	connection := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?parseTime=true"
	DB, err = sql.Open("mysql", connection)
	if err != nil {
		return nil, err
	}
	return DB, nil
}

func NewTestDB() (*sql.DB, error) {
	//viper.AddConfigPath("../config")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("cannot read from a config")
	}
	host := viper.Get("database_test_host").(string)
	port := viper.Get("database_test_port").(string)
	user := viper.Get("database_test_user").(string)
	dbname := viper.Get("database_test_dbname").(string)
	password := viper.Get("database_test_password").(string)

	// Starting a database
	connection := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?parseTime=true"
	DB, err = sql.Open("mysql", connection)
	if err != nil {
		return nil, err
	}
	return DB, nil
}

func (d *dbHandler) NewUserQuery() UserQuery {
	return &userQuery{}
}

func (d *dbHandler) NewBusinessQuery() BusinessQuery {
	return &businessQuery{}
}

func (d *dbHandler) NewProductQuery() ProductQuery {
	return &productQuery{}
}

func (d *dbHandler) NewAdminQuery() AdminQuery {
	return &adminQuery{}
}
