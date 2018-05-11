package datastore

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"encoding/binary"
)
func TestBlock(t *testing.T) {
	_, err := NewBlock(1,[]byte(strings.Repeat("x",1000)))
	Convey("Fail on too big input",t,func() {
		So(err, ShouldNotBeNil)
	})
	out, err := NewBlock(1,[]byte(strings.Repeat("t",Datasize())))
	Convey("test",t,func() {
		So(err,ShouldBeNil)
		So(len(out),ShouldEqual,Blocksize())
		So(string(out),ShouldStartWith,"\000\001")
		So(string(out),ShouldContainSubstring,"tttt")
	})

}

func TestErasureset(t *testing.T) {
	testblkData := []byte(strings.Repeat("x",blockDataSize-4))
	es, err := NewErasureset(testblkData)
	testblkLen := make([]byte,4)
	binary.BigEndian.PutUint32(testblkLen,uint32(len(testblkData)))
	expected := append(testblkLen, []byte("xxxx")...)

	Convey("Create Erasureset with 1 block", t, func() {
		So(err,ShouldBeNil)
		So(string(es[0]),ShouldContainSubstring,string(expected))
	})

	testblkData2 := append([]byte(strings.Repeat("x",blockDataSize  - 4)),
		[]byte(strings.Repeat("y",blockDataSize))...
	)
	testblkLen2 := make([]byte,4)
	binary.BigEndian.PutUint32(testblkLen2,uint32(len(testblkData2)))
	es, err = NewErasureset(testblkData2)
	expected2 := append(testblkLen2, []byte("xxxx")...)
	Convey("Create Erasureset with 2 blocks", t, func() {
		So(err,ShouldBeNil)
		So(string(es[0]),ShouldContainSubstring,string(expected2))
		So(string(es[1]),ShouldStartWith,string("yyyy"))
	})

	testblkData3 := append([]byte(strings.Repeat("x",blockDataSize  - 4)),
		[]byte(strings.Repeat("y",8))...
	)
	testblkLen3 := make([]byte,4)
	binary.BigEndian.PutUint32(testblkLen3,uint32(len(testblkData3)))
	es, err = NewErasureset(testblkData3)
	expected3 := append(testblkLen3, []byte("xxxx")...)
	Convey("Create Erasureset with partial block", t, func() {
		So(err,ShouldBeNil)
		So(string(es[0]),ShouldContainSubstring,string(expected3))
		So(string(es[1]),ShouldStartWith,string("yyyy"))
	})

	_, err = NewErasureset([]byte(strings.Repeat("x", erasureSetDataSize + 1)))
	Convey("Create Erasureset too much data", t, func() {
		So(err, ShouldNotBeNil)
	})
}