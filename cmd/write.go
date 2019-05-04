package cmd

import (
	"github.com/XANi/flashtest/blockdev"
	"github.com/XANi/flashtest/datastore"
	"github.com/op/go-logging"
	"fmt"
	"runtime"
	"sync"
	"time"
)
var log = logging.MustGetLogger("main")

type writeQ struct {
	block []byte
	offset int
}
type encodeQ struct {
	block []byte
	offset int
}
func writeQWorker(in chan writeQ, dev blockdev.Blockdev) {
	for qe := range in{
		err := dev.Write(qe.offset,qe.block)
		if err != nil {log.Errorf("write error at %d:%s", qe.offset,err)}
	}
	log.Infof("Write done")
}

func encodeQWorker(in chan encodeQ, out chan writeQ) {
	for e := range in {
		 data,err := datastore.EncodeDataBlock(e.block)
		 if err != nil {
			 log.Warningf("error encoding %s, data:%s", err, data)
		 } else {
			 out <- writeQ{block: data, offset: e.offset}
		 }
	}
}

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
	log.Infof("blocks to write: %d",dataBlocks)
	lastOutput:=time.Now()
	var wgWrite sync.WaitGroup
	var wgEncode sync.WaitGroup
	wq := make(chan writeQ,runtime.NumCPU()*4)
	eq := make(chan encodeQ,runtime.NumCPU()*4)
	for i:=0; i < runtime.NumCPU();i++ {
		go func() {
			wgEncode.Add(1)
			encodeQWorker(eq, wq)
			wgEncode.Done()
		}()
	}
	go func() {
		wgWrite.Add(1)
		writeQWorker(wq,dev)
		wgWrite.Done()

	}()
	for i := 0; i < dataBlocks; i++ {
		offset := i * datastore.DataBlockSize
		if lastOutput.Add(time.Second * 10).Before(time.Now()) {
			lastOutput = time.Now()
			log.Infof("Writing block %d/%d at offset %d",i,dataBlocks, offset)
		}
		eq <- encodeQ{block: []byte(fmt.Sprintf("Block %d",i)),offset:offset}
		//data, err := datastore.EncodeDataBlock([]byte(fmt.Sprintf("Block %d",i)))


		_ = err // handle err
		//wq <- writeQ{block:data,offset:offset}
		//err = dev.Write(offset,data)
		_ = err // handle write errors
	}
	close(eq)
	wgEncode.Wait()
	close(wq)
	wgWrite.Wait()
	//_ = <- end
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
		  	_ = out
			  if ok {
				  //log.Infof("decoded block %d: %s", qe.offset, string(out))
			  }
			  if len(errlist) > 0 {
				  log.Warning("  found errors")
				  for _, e := range errlist {
				  	_ = e
					//  log.Warningf("  [%d] error: %+v", qe.offset, e)
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
	time.Sleep(time.Second*10)
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

