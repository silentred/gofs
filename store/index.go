package store

import (
	"bytes"
	"encoding/binary"

	"github.com/golang/glog"
)

const (
	idxIDLen     = 8
	idxOffsetLen = 4
	idxSizeLen   = 4
)

var (
	indexLen = idxIDLen + idxOffsetLen + idxSizeLen
)

type indexProvider interface {
	// LoadIndex to memery during recovery from failure
	LoadIndex() error
	// Get ID of next needle
	NextID() uint64
	FindByID(id uint64) *Index
	// Persistent the index
	Append(Index) error
	Close() error
}

// IndexManager manages index. Read, write index, give id to needle
// regernate index file when compacting superblock(cleaning deleted needle)
type IndexManager struct {
	block    *Superblock
	provider *indexProvider
}

func NewIndexManager() {

}

// Index item in index file
type Index struct {
	ID     uint64
	Offset uint32
	Size   uint32 // data size of needle
}

func NewIndex(id uint64, offset, size uint32) *Index {
	return &Index{
		ID:     id,
		Offset: offset,
		Size:   size,
	}
}

// Bytes of index item
func (i *Index) Bytes() []byte {
	var n int
	var err error

	b := make([]byte, 0, indexLen)
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
	var idx int

	if len(b) != indexLen {
		return nil, errInvalidIndexByte
	}

	i.ID, _ = binary.Uvarint(b[:idLen])
	idx += idLen

	i.Offset = binary.LittleEndian.Uint32(b[idx : idx+idxOffsetLen])
	idx += idxOffsetLen

	i.Size = binary.LittleEndian.Uint32(b[idx : idx+sizeLen])

	return &i, nil
}
