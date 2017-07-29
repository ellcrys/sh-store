package oauth

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ellcrys/elldb/servers/common"
	"github.com/ellcrys/elldb/servers/db"
	"github.com/ellcrys/util"
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
		Convey(".GetToken", func() {
			Convey("Should return error if grant type is not provided", func() {
				req, err := http.NewRequest("POST", "/token", nil)
				So(err, ShouldBeNil)
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(common.EasyHandle("POST", o.GetToken))
				handler.ServeHTTP(rr, req)
				So(rr.Code, ShouldEqual, 400)
				m, _ := util.JSONToMap(rr.Body.String())
				errs := m["errors"].([]interface{})
				So(errs, ShouldHaveLength, 1)
				So(errs[0], ShouldResemble, map[string]interface{}{
					"status":  "400",
					"message": "Grant type is required",
				})
			})

			Convey("Should return error if grant type is unknown", func() {
				req, err := http.NewRequest("POST", "/token?grant_type=unknown", nil)
				So(err, ShouldBeNil)
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(common.EasyHandle("POST", o.GetToken))
				handler.ServeHTTP(rr, req)
				So(rr.Code, ShouldEqual, 400)
				m, _ := util.JSONToMap(rr.Body.String())
				errs := m["errors"].([]interface{})
				So(errs, ShouldHaveLength, 1)
				So(errs[0], ShouldResemble, map[string]interface{}{
					"status":  "400",
					"message": "Invalid grant type",
				})
			})
		})

		Convey(".getAppToken", func() {

			account := db.NewAccount()
			account.ClientID = util.RandString(10)
			account.ClientSecret = util.RandString(10)
			err := dbCon.Create(account).Error
			So(err, ShouldBeNil)

			Convey("Should return error if client id is not provided in query parameter", func() {
				req, err := http.NewRequest("POST", "/token?grant_type=client_credentials", nil)
				So(err, ShouldBeNil)
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(common.EasyHandle("POST", o.getAppToken))
				handler.ServeHTTP(rr, req)
				So(rr.Code, ShouldEqual, 400)
				m, _ := util.JSONToMap(rr.Body.String())
				errs := m["errors"].([]interface{})
				So(errs[0], ShouldResemble, map[string]interface{}{
					"status":  "400",
					"code":    "invalid_parameter",
					"message": "Client id is required",
				})
			})

			Convey("Should return error if client secret is not provided in query parameter", func() {
				req, err := http.NewRequest("POST", "/token?grant_type=client_credentials&client_id=abc", nil)
				So(err, ShouldBeNil)
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(common.EasyHandle("POST", o.getAppToken))
				handler.ServeHTTP(rr, req)
				So(rr.Code, ShouldEqual, 400)
				m, _ := util.JSONToMap(rr.Body.String())
				errs := m["errors"].([]interface{})
				So(errs, ShouldHaveLength, 1)
				So(errs[0], ShouldResemble, map[string]interface{}{
					"status":  "400",
					"code":    "invalid_parameter",
					"message": "Client secret is required",
				})
			})

			Convey("Should return error if client id is invalid", func() {
				req, err := http.NewRequest("POST", fmt.Sprintf("/token?grant_type=client_credentials&client_id=abc&client_secret=abc"), nil)
				So(err, ShouldBeNil)
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(common.EasyHandle("POST", o.getAppToken))
				handler.ServeHTTP(rr, req)
				So(rr.Code, ShouldEqual, 400)
				m, _ := util.JSONToMap(rr.Body.String())
				errs := m["errors"].([]interface{})
				So(errs, ShouldHaveLength, 1)
				So(errs[0], ShouldResemble, map[string]interface{}{
					"status":  "400",
					"code":    "invalid_parameter",
					"message": "Client id or secret are invalid",
				})
			})

			Convey("Should return error if client secret is invalid", func() {
				req, err := http.NewRequest("POST", fmt.Sprintf("/token?grant_type=client_credentials&client_id=%s&client_secret=abc", account.ClientID), nil)
				So(err, ShouldBeNil)
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(common.EasyHandle("POST", o.getAppToken))
				handler.ServeHTTP(rr, req)
				So(rr.Code, ShouldEqual, 400)
				m, _ := util.JSONToMap(rr.Body.String())
				errs := m["errors"].([]interface{})
				So(errs, ShouldHaveLength, 1)
				So(errs[0], ShouldResemble, map[string]interface{}{
					"status":  "400",
					"code":    "invalid_parameter",
					"message": "Client id or secret are invalid",
				})
			})

			Convey("Should successfully create an app token", func() {
				req, err := http.NewRequest("POST", fmt.Sprintf("/token?grant_type=client_credentials&client_id=%s&client_secret=%s", account.ClientID, account.ClientSecret), nil)
				So(err, ShouldBeNil)
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(common.EasyHandle("POST", o.getAppToken))
				handler.ServeHTTP(rr, req)
				So(rr.Code, ShouldEqual, 201)
				res, _ := util.JSONToMap(rr.Body.String())
				So(res, ShouldContainKey, `access_token`)
				So(res["access_token"], ShouldNotBeEmpty)
			})
		})
	})
}
