package store

import (
	"io"
	"os"
	"sync/atomic"

	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/golang/glog"
)

// ===== File provider =====
// every superblock has one Index File

type FileProvider struct {
	block    *Superblock
	file     string
	offset   map[uint64]uint32 // [id]offset in index file
	itemCnt  uint32
	reader   *os.File
	appender *os.File
	flake    *snowflake.Node
	mut      *sync.Mutex
}

func NewFileProvider(file string, block *Superblock) (*FileProvider, error) {
	var err error
	var p *FileProvider
	var r, w *os.File
	var node *snowflake.Node

	node, err = snowflake.NewNode(int64(block.id))
	if err != nil {
		return nil, err
	}

	r, err = os.OpenFile(file, os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	w, err = os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return nil, err
	}

	p = &FileProvider{
		block:    block,
		file:     file,
		offset:   make(map[uint64]uint32),
		reader:   r,
		appender: w,
		flake:    node,
		mut:      &sync.Mutex{},
	}

	return p, nil
}

// Close files
func (fp *FileProvider) Close() error {
	var err error

	err = fp.reader.Close()
	if err != nil {
		glog.Error(err)
	}

	err = fp.appender.Close()
	if err != nil {
		glog.Error(err)
	}

	return err
}

// NextID for needle
func (fp *FileProvider) NextID() uint64 {
	return uint64(fp.flake.Generate())
}

func (fp *FileProvider) Append(i *Index) error {
	var err error
	_, err = fp.appender.Write(i.Bytes())
	if err != nil {
		return err
	}
	fp.offset[i.ID] = fp.itemCnt
	atomic.AddUint32(&fp.itemCnt, 1)

	return nil
}

func (fp *FileProvider) FindByID(id uint64) *Index {
	if off, has := fp.offset[id]; has {
		b := make([]byte, indexLen)
		n, err := fp.reader.ReadAt(b, int64(off)*int64(indexLen))
		if err != nil || n < indexLen {
			glog.Error(err)
			return nil
		}
		i, err := byteToIndex(b)
		if err != nil || n < indexLen {
			glog.Error(err)
			return nil
		}
		return i
	}

	return nil
}

func (fp *FileProvider) LoadIndex() error {
	var err, idxErr error
	var pageSize = 4096
	var tmpBytes = make([]byte, pageSize)
	var n int
	var idx *Index

	fp.mut.Lock()
	defer fp.mut.Unlock()

	// rewind to 0 point
	_, err = fp.reader.Seek(0, 0)
	if err != nil {
		glog.Error(err)
		return err
	}

	for err != io.EOF {
		//TODO: read n bytes, n may be less than len(tmpBytes)
		n, err = fp.reader.Read(tmpBytes)
		var times = n / indexLen
		for i := 0; i < times; i++ {
			idx, idxErr = byteToIndex(tmpBytes[i*indexLen : (i+1)*indexLen])
			if idxErr == nil {
				fp.offset[idx.ID] = fp.itemCnt
				fp.itemCnt++
			} else {
				glog.Error(idxErr)
				return idxErr
			}
		}
	}

	return nil
}
