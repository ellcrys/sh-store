package common

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ellcrys/util"
)

// OpenDB returns a connection to a database
func OpenDB(conStr string) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", conStr)
	return
}

// CreateRandomDB creates a random database. Returns the database name.
func CreateRandomDB(con *sql.DB) (string, error) {
	dbName := util.RandString(5)
	_, err := con.Query(fmt.Sprintf("CREATE DATABASE %s;", dbName))
	return dbName, err
}

// DropDB drops a database
func DropDB(con *sql.DB, dbName string) error {
	_, err := con.Exec(fmt.Sprintf("DROP DATABASE %s;", dbName))
	return err
}

// ClearTable truncates tables
func ClearTable(db *sql.DB, tables ...string) error {
	_, err := db.Exec("TRUNCATE " + strings.Join(tables, ","))
	if err != nil {
		return err
	}
	return nil
}
