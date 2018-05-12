package datastore

import (
	"unsafe"
	"fmt"
)

var enc *encoder

func init() {
	var b Block
	sizeOfBlock := unsafe.Sizeof(b)
	if int(sizeOfBlock) != blockSize {
		// test in case we fuck up alignment in struct or some arch messes it up
		panic(fmt.Sprintf("Size of block should equal 512 not %d",sizeOfBlock))
	}
	enc = newEncoder()
}