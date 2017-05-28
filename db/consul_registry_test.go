package db

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConsulRegistry(t *testing.T) {
	r, err := NewConsulRegistry()
	if err != nil {
		t.Fatalf("failed to create consul: %s", err)
	}

	Convey("TestConsulRegistry", t, func() {
		Convey(".Add", func() {
			Convey("Should successfully add an item", func() {
				item := RegItem{
					Address: "127.0.0.1",
					Port:    9000,
					SID:     "abc",
					Meta: map[string]interface{}{
						"something": "here",
					},
				}
				err := r.Add(item)
				So(err, ShouldBeNil)

				Convey(".Get", func() {
					Convey("Should successfully get the registry item", func() {
						regItem, err := r.Get("abc")
						So(err, ShouldBeNil)
						So(regItem, ShouldResemble, &item)
					})

					Convey("Should return ErrNotFound if not found", func() {
						_, err = r.Get("abcd")
						So(err, ShouldNotBeNil)
						So(err, ShouldEqual, ErrNotFound)
					})
				})

				Convey(".Del", func() {
					Convey("Should return nil when session does not exist", func() {
						err := r.Del("xyz")
						So(err, ShouldBeNil)
					})

					Convey("Should return nil after deleting an item", func() {
						_, err = r.Get("abc")
						So(err, ShouldBeNil)
						err = r.Del("abc")
						So(err, ShouldBeNil)
						_, err = r.Get("abc")
						So(err, ShouldEqual, ErrNotFound)
					})
				})
			})
		})
	})
}
