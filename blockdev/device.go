package blockdev

import (
	"log"
	"os"
	"fmt"
)
var Debug bool
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
func (d *Device)GetSize() int {
	stat, err := d.file.Stat()
	if err != nil {
		log.Panicf("error seeking file: %s", err)
	}
	// block devices won't return size via stat call, need to check it in other way
	if stat.Size() == 0 {
		end, err := d.file.Seek(0, os.SEEK_END)
		if err != nil {
			log.Panicf("error checking file size: %s", err)
		}
		return int(end)
	} else {
		return int(stat.Size())
	}
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
	if Debug {fmt.Printf("W: %d - %d\n",pos, len(data))}
	return err
}
func (d *Device)Read(pos int,size int) (data []byte, err error) {
	b := make([]byte,size)
	n, err := d.file.ReadAt(b,int64(pos))
	if Debug {fmt.Printf("R: %d - %d\n",pos, size)}
	if n != size {
		return b, fmt.Errorf("only %d out of %d bytes were read: %s", n, len(data), err)
	}
	return b,err
}