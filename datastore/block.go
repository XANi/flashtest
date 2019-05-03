package datastore

import (
	"fmt"
	"encoding/binary"
)

// returns Size of Data block
func GetBlockSize() int {
	return DataBlockSize
}

// returns the Size of Data that each block carries
func GetDataSize() int {
	return erasureSetDataSize
}

// block error container
type DecodeError struct {
	Id int
	Offset int
	Msg string


}

// encode Data into datablock. Max Data Size is GetDataSize and output Size will be GetBlockSize()
func EncodeDataBlock(data []byte) ([]byte, error) {
	esBlock, err := newErasureset(data)
	if err != nil { return nil, err }
	err = enc.Encode(&esBlock)
	if err != nil {	return nil, err }
	out := make([]byte,0, DataBlockSize)
	for id, chunk := range esBlock {
		block, err := newBlock(uint8(id),chunk)
		if err != nil { return nil, err }
		out = append(out, block...)
	}
	return out,nil
}
// Decode Data block. Returns Data and list of errors
func DecodeDataBlock(data []byte) (out []byte, errlist []DecodeError,err error) {
	if len(data) != DataBlockSize {
		return out, errlist, fmt.Errorf("only full sized (%d) blocks are accepted, got [%d]", DataBlockSize, len(data))
	}
	shards := make([]*Block, totalShards)
	erasureBlock := make([][]byte, totalShards)
	var badBlocks []uint8
	for a := 0; a < totalShards; a++ {
		block, err := loadBlock(data[a*blockSize : (a+1)*blockSize])
		if err != nil {
			blkErr := DecodeError{
				Id:     a,
				Offset: blockSize * a,
				Msg:    fmt.Sprintf("Error decoding: %s", err),
			}
			errlist = append(errlist, blkErr)
			badBlocks = append(badBlocks,uint8(a))
			erasureBlock[a] = nil
			continue
		}
		if block.Id != uint16(a) {
			blkErr := DecodeError{
				Id:     a,
				Offset: blockSize * a,
				Msg:    fmt.Sprintf("Block ID mismatch: wanted %d, got %d", a, block.Id),
			}
			errlist = append(errlist, blkErr)
			continue
		}
		shards[a] = block
		erasureBlock[a] = block.GetData()
	}
	var dataGood bool
	verified, errVerify := enc.Verify(erasureBlock)
	errMsg := ""
	if verified {
		dataGood = true
	} else {
		errRepair := enc.Repair(&erasureBlock, badBlocks)
		if errRepair == nil {
			errMsg = fmt.Sprintf("verify error[%s], repaired", errVerify)
			dataGood = true
		} else {
			errMsg = fmt.Sprintf("verify error[%s], couldn't repair: %s", errVerify,errRepair)
		}
		errlist = append(errlist,DecodeError{
			Id:     -1,
			Offset: -1,
			Msg: errMsg,
		})
	}
	if !dataGood {
		return out, errlist, fmt.Errorf("err: %s",errMsg)
	}
	var verifiedData []byte
	for idx := 0;  idx < dataShards; idx++ {
		verifiedData = append(verifiedData, erasureBlock[idx]...)
	}
	dataLength := binary.BigEndian.Uint32(verifiedData[:4])
	if len(verifiedData) < int(dataLength + 4) {
		return out, errlist, fmt.Errorf("Data length mismatch: read %d bytes but header indicates %d",len(verifiedData), dataLength)
	}
	return verifiedData[4:dataLength+4], errlist,nil
}

