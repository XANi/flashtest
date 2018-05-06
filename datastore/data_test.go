package datastore

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)
func TestBlock(t *testing.T) {
	out, err := NewBlock(1,[]byte("test"))
	Convey("test",t,func() {
		So(err,ShouldBeNil)
		So(len(out),ShouldEqual,512)
		So(string(out),ShouldContainSubstring,"\000\001\000\004test")
	})

}