package blockdev

type Blockdev interface {
	GetBlocksize() int
	// write data aligned to block size.
	// Will fail if position does not align to the block size or if data size is not a multiple of
	// the blocksize
	WriteAligned(pos int, data []byte) error
	// Read blocksize of data starting at pos. Will err out if start of read is unaligned.
	// If size doesn't align with block size it will return size + last block in full
	ReadAligned(pos int, size int) (data []byte, err error)
	Write(pos int, data []byte) error
	Read(pos int,size int) (data []byte, err error)
}