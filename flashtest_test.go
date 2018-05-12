package main

import (
	"testing"
	"io/ioutil"
	"strings"
	"github.com/XANi/flashtest/blockdev"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/XANi/flashtest/datastore"
)

var testStrings []string

func TestDataIO(t *testing.T) {
	// print out debug data during test
	blockdev.Debug = true
	testfileName := `./t-data/flashblock.test`
	testData := []byte("testcat")
	errIo := ioutil.WriteFile(testfileName, []byte(strings.Repeat("y",1024*1024)),0600)
	f, err := blockdev.NewFromFile(testfileName)
	Convey("TestOpen", t, func() {
		So(errIo, ShouldBeNil)
		So(err,   ShouldBeNil)
	})
	out, errDatastore := datastore.EncodeDataBlock(testData)
	errWrite := f.WriteAligned(3*datastore.GetBlockSize(),out)
	Convey("TestWrite", t, func() {
		So(errDatastore, ShouldBeNil)
		So(errWrite,   ShouldBeNil)
	})
	f.Sync()
	readDataAligned,err := f.ReadAligned(3*datastore.GetBlockSize(),datastore.GetBlockSize())
	out, errlist, err := datastore.DecodeDataBlock(readDataAligned)
	Convey("TestRead",t,func() {
		So(err,ShouldBeNil)
		So(len(errlist),ShouldEqual,0)
		So(out,ShouldResemble,testData)
	})

}
