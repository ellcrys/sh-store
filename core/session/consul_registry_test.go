package session

import (
	"testing"

	"github.com/ellcrys/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestConsulRegistry(t *testing.T) {
	r, err := NewConsulRegistry()
	if err != nil {
		t.Fatalf("failed to create consul: %s", err)
	}

	Convey("TestConsulRegistry", t, func() {
		Convey(".Add", func() {
			Convey("Should return error if session item id is not a UUID", func() {
				item := RegItem{
					Address: "127.0.0.1",
					Port:    9000,
					SID:     "abc",
					Meta: map[string]interface{}{
						"something": "here",
					},
				}
				err := r.Add(item)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "session id must be a UUID")
			})

			Convey("Should successfully add an item", func() {
				item := RegItem{
					Address: "127.0.0.1",
					Port:    9000,
					SID:     util.UUID4(),
					Meta: map[string]interface{}{
						"something": "here",
					},
				}
				err := r.Add(item)
				So(err, ShouldBeNil)

				Convey(".Get", func() {
					Convey("Should successfully get the registry item", func() {
						regItem, err := r.Get(item.SID)
						So(err, ShouldBeNil)
						So(regItem, ShouldResemble, &item)
					})

					Convey("Should return ErrNotFound if not found", func() {
						_, err = r.Get(util.UUID4())
						So(err, ShouldNotBeNil)
						So(err, ShouldEqual, ErrNotFound)
					})
				})

				Convey(".Del", func() {
					Convey("Should return error if session id is not a UUID4", func() {
						err := r.Del("abc")
						So(err, ShouldNotBeNil)
						So(err.Error(), ShouldEqual, "session id must be a UUID")
					})

					Convey("Should return nil when session does not exist", func() {
						err := r.Del(util.UUID4())
						So(err, ShouldBeNil)
					})

					Convey("Should return nil after deleting an item", func() {
						_, err = r.Get(item.SID)
						So(err, ShouldBeNil)
						err = r.Del(item.SID)
						So(err, ShouldBeNil)
						_, err = r.Get(item.SID)
						So(err, ShouldEqual, ErrNotFound)
					})
				})
			})
		})
	})
}
