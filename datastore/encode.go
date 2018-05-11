package datastore

import (
	"github.com/klauspost/reedsolomon"
	"fmt"
)

type Encoder struct {
	encoder reedsolomon.Encoder
}

func NewEncoder() *Encoder {
	enc, err := reedsolomon.New(dataShards,parityShards)
	if err != nil {
		panic(fmt.Sprintf("Can't initalize encoder: %s", err))
	}
	var e Encoder
	e.encoder = enc
	return &e
}
func (e *Encoder) Encode(data *[][]byte) error {
	return e.encoder.Encode(*data)
}

func (e *Encoder) Verify(data [][]byte)(bool, error) {
	return e.encoder.Verify(data)
}

func (e *Encoder) Repair(data *[][]byte, failedShards []uint8) error {
	// mark failed shards
	for a := range(failedShards) {
		(*data)[a] = nil
	}
	return e.encoder.Reconstruct(*data)
}