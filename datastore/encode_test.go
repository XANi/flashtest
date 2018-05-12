package datastore

import (
		. "github.com/smartystreets/goconvey/convey"
	"testing"
	"fmt"
)

func TestEncoder(t *testing.T) {
	testblkData2 := append([]byte("xxxxxx"))
	es, _ := newErasureset(testblkData2)
	enc := newEncoder()
	err := enc.Encode(&es)
	verify, err2 := enc.Verify(es)
	Convey("encode", t, func() {
		So(err,ShouldBeNil)
		So(err2,ShouldBeNil)
		So(verify,ShouldBeTrue)
	})
	// drop Data
	oldBlkVal := es[0]
	es[0]=nil
	verify, err2 = enc.Verify(es)
	Convey("verify err on missing", t, func() {
		So(err2,ShouldNotBeNil)
		So(verify,ShouldBeFalse)
	})
	// corrupt shard
	es[0]=make([]byte,len(oldBlkVal))
	verify, err2 = enc.Verify(es)
	Convey("verify fail on corrupted", t, func() {
		So(err2,ShouldBeNil)
		So(verify,ShouldBeFalse)
	})
	// regen shard
	err = enc.Repair(&es,[]uint8{0})
	verify, err2 = enc.Verify(es)
	Convey("repair shard", t, func() {
		So(err, ShouldBeNil)
		//ShouldEqual is picky...
		So(fmt.Sprintf("%+v", es[0]),ShouldEqual,fmt.Sprintf("%+v", oldBlkVal))
		So(err2, ShouldBeNil)
		So(verify,ShouldBeTrue)
	})

}