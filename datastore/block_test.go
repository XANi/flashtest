package datastore

import (
	"testing"
	"strings"
	. "github.com/smartystreets/goconvey/convey"
	"fmt"
)

func TestEncodeDataBlock(t *testing.T) {
	testdata := []byte(strings.Repeat("x",500))
	encoded, err := EncodeDataBlock(testdata)
	Convey("New encoded block",t,func() {
		So(err,ShouldBeNil)
		So(len(encoded),ShouldEqual, GetBlockSize())
	})


	decoded, errlist, err := DecodeDataBlock(encoded)
	if len(errlist) > 0 {
		for i, e := range errlist {
			fmt.Printf("e: %d [%+v]\n", i, e)
		}
	}
	Convey("Decode encoded block",t,func() {
		So(err,ShouldBeNil)
		So(decoded,ShouldResemble, testdata)
		So(len(errlist),ShouldEqual,0)
	})

	encoded[513] = 'z'
	decoded, errlist, err = DecodeDataBlock(encoded)
	if len(errlist) > 0 {
		for i, e := range errlist {
			fmt.Printf("e: %d [%+v]\n", i, e)
		}
	}
	Convey("Decode encoded block with error",t,func() {
		So(err,ShouldBeNil)
		So(len(errlist),ShouldEqual,2)
		So(string(decoded),ShouldEqual, string(testdata))
	})


}
// from fib_test.go
func BenchmarkBlock_UpdateChecksum(b *testing.B) {
	blockRaw, _ := newBlock(5, []byte(strings.Repeat("x",500)))
	block, _ := loadBlock(blockRaw)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		block.UpdateChecksum()
	}
}
func BenchmarkBlock_VerifyChecksum(b *testing.B) {
	blockRaw, _ := newBlock(5, []byte(strings.Repeat("x",500)))
	block, _ := loadBlock(blockRaw)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		block.VerifyChecksum()
	}
}