package datastore

import (
	"github.com/klauspost/reedsolomon"
	"fmt"
)
type encoder struct {
	encoder reedsolomon.Encoder
}

func newEncoder() *encoder {
	enc, err := reedsolomon.New(dataShards,parityShards)
	if err != nil {
		panic(fmt.Sprintf("Can't initalize encoder: %s", err))
	}
	var e encoder
	e.encoder = enc
	return &e
}
func (e *encoder) Encode(data *[][]byte) error {
	return e.encoder.Encode(*data)
}

func (e *encoder) Verify(data [][]byte)(bool, error) {
	return e.encoder.Verify(data)
}

func (e *encoder) Repair(data *[][]byte, failedShards []uint8) error {
	// mark failed shards
	for a := range(failedShards) {
		(*data)[a] = nil
	}
	return e.encoder.Reconstruct(*data)
}