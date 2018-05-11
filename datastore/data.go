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
	if int(sizeOfBlock) != blockSize {
		// test in case we fuck up alignment in struct or some arch messes it up
		panic(fmt.Sprintf("size of block should equal 512 not %d",sizeOfBlock))
	}
}

// Split input data into shards to be used with encoder
func NewErasureset(data []byte) (out [][]byte, err error) {
	sizeRaw := make([]byte,4)
	binary.BigEndian.PutUint32(sizeRaw,uint32(len(data)))
	data = append(sizeRaw, data...)
	if len(data) >  (blockDataSize * dataShards){
		return out, fmt.Errorf("Data + header is bigger[%d] than total Erasureset capacity, should get at most %d bytes", len(data), erasureSetDataSize)
	}
	out = make([][]byte, totalShards)
	var chunk []byte
	blkid := 0
	for len(data) >= blockDataSize {
		chunk, data = data[:blockDataSize], data[blockDataSize:]
		out[blkid] = chunk
		blkid++
	}
	if len(data) > 0 {
		out[blkid] = make([]byte,blockDataSize)
		copy(out[blkid], data)
		blkid++
	}
	for blkid < totalShards {
		out[blkid] = make([]byte,blockDataSize)
		blkid++
	}
	return out,err
}


// create new data block (block that is supposed to be saved in target file/device)
func NewBlock(id uint8, data []byte) ( out []byte, err  error) {
    if len(data) > blockDataSize {
		return out, fmt.Errorf("block too large: %d > 500", len(data))
	}
    var b Block
    b.id=uint16(id)
    copy(b.data[:], data[:])
	b.size = uint16(len(data))
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian,b)
	out = make([]byte,512)
	n, err := buf.Read(out)
	if n != blockSize {
		return out, fmt.Errorf("read only %d bytes from buffer instead of 512", n)
	}
	return out, err
}
func Datasize() int {
	return blockDataSize
}
func Blocksize() int {
	return blockSize
}