package repository_test

import (
	"os"
	"ownify_api/internal/repository"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func DBSetup() error {
	_, err := repository.NewTestDB()
	if err != nil {
		return err
	}
	//delete database
	deleteDBSql := "DROP DATABASE ownify"
	repository.DB.Exec(deleteDBSql)

	//create new database
	createDBSql := "CREATE DATABASE ownify"
	_, err = repository.DB.Exec(createDBSql)
	if err != nil {
		return nil
	}
	//select database
	useDBSql := "USE ownify"
	_, err = repository.DB.Exec(useDBSql)
	if err != nil {
		return nil
	}

	//create product_test table
	productTestSql := "CREATE TABLE products_test(chain_id int NOT NULL,asset_id BIGINT PRIMARY KEY COMMENT 'Primary Key',owner VARCHAR(255) NOT NULL,barcode VARCHAR(30) NOT NULL,item_name VARCHAR(255),brand_name VARCHAR(255),additional_data TEXT,location VARCHAR(50),issued_date BIGINT,transfer_date BIGINT,manufacture VARCHAR(50),origin_country VARCHAR(50),status VARCHAR(20),warranty BOOLEAN,warranty_time TIMESTAMP,insurance BOOLEAN,insurance_time TIMESTAMP,data_sharing INT,recyclable BOOLEAN,ecommerce BOOLEAN,category VARCHAR(30),price FLOAT,currency VARCHAR(10),points FLOAT,authenticity BOOLEAN,ownership BOOLEAN,is_gs1 BOOLEAN,create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Create Time',name VARCHAR(255)) COMMENT '';"
	_, err = repository.DB.Exec(productTestSql)
	if err != nil {
		return nil
	}

	return nil
}
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
func TestNewDB(t *testing.T) {
	_, err := repository.NewTestDB()
	require.Nil(t, err)
}
