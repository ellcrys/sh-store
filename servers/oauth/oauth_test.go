package oauth

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ellcrys/util"
	"github.com/jinzhu/gorm"
	"github.com/ellcrys/safehold/core/db/cockroach"
	"github.com/ellcrys/safehold/core/object"
	"github.com/ellcrys/safehold/core/servers/common"
	logging "github.com/op/go-logging"
	. "github.com/smartystreets/goconvey/convey"
)

var testDB *sql.DB

var dbName = "test_" + strings.ToLower(util.RandString(5))
var conStr = "postgresql://root@localhost:26257?sslmode=disable"
var conStrWithDB = "postgresql://root@localhost:26257/" + dbName + "?sslmode=disable"

func init() {
	var err error
	testDB, err = sql.Open("postgres", conStr)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %s", err))
	}
}

func createDb(t *testing.T) error {
	_, err := testDB.Query(fmt.Sprintf("CREATE DATABASE %s;", dbName))
	return err
}

func dropDB(t *testing.T) error {
	_, err := testDB.Query(fmt.Sprintf("DROP DATABASE %s;", dbName))
	return err
}

func clearTable(db *gorm.DB, tables ...string) error {
	_, err := db.CommonDB().Exec("TRUNCATE " + strings.Join(tables, ","))
	if err != nil {
		return err
	}
	return nil
}

func TestOAuth(t *testing.T) {

	if err := createDb(t); err != nil {
		t.Fatalf("failed to create test database. %s", err)
	}
	defer dropDB(t)

	cdb := &cockroach.DB{ConnectionString: conStrWithDB}
	logging.SetLevel(logging.CRITICAL, cockroach.GetLogger().Module)
	if err := cdb.Connect(10, 5); err != nil {
		t.Fatalf("failed to connect to database. %s", err)
	}

	if err := cdb.CreateTables(); err != nil {
		t.Fatalf("failed to create tables. %s", err)
	}

	oauth := NewOAuth(cdb)

	Convey("TestOAuth", t, func() {
		Convey(".GetToken", func() {
			Convey("Should return error if grant type is not provided", func() {
				req, err := http.NewRequest("POST", "/token", nil)
				So(err, ShouldBeNil)
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(common.EasyHandle("POST", oauth.GetToken))
				handler.ServeHTTP(rr, req)
				So(rr.Code, ShouldEqual, 400)
				So(rr.Body.String(), ShouldEqual, `{"errors":[{"status":"400","message":"Grant type is required"}]}`)
			})

			Convey("Should return error if grant type is unknown", func() {
				req, err := http.NewRequest("POST", "/token?grant_type=unknown", nil)
				So(err, ShouldBeNil)
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(common.EasyHandle("POST", oauth.GetToken))
				handler.ServeHTTP(rr, req)
				So(rr.Code, ShouldEqual, 400)
				So(rr.Body.String(), ShouldEqual, `{"errors":[{"status":"400","message":"Invalid grant type"}]}`)
			})
		})

		Convey(".getAppToken", func() {

			identity := object.MakeDeveloperIdentityObject("owner_id", "creator_id", "e@email.com", "some_pass", true)
			err := oauth.obj.Create(identity)
			So(err, ShouldBeNil)

			Convey("Should return error if client id is not provided in query parameter", func() {
				req, err := http.NewRequest("POST", "/token?grant_type=client_credentials", nil)
				So(err, ShouldBeNil)
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(common.EasyHandle("POST", oauth.getAppToken))
				handler.ServeHTTP(rr, req)
				So(rr.Code, ShouldEqual, 400)
				So(rr.Body.String(), ShouldEqual, `{"errors":[{"status":"400","code":"invalid_parameter","message":"Client id is required"}]}`)
			})

			Convey("Should return error if client secret is not provided in query parameter", func() {
				req, err := http.NewRequest("POST", "/token?grant_type=client_credentials&client_id=abc", nil)
				So(err, ShouldBeNil)
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(common.EasyHandle("POST", oauth.getAppToken))
				handler.ServeHTTP(rr, req)
				So(rr.Code, ShouldEqual, 400)
				So(rr.Body.String(), ShouldEqual, `{"errors":[{"status":"400","code":"invalid_parameter","message":"Client secret is required"}]}`)
			})

			Convey("Should return error if client id is invalid", func() {
				req, err := http.NewRequest("POST", fmt.Sprintf("/token?grant_type=client_credentials&client_id=abc&client_secret=abc"), nil)
				So(err, ShouldBeNil)
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(common.EasyHandle("POST", oauth.getAppToken))
				handler.ServeHTTP(rr, req)
				So(rr.Code, ShouldEqual, 400)
				So(rr.Body.String(), ShouldEqual, `{"errors":[{"status":"400","code":"invalid_parameter","message":"Client id or secret are invalid"}]}`)
			})

			Convey("Should return error if client secret is invalid", func() {
				req, err := http.NewRequest("POST", fmt.Sprintf("/token?grant_type=client_credentials&client_id=%s&client_secret=abc", identity.Ref1), nil)
				So(err, ShouldBeNil)
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(common.EasyHandle("POST", oauth.getAppToken))
				handler.ServeHTTP(rr, req)
				So(rr.Code, ShouldEqual, 400)
				So(rr.Body.String(), ShouldEqual, `{"errors":[{"status":"400","code":"invalid_parameter","message":"Client id or secret are invalid"}]}`)
			})

			Convey("Should successfully create an app token", func() {
				req, err := http.NewRequest("POST", fmt.Sprintf("/token?grant_type=client_credentials&client_id=%s&client_secret=%s", identity.Ref1, identity.Ref2), nil)
				So(err, ShouldBeNil)
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(common.EasyHandle("POST", oauth.getAppToken))
				handler.ServeHTTP(rr, req)
				So(rr.Code, ShouldEqual, 201)
				res, _ := util.JSONToMap(rr.Body.String())
				So(res, ShouldContainKey, `access_token`)
				So(res["access_token"], ShouldNotBeEmpty)
			})
		})
	})
}
