package store

import (
	"io"
	"os"
	"sync/atomic"

	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/golang/glog"
)

const (
	sizeBitLen = 32
)

var (
	// PageSize of fielsystem
	PageSize = 4096
)

// ===== File provider =====
// every superblock has one Index File

type FileProvider struct {
	block      *Superblock
	file       string
	indexCache map[uint64]uint64
	itemCnt    uint32
	reader     *os.File
	appender   *os.File
	flake      *snowflake.Node
	mut        *sync.Mutex
}

func NewFileProvider(file string, b *Superblock) (*FileProvider, error) {
	var err error
	var p *FileProvider
	var r, w *os.File
	var node *snowflake.Node

	node, err = snowflake.NewNode(int64(b.id))
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
		block:      b,
		file:       file,
		indexCache: make(map[uint64]uint64),
		reader:     r,
		appender:   w,
		flake:      node,
		mut:        &sync.Mutex{},
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
	if _, err = fp.appender.Write(i.Bytes()); err != nil {
		return err
	}
	if err = Syncfilerange(fp.appender.Fd(), int64(fp.itemCnt)*int64(indexLen), int64(indexLen), SYNC_FILE_RANGE_WRITE); err != nil {
		return err
	}

	fp.SetCache(i.ID, i.Offset, i.Size)
	atomic.AddUint32(&fp.itemCnt, 1)

	return nil
}

func (fp *FileProvider) FindByID(id uint64) *Index {
	// var n int
	// var baseOffset int64
	var err error
	var offset, size uint32
	var i Index

	offset, size, err = fp.GetCache(id)
	if err != nil {
		glog.Error(err)
		return nil
	}
	i.ID = id
	i.Offset = offset
	i.Size = size

	return &i

	// if off, has := fp.offset[id]; has {
	// 	b := make([]byte, indexLen)
	// 	baseOffset = int64(off) * int64(indexLen)
	// 	for n < indexLen && err == nil {
	// 		var nn int
	// 		nn, err = fp.reader.ReadAt(b[n:], baseOffset+int64(n))
	// 		n += nn
	// 	}
	// 	glog.Infof("readat offset=%d read=%d err=%v", int64(off)*int64(indexLen), n, err)
	// 	if err != nil {
	// 		glog.Error(err)
	// 		return nil
	// 	}

	// 	err := i.Parse(b)
	// 	if err != nil {
	// 		glog.Error(err)
	// 		return nil
	// 	}
	//	return i

}

func (fp *FileProvider) LoadIndex() error {
	var err, idxErr error
	var tmpBytes = make([]byte, PageSize)
	var n int
	var idx Index

	fp.mut.Lock()
	defer fp.mut.Unlock()

	// fadvise sequential read
	err = Fadvise(fp.reader.Fd(), 0, int64(fp.itemCnt*uint32(indexLen)), POSIX_FADV_SEQUENTIAL)
	if err != nil {
		glog.Error(err)
		return err
	}

	// rewind to 0 point
	_, err = fp.reader.Seek(0, 0)
	if err != nil {
		glog.Error(err)
		return err
	}

	for err == nil {
		// reads at least pageSize bytes from reader
		n, err = io.ReadFull(fp.reader, tmpBytes)
		if err == io.ErrUnexpectedEOF {
			glog.Infof("readfull read=%d err=%v, maybe EOF", n, err)
		}

		var times = n / indexLen
		for i := 0; i < times; i++ {
			idxErr = idx.Parse(tmpBytes[i*indexLen : (i+1)*indexLen])
			if idxErr == nil {
				fp.SetCache(idx.ID, idx.Offset, idx.Size)
				fp.itemCnt++
			} else {
				glog.Error(idxErr)
				return idxErr
			}
		}
	}

	return nil
}

func (fp *FileProvider) SetCache(key uint64, offset, size uint32) {
	fp.indexCache[key] = encodeCacheVal(offset, size)
}

func (fp *FileProvider) GetCache(key uint64) (offset uint32, size uint32, err error) {
	if val, has := fp.indexCache[key]; has {
		offset, size = decodeCacheVal(val)
		return
	}
	err = errMissIndexCache
	return
}

func encodeCacheVal(offset, size uint32) uint64 {
	return uint64(offset)<<sizeBitLen + uint64(size)
}

func decodeCacheVal(val uint64) (uint32, uint32) {
	return uint32(val >> sizeBitLen), uint32(val)
}
