package common

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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
				So(rr.Body.String(), ShouldContainSubstring, `{"status":"404","code":"unsupported_method","message":"method not supported for this endpoint"}`)
			})

			Convey("Should return json body for APIError", func() {
				req, err := http.NewRequest("GET", "/some_endpoint", nil)
				So(err, ShouldBeNil)
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(EasyHandle("GET", func(w http.ResponseWriter, r *http.Request) (interface{}, int) {
					return NewSingleAPIErr(200, "some_code", "field_name", "some text", nil), 0
				}))
				handler.ServeHTTP(rr, req)
				So(rr.Code, ShouldEqual, 200)
				So(rr.HeaderMap["Content-Type"][0], ShouldEqual, "application/json")
				So(rr.Body.String(), ShouldContainSubstring, `{"errors":[{"status":"200","code":"some_code","field":"field_name","message":"some text"}]}`)
			})

			Convey("Should return string if response is an error whose message is not a json string", func() {
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
					return struct {
						Name string
					}{Name: "ben"}, 200
				}))
				handler.ServeHTTP(rr, req)
				So(rr.Code, ShouldEqual, 200)
				So(rr.HeaderMap["Content-Type"][0], ShouldEqual, "application/json")
				So(rr.Body.String(), ShouldContainSubstring, `{"Name":"ben"}`)
			})
		})
	})
}
