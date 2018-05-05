package blockdev

import (
	"os"
	"fmt"
)

type Device struct {
	path string
	blocksize int
	file *os.File
}


func NewFromFile(filename string) (Blockdev, error) {
	d :=  Device{
		path:      filename,
		blocksize: 4096,
	}
	f, err := os.OpenFile(filename, os.O_RDWR | os.O_SYNC, 0644)
	d.file  = f
	return &d, err
}
func (d *Device)GetBlocksize() int {
	return d.blocksize
}

func (d *Device)Sync() {
	d.file.Sync()
}

func (d *Device)WriteAligned(pos int, data []byte) error {
	if blkOffset := pos % d.blocksize; blkOffset > 0 {
		return fmt.Errorf("Unaligned write: %d (b: %d bs:%d)", pos, blkOffset, d.blocksize)
	}
	return d.Write(pos,data)
}
func (d *Device)ReadAligned(pos int, size int) (data []byte, err error) {
	blkOffset := pos % d.blocksize
	sizeOffset := size % d.blocksize
	pos = pos - blkOffset
	size = size - sizeOffset
	if sizeOffset > 0 {
		size += d.blocksize
	}

	return d.Read(pos,size)
}
func (d *Device)Write(pos int, data []byte) error {
	n, err := d.file.WriteAt(data,int64(pos))
	if n != len(data) {
		return fmt.Errorf("only %d out of %d bytes were written: %s", n, len(data), err)
	}
	fmt.Printf("W: %d - %d\n",pos, len(data))
	return err
}
func (d *Device)Read(pos int,size int) (data []byte, err error) {
	b := make([]byte,size)
	n, err := d.file.ReadAt(b,int64(pos))
	fmt.Printf("R: %d - %d\n",pos, size)
	if n != size {
		return b, fmt.Errorf("only %d out of %d bytes were read: %s", n, len(data), err)
	}
	return b,err
}