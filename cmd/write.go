package cmd

import (
	"github.com/XANi/flashtest/blockdev"
	"github.com/XANi/flashtest/datastore"
	"github.com/op/go-logging"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
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
		log.Errorf("file too small[%d], must be at least %d bytes",filesize, datastore.DataBlockSize)
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
		wgEncode.Add(1)
		go func() {
			encodeQWorker(eq, wq)
			wgEncode.Done()
		}()
	}
	wgWrite.Add(1)
	go func() {
		writeQWorker(wq,dev)
		wgWrite.Done()
	}()
	for i := 0; i < dataBlocks; i++ {
		offset := i * datastore.DataBlockSize
		if lastOutput.Add(time.Second * 10).Before(time.Now()) {
			lastOutput = time.Now()
			progressPct := float32(i) / float32(dataBlocks) * 100
			log.Infof("Writing block %d/%d at offset %d (%.0f)",i,dataBlocks, offset, progressPct)
		}
		eq <- encodeQ{block: []byte(fmt.Sprintf("Block %d",i)),offset:offset}
	}
	close(eq)
	wgEncode.Wait()
	close(wq)
	wgWrite.Wait()
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
		log.Errorf("file too small[%d], must be at least %d bytes",filesize,datastore.DataBlockSize)
		return
	}
	if filesize > checkedSize {
		log.Warningf("File is not aligned to data block size[%d]. Will not touch last %d bytes", datastore.DataBlockSize, filesize-checkedSize)
	}
	vq := make(chan verifyQ,5)
	end := make(chan bool,1)
	var goodCount int64
	var badCount int64
	var count int64
	go func() {
		  for qe := range vq {
		  	ok,out,errlist,err:=verifyBlock(qe.block,qe.offset)
			  atomic.AddInt64(&count,1)
		  	_ = out
			  if ok {
				  atomic.AddInt64(&goodCount,1)
				  //log.Infof("decoded block %d: %s", qe.offset, string(out))
			  }
			  if len(errlist) > 0 {
				  atomic.AddInt64(&badCount,1)
				  log.Warning("  found errors")
				  for _, e := range errlist {
				  	_ = e
					//  log.Warningf("  [%d] error: %+v", qe.offset, e)
				  }
			  }
			  if err != nil {
				  atomic.AddInt64(&badCount,1)
				  log.Errorf("error decoding block %d: %s",qe.offset,err)
			  }
		  }
		  end <- true
	} ()
	for i := 0; i < dataBlocks; i++ {
		offset := i * datastore.DataBlockSize
		data,_ := dev.Read(offset,datastore.DataBlockSize)
		vq <- verifyQ{block: data,offset: offset}
	}
	close(vq)
	log.Noticef("good/bad/total: %d/%d/%d",goodCount,badCount,count)
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

