package datastore

import (
	"fmt"
	"unsafe"
	"bytes"
	"encoding/binary"
)

// ID is uint16 because of byte alignment
type Block struct {
	id    uint16
	size  uint16
	data  [500]byte
	crc64 [8]byte
}
func init() {
	var b Block
	sizeOfBlock := unsafe.Sizeof(b)
	if sizeOfBlock != 512 {
		// test in case we fuck up alignment in struct or some arch messes it up
		panic(fmt.Sprintf("size of block should equal 512 not %d",sizeOfBlock))
	}
}


func NewBlockset(data []byte, blocksize int) (out [][]byte, err error) {
	blocks := dataShards + parityShards
	if blocksize%blocks != 0 {
		return out, fmt.Errorf("Blocksize %d does not divide evenly into %d encoding blocks")
	}
	if blocksize%blocks != 0 {
		return out, fmt.Errorf("Blocksize %d does not divide evenly into %d encoding blocks")
	}
	maxBlockSize := blocksize / blocks
	_ = maxBlockSize
	return out,err
}
func NewBlock(id uint8, data []byte) ( out []byte, err  error) {
    if len(data) > 500 {
		return out, fmt.Errorf("block too large: %d > 246", len(data))
	}
    var b Block
    b.id=uint16(id)
    copy(b.data[:], data[:])
	b.size = uint16(len(data))
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian,b)
	out = make([]byte,512)
	n, err := buf.Read(out)
	if n != 512 {
		return out, fmt.Errorf("read only %d bytes from buffer instead of 512", n)
	}
	return out, err
}