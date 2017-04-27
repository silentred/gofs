package store

import (
	"os"
)

const (
	padding = 8
)

var (
	// MaxBlockSize is the max size of superblock
	// considering the need of compacting and backup, this value should not be too large.
	// default is 256 MB, should be less than 8*4=32GB (theoretical max value)
	MaxBlockSize = 1 << 28

	blockMagic     = []byte{0x00, 0x09, 0x01, 0x11} // 4 bytes
	maxLogicOffset = uint32(MaxBlockSize / align)
)

// Superblock is the logical volume. Its id should be globally set.
// To fetch a image, it should provides superblock_id, offset.
type Superblock struct {
	id       uint32
	size     uint64 // current block size
	file     string // file name
	reader   *os.File
	writer   *os.File // mark delete
	appender *os.File // append needle

	writeOffset uint32 // logical offset = real byte offset / align
	writeEnable bool
	deletedSize uint32 // the total size of deleted needle
	//index
	index indexProvider
	store *Store
}

// NewSuperblock creates new a logical volume.
func NewSuperblock() *Superblock {
	//
	return nil
}

func (sb *Superblock) AppendNeedle() error {
	return nil
}

func (sb *Superblock) ReadNeedleByID() ([]byte, error) {
	return nil, nil
}

func (sb *Superblock) MarkNeedleDeleted() error {
	return nil
}

// MakeIndex reads through superblock file to generate a new index file
func (sb *Superblock) MakeIndex() error {
	return nil
}

func Init() {

}

func Shutdown() {

}

func stopWrite() {

}

func flush() {

}
