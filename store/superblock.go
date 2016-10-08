package store

import "os"

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
	id     uint32
	size   uint64 // block size
	reader *os.File
	writer *os.File

	writeOffset uint64
	writeEnable bool
	//index
	index *IndexManager
	store *Store
}

// NewSuperblock creates new a logical volume. Its id is fetched from etcd
func NewSuperblock() *Superblock {
	return nil
}

func AppendNeedle() {
}

func ReadNeedleByKey() {

}

func MarkNeedleDeleted() {

}

func Init() {

}

func Terminate() {

}

func stopWrite() {

}

func flush() {

}
