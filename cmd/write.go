package cmd

import (
	"github.com/XANi/flashtest/blockdev"
	"github.com/XANi/flashtest/datastore"
	"github.com/op/go-logging"
	"fmt"
)
var log = logging.MustGetLogger("main")
func WriteFile(filename string, filesize int) {
	dev, err := blockdev.NewFromFile(filename)
	if err != nil {
		log.Criticalf("error opening [%s]:%s", filename, err)
		return
	}
	if filesize < 1 {
		filesize = dev.GetSize()
	}
	dataBlocks := int(filesize/datastore.DataBlockSize)
	checkedSize := dataBlocks * datastore.DataBlockSize
	if dataBlocks < 1 {
		log.Errorf("file too small, must be at least %d bytes",datastore.DataBlockSize)
		return
	}
	if filesize > checkedSize {
		log.Warningf("File is not aligned to data block size[%d]. Will not touch last %d bytes", datastore.DataBlockSize, filesize-checkedSize)
	}
	for i := 0; i < dataBlocks; i++ {
		offset := i * datastore.DataBlockSize
		log.Infof("Writing block %d at offset %d",i, offset)
		data, err := datastore.EncodeDataBlock([]byte(fmt.Sprintf("Block %d",i)))
		_ = err // handle err
		err = dev.Write(offset,data)
		_ = err // handle write errors
	}
}
type verifyQ struct {
	block []byte
	offset int
}
type reportQ struct {
	pos int
	size int
	ok bool
}
func VerifyFile(filename string, filesize int) {
	dev, err := blockdev.NewFromFile(filename)
	if err != nil {
		log.Criticalf("error opening [%s]:%s", filename, err)
		return
	}
	if filesize < 1 {
		filesize = dev.GetSize()
	}
	dataBlocks := int(filesize/datastore.DataBlockSize)
	checkedSize := dataBlocks * datastore.DataBlockSize
	if dataBlocks < 1 {
		log.Errorf("file too small, must be at least %d bytes",datastore.DataBlockSize)
		return
	}
	if filesize > checkedSize {
		log.Warningf("File is not aligned to data block size[%d]. Will not touch last %d bytes", datastore.DataBlockSize, filesize-checkedSize)
	}
	vq := make(chan verifyQ,5)
	end := make(chan bool,1)
	go func() {
		  for qe := range vq {
		  	ok,out,errlist,err:=verifyBlock(qe.block,qe.offset)
			  if ok {
				  log.Infof("decoded block %d: %s", qe.offset, string(out))
			  }
			  if len(errlist) > 0 {
				  log.Warning("  found errors")
				  for _, e := range errlist {
					  log.Warningf("  [%d] error: %+v", qe.offset, e)
				  }
			  }
			  if err != nil {
				  log.Errorf("error decoding block %d: %s",qe.offset,err)
			  }
		  }
		  end <- true
	} ()
	for i := 0; i < dataBlocks; i++ {
		offset := i * datastore.DataBlockSize
		data,_ := dev.Read(offset,datastore.DataBlockSize)
		vq <- verifyQ{block:data,offset:offset}
	}
	close(vq)
	_ = <- end
}

func verifyBlock(block []byte,offset int) (ok bool, data []byte,errl []datastore.DecodeError, err error) {
	out,errlist,err := datastore.DecodeDataBlock(block)
	if len(errlist) < 1 && err == nil {
		return true, out, errlist, nil
	}
	if err != nil {
		return false, []byte{}, errlist, err
	}
	return true,out,errlist,nil

}

