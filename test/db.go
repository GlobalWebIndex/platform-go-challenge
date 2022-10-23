package test

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql" // load mysql driver
)

// NewDBClient returns a new sql client.
func NewDBClient() *sql.DB {
	cs := "root:secret@tcp(localhost:3308)/challenge"

	client, _ := sql.Open("mysql", cs+"?loc=Europe%2FAthens&parseTime=true")
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(0)

	return client
}

// NewDBConnection returns a new bpsql Client.
func NewDBConnection() *sql.DB {
	return NewDBClient()
}

//
// DBInit creates a database connection executes the given tearUp function in it and returns a function
// that executes the tearDown function and closes the connection (to be used in defer).
//
func DBInit(tearUp func(db *sql.DB), tearDown func(db *sql.DB)) (*sql.DB, func()) {
	db := NewDBConnection()
	tearUp(db)
	return db, func() {
		tearDown(db)
		db.Close()
	}
}

//
// DBInitQueries creates a database connection executes the given tearUp query in it and returns a function
// that executes the tearDown query and closes the connection (to be used in defer).
//
func DBInitQueries(tearUp string, tearDown string) (*sql.DB, func()) {
	return DBInit(
		func(db *sql.DB) {
			_, _ = db.Exec(tearUp)
		},
		func(db *sql.DB) {
			_, _ = db.Exec(tearDown)
		},
	)
}

// DBInitQueriesWithTest function.
func DBInitQueriesWithTest(t *testing.T, tearUp []string, tearDown []string) (*sql.DB, func()) {
	t.Helper()
	db := NewDBConnection()
	down := func(db *sql.DB) {
		t.Helper()
		for _, q := range tearDown {
			_, err := db.Exec(q)
			if err != nil {
				t.Errorf("DB down failed. \n\tquery: %s \n\terror: %s", q, err.Error())
			}
		}
	}

	up := func(db *sql.DB) {
		t.Helper()
		atLeastOneFailed := false
		for _, q := range tearUp {
			_, err := db.Exec(q)
			if err != nil {
				atLeastOneFailed = true
				t.Errorf("DB init failed. \n\tquery: %s \n\terror: %s", q, err.Error())
			}
		}
		if atLeastOneFailed {
			down(db)
			t.Fatalf("DB init failed.")
		}
	}

	up(db)
	return db, func() {
		down(db)
		db.Close()
	}
}
