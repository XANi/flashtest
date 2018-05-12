package datastore

import (
	"fmt"
	"bytes"
	"encoding/binary"
	"hash/crc64"
)

// ID is uint16 because of byte alignment
type Block struct {
	Id    uint16
	Size  uint16
	Data  [500]byte
	Crc64 uint64
}

// Split input Data into shards to be used with encoder
func newErasureset(data []byte) (out [][]byte, err error) {
	sizeRaw := make([]byte,4)
	binary.BigEndian.PutUint32(sizeRaw,uint32(len(data)))
	data = append(sizeRaw, data...)
	if len(data) >  (blockDataSize * dataShards){
		return out, fmt.Errorf("GetData + header is bigger[%d] than total Erasureset capacity, should get at most %d bytes", len(data), erasureSetDataSize)
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


// create new Data block (block that is supposed to be saved in target file/device)
func newBlock(id uint8, data []byte) ( out []byte, err  error) {
    if len(data) > blockDataSize {
		return out, fmt.Errorf("block too large: %d > 500", len(data))
	}
    var b Block
    b.Id =uint16(id)
    copy(b.Data[:], data[:])
	b.Size = uint16(len(data))
	b.UpdateChecksum()
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian,b)
	out = make([]byte,512)
	n, err := buf.Read(out)
	if n != blockSize {
		return out, fmt.Errorf("read only %d bytes from buffer instead of 512", n)
	}
	return out, err
}
func loadBlock(data []byte) (*Block, error) {
	var b Block
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.BigEndian, &b)
	if err != nil {
		return nil, err
	}
	if !b.VerifyChecksum() {
		return nil, fmt.Errorf("Block checksum mismatch")
	}
	return &b, err
}
func(b *Block) GetData() []byte {
	return b.Data[:b.Size]
}
func(b *Block) encode()[]byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian,b)
	out := make([]byte,512)
	buf.Read(out)
	return out
}
func(b *Block) UpdateChecksum() {
	data := (b.encode())[:512-8]
	b.Crc64 = crc64.Checksum(data, crc64table)
}
func(b *Block) VerifyChecksum() bool{
	data := (b.encode())[:512-8]
	return b.Crc64 == crc64.Checksum(data,crc64table)
}
func datasize() int {
	return blockDataSize
}
func blocksize() int {
	return blockSize
}