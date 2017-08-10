package common

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ellcrys/util"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHTTPCommon(t *testing.T) {
	Convey("TestHTTPCommon", t, func() {
		Convey(".EasyHandle", func() {

			Convey("Should return error if method of request does not match the request method associated with the handler", func() {
				req, err := http.NewRequest("GET", "/some_endpoint", nil)
				So(err, ShouldBeNil)
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(EasyHandle("POST", func(w http.ResponseWriter, r *http.Request) (interface{}, int) {
					return nil, 0
				}))
				handler.ServeHTTP(rr, req)
				So(rr.Code, ShouldEqual, 404)
				m, _ := util.JSONToMap(rr.Body.String())
				So(m, ShouldResemble, map[string]interface{}{
					"code":    "unsupported_method",
					"message": "method not supported for this endpoint",
					"status":  "404",
				})
			})

			Convey("Should return expected Error json body ", func() {
				req, err := http.NewRequest("GET", "/some_endpoint", nil)
				So(err, ShouldBeNil)
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(EasyHandle("GET", func(w http.ResponseWriter, r *http.Request) (interface{}, int) {
					return Error(200, "some_code", "field_name", "some text"), 0
				}))
				handler.ServeHTTP(rr, req)
				So(rr.Code, ShouldEqual, 200)
				So(rr.HeaderMap["Content-Type"][0], ShouldEqual, "application/json")
				m, _ := util.JSONToMap(rr.Body.String())
				So(m, ShouldContainKey, "errors")
				errs := m["errors"].([]interface{})
				So(errs, ShouldHaveLength, 1)
				So(errs[0], ShouldResemble, map[string]interface{}{
					"status": "200",
					"code":   "some_code",
					"source": "field_name",
					"detail": "some text",
				})
			})

			Convey("Should return string if response is an Error whose message is not a json string", func() {
				req, err := http.NewRequest("GET", "/some_endpoint", nil)
				So(err, ShouldBeNil)
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(EasyHandle("GET", func(w http.ResponseWriter, r *http.Request) (interface{}, int) {
					return fmt.Errorf("some error text"), 0
				}))
				handler.ServeHTTP(rr, req)
				So(rr.Code, ShouldEqual, 0)
				So(rr.Body.String(), ShouldContainSubstring, `some error text`)
			})

			Convey("Should return string if response is a text", func() {
				req, err := http.NewRequest("GET", "/some_endpoint", nil)
				So(err, ShouldBeNil)
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(EasyHandle("GET", func(w http.ResponseWriter, r *http.Request) (interface{}, int) {
					return "some error text", 0
				}))
				handler.ServeHTTP(rr, req)
				So(rr.Code, ShouldEqual, 200)
				So(rr.HeaderMap["Content-Type"][0], ShouldEqual, "text/plain; charset=utf-8")
				So(rr.Body.String(), ShouldEqual, `some error text`)
			})

			Convey("Should return json body if response is a map", func() {
				req, err := http.NewRequest("GET", "/some_endpoint", nil)
				So(err, ShouldBeNil)
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(EasyHandle("GET", func(w http.ResponseWriter, r *http.Request) (interface{}, int) {
					return map[string]string{
						"name": "ben",
					}, 200
				}))
				handler.ServeHTTP(rr, req)
				So(rr.Code, ShouldEqual, 200)
				So(rr.HeaderMap["Content-Type"][0], ShouldEqual, "application/json")
				So(rr.Body.String(), ShouldContainSubstring, `{"name":"ben"}`)
			})

			Convey("Should return json body if response is a struct", func() {
				req, err := http.NewRequest("GET", "/some_endpoint", nil)
				So(err, ShouldBeNil)
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(EasyHandle("GET", func(w http.ResponseWriter, r *http.Request) (interface{}, int) {
					return &struct {
						Name string
					}{Name: "ben"}, 200
				}))
				handler.ServeHTTP(rr, req)
				So(rr.Code, ShouldEqual, 200)
				So(rr.HeaderMap["Content-Type"][0], ShouldEqual, "application/json")
				So(rr.Body.String(), ShouldContainSubstring, `{"data":{"type":""}`)
			})
		})
	})
}
