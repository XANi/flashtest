package datastore

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"encoding/binary"
)
func TestBlock(t *testing.T) {
	_, err := newBlock(1,[]byte(strings.Repeat("x",1000)))
	Convey("Fail on too big input",t,func() {
		So(err, ShouldNotBeNil)
	})
	out, err := newBlock(1,[]byte(strings.Repeat("t", datasize())))
	Convey("test",t,func() {
		So(err,ShouldBeNil)
		So(len(out),ShouldEqual, blocksize())
		So(string(out),ShouldStartWith,"\000\001")
		So(string(out),ShouldContainSubstring,"tttt")
	})

}

func TestErasureset(t *testing.T) {
	testblkData := []byte(strings.Repeat("x",blockDataSize-4))
	es, err := newErasureset(testblkData)
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
	es, err = newErasureset(testblkData2)
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
	es, err = newErasureset(testblkData3)
	expected3 := append(testblkLen3, []byte("xxxx")...)
	Convey("Create Erasureset with partial block", t, func() {
		So(err,ShouldBeNil)
		So(string(es[0]),ShouldContainSubstring,string(expected3))
		So(string(es[1]),ShouldStartWith,string("yyyy"))
	})

	_, err = newErasureset([]byte(strings.Repeat("x", erasureSetDataSize + 1)))
	Convey("Create Erasureset too much Data", t, func() {
		So(err, ShouldNotBeNil)
	})
}

func BenchmarkEncoder_Encode(b *testing.B) {
	testblkData := []byte(strings.Repeat("x",erasureSetDataSize))
	es, _ := newErasureset(testblkData)
	bb := make([][]byte,len(es))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		copy(bb, es)
		enc.Encode(&bb)
	}
}
func BenchmarkEncoder_Repair(b *testing.B) {
	testblkData := []byte(strings.Repeat("x",erasureSetDataSize))
	es, _ := newErasureset(testblkData)
	enc.Encode(&es)
	bb := make([][]byte,len(es))
	es[0] = nil
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		copy(bb, es)
		enc.Repair(&bb,[]uint8{})
	}
}
func BenchmarkEncoder_Verify(b *testing.B) {
	testblkData := []byte(strings.Repeat("x",erasureSetDataSize))
	es, _ := newErasureset(testblkData)
	enc.Encode(&es)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		enc.Verify(es)
	}
}
