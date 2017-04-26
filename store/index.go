package store

import (
	"bytes"
	"encoding/binary"
	"os"

	"github.com/golang/glog"
)

const (
	idLen     = 8
	offsetLen = 4
	sizeLen   = 4
)

var (
	indexLen = idLen + offsetLen + sizeLen
)

type indexProvider interface {
	// LoadIndex during recovery from failure
	LoadIndex() error
	// Get ID of next needle
	NextID() uint64
	FindByID(id uint64) *Index
	// Persistent the index
	Append(Index) error
	Close() error
}

// IndexManager manages index. Read, write index, give id to needle
type IndexManager struct {
	store    *Store
	provider *indexProvider
}

// Index item in index file
type Index struct {
	ID     uint64
	Offset uint32
	Size   uint32 // data size of needle
}

func NewIndex(id uint64, offset, size uint32) Index {
	return Index{
		ID:     id,
		Offset: offset,
		Size:   size,
	}
}

// Bytes of index item
func (i *Index) Bytes() []byte {
	var n int
	var err error

	b := make([]byte, indexLen)
	buf := bytes.NewBuffer(b)
	tmpBytes := make([]byte, 8)
	byte32 := tmpBytes[:4]

	binary.PutUvarint(tmpBytes, i.ID)
	n, err = buf.Write(tmpBytes)

	binary.LittleEndian.PutUint32(byte32, i.Offset)
	n, err = buf.Write(byte32)

	binary.LittleEndian.PutUint32(byte32, i.Size)
	n, err = buf.Write(byte32)
	_ = n
	if err != nil {
		glog.Error(err)
	}

	return buf.Bytes()
}

func byteToIndex(b []byte) (*Index, error) {
	var i Index
	var err error
	var idx int

	if len(b) != indexLen {
		return nil, errInvalidIndexByte
	}

	i.ID, _ = binary.Uvarint(b[:idLen])
	idx += idLen

	i.Offset = binary.LittleEndian.Uint32(b[idx : idx+offsetLen])
	idx += offsetLen

	i.Size = binary.LittleEndian.Uint32(b[idx : idx+sizeLen])

	return &i, nil
}

// ===== File provider =====
// every superblock has one Index File

type FileProvider struct {
	block    *Superblock
	file     string
	offset   map[uint64]uint32 // [id]offset in index file
	reader   *os.File
	appender *os.File
}

func NewFileProvider(file string, block *Superblock) (*FileProvider, error) {
	var err error
	var p *FileProvider
	var r, w *os.File

	r, err = os.OpenFile(file, os.O_RDONLY|os.O_CREATE, 0766)
	if err != nil {
		return nil, err
	}

	w, err = os.OpenFile(file, os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	p = &FileProvider{
		block:    block,
		file:     file,
		offset:   make(map[uint64]uint32),
		reader:   r,
		appender: w,
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
