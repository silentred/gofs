package store

import (
	"os"
)

const (
	maxSize = 1 << 28 // 256 MB , do not larger than 4 GB
	padding = 8
)

var (
	blockIDIncrement uint32

	blockMagic = []byte{0x00, 0x09, 0x01, 0x11} // 4 bytes
	blockLen   = uint32(maxSize)                // 4 bytes
)

// Superblock is the logical volume. Its id should be globally set.
// To fetch a image, it should provides superblock_id, offset.
type Superblock struct {
	id       uint32
	size     uint64 // current block size
	reader   *os.File
	writer   *os.File // mark delete
	appender *os.File // append needle

	writeOffset uint32 // cnt of needle
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
