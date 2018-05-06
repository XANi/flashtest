package datastore

// must add to 16
var dataShards = 1
var parityShards = 16 - dataShards

//func Generate(blockid int, blocksize int, metadata interface{}) ([]byte, error) {
//	enc, err := reedsolomon.New(dataShards, parityShards)
//
//}
