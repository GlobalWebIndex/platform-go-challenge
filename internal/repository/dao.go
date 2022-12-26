package repository

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Masterminds/squirrel"
	"github.com/spf13/viper"
)

type DAO interface {
	NewUserQuery() UserQuery
	NewTransactionQuery() TransactionQuery
}

type dao struct {
	db *sql.DB
}

var DB *sql.DB

func pgQb() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(DB)
}

func NewDAO(db *sql.DB) DAO {
	return &dao{db}
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

func (d *dao) NewTransactionQuery() TransactionQuery {
	return &transactionQuery{}
}

func (d *dao) NewUserQuery() UserQuery {
	return &userQuery{}
}
