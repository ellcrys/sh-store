package oauth

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/ellcrys/elldb/core/servers/common"
	"github.com/ellcrys/elldb/core/servers/db"
	"github.com/jinzhu/gorm"
	. "github.com/smartystreets/goconvey/convey"
)

var err error
var testDB *sql.DB
var dbName string
var conStr = "postgresql://root@localhost:26257?sslmode=disable"
var conStrWithDB string

func init() {
	SigningSecret = "secret"
	testDB, err = common.OpenDB(conStr)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %s", err))
	}
	dbName, err = common.CreateRandomDB(testDB)
	if err != nil {
		panic(fmt.Errorf("failed to create test database: %s", err))
	}
	conStrWithDB = "postgresql://root@localhost:26257/" + dbName + "?sslmode=disable"
}

func TestOAuth(t *testing.T) {

	dbCon, err := gorm.Open("postgres", conStrWithDB)
	if err != nil {
		t.Fatal("failed to connect to test database")
	}

	err = dbCon.CreateTable(&db.Bucket{}, &db.Object{}, &db.Account{}, &db.Mapping{}).Error
	if err != nil {
		t.Fatalf("failed to create database tables. %s", err)
	}

	defer func() {
		common.DropDB(testDB, dbName)
	}()

	o := NewOAuth(dbCon)

	Convey("TestOAuth", t, func() {

	})
}
