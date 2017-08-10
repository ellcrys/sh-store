package common

import (
	"database/sql"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ellcrys/util"
)

// OpenDB returns a connection to a database
func OpenDB(conStr string) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", conStr)
	return
}

// CreateRandomDB creates a random database. Returns the database name.
func CreateRandomDB() (string, error) {
	dbName := "test_" + util.RandString(5)
	return dbName, exec.Command("createdb", dbName).Run()
}

// DropDB drops a database
func DropDB(dbName string) error {
	fmt.Println("Dropping database: ", dbName)
	out, err := exec.Command("dropdb", dbName).CombinedOutput()
	fmt.Println(string(out))
	return err
	// return exec.Command("dropdb", dbName).Run()
}

// ClearTable truncates tables
func ClearTable(db *sql.DB, tables ...string) error {
	_, err := db.Exec("TRUNCATE " + strings.Join(tables, ","))
	if err != nil {
		return err
	}
	return nil
}
