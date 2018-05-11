package datastore

// must add to 16
var dataShards = 2
var parityShards = 16 - dataShards
var totalShards = dataShards + parityShards
// this is repeated few times in code for various reasons (mostly because struct defs dont take variables as size arguments)
var blockDataSize = 500
var blockSize = 512
var erasureSetDataSize = (blockDataSize * dataShards) - 4
//func Generate(blockid int, blocksize int, metadata interface{}) ([]byte, error) {
//	enc, err := reedsolomon.New(dataShards, parityShards)
//
//}
