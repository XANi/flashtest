package blockdev
import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"strings"
	"io/ioutil"
)
func TestFile(t *testing.T) {
	testfileName := `../t-data/block.test`
	ioutil.WriteFile(testfileName, []byte(strings.Repeat("y",1024*1024)),0600)
	f, err := NewFromFile(testfileName)
	Convey("TestOpen", t, func() {
		So(err, ShouldEqual, nil)
	})
	readData,err := f.Read(1,1)
	Convey("TestRead",t,func() {
		So(err, ShouldEqual, nil)
		So(string(readData), ShouldEqual, "y")
	})
	readDataAligned,err := f.ReadAligned(1,1)
	Convey("TestReadAligned",t,func() {
		So(err, ShouldEqual, nil)
		So(string(readDataAligned), ShouldEqual, strings.Repeat("y",f.GetBlocksize()))

	})

	err1 := f.Write(1,[]byte("z"))
	readData,err2 := f.Read(1,1)
	Convey("TestWriteAligned", t, func() {
		So(err1, ShouldEqual, nil)
		So(err2, ShouldEqual, nil)
		So(string(readData),ShouldEqual,"z")
	})
	testdata := []byte(strings.Repeat("x", f.GetBlocksize()))
	err1 = f.WriteAligned(10 * f.GetBlocksize(),testdata)
	readDataAligned,err2 = f.ReadAligned(1,1)
	Convey("TestWriteAligned", t, func() {
		So(err1, ShouldEqual, nil)
		So(err2, ShouldEqual, nil)
		So(string(readDataAligned),ShouldEqual,string(testdata))
	})

}
