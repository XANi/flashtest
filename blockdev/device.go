package blockdev


type Device struct {
	path string
	blocksize int
}


func NewFromFile(filename string) (Blockdev, error) {
	d :=  Device{
		path:      filename,
		blocksize: 4096,
	}
	return &d, nil
}
func (d *Device)GetBlocksize() int {
	return d.blocksize
}

func (d *Device)WriteAligned(pos int, data []byte) error {
	return nil
}
func (d *Device)ReadAligned(pos int, size int) (data []byte, err error) {
	return []byte{}, nil
}
func (d *Device)Write(pos int, data []byte) error {
	return nil
}
func (d *Device)Read(pos int,size int) (data []byte, err error) {
	return []byte{},nil
}