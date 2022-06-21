package intetests

import (
	"platform-go-challenge/domain"
	"platform-go-challenge/sqldb"
	"testing"
)

func setupSuite(tb testing.TB) (*domain.Domain, func(tb testing.TB)) {
	db, err := sqldb.NewDB("user", "user", "127.0.0.1:3306", "mydb")
	if err != nil {
		tb.Fatal(err)
	}
	db.DropTablesIfExist()
	sqldb, _ := db.GormDB().DB()
	db.CreateTables()
	domain := domain.NewDomain(db)
	// Return a function to teardown the test
	return domain, func(tb testing.TB) {
		db.DropTablesIfExist()
		sqldb.Close()
	}
}
