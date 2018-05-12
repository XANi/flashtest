package datastore

import "hash/crc64"

// must add to 16
var dataShards = 2
var parityShards = 16 - dataShards
var totalShards = dataShards + parityShards
// this is repeated few times in code for various reasons (mostly because struct defs dont take variables as Size arguments)
var blockDataSize = 500
var blockSize = 512
var dataBlockSize =  totalShards * blockSize
var erasureSetDataSize = (blockDataSize * dataShards) - 4
var crc64table = crc64.MakeTable(crc64.ECMA)
