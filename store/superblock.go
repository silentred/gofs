package store

import "os"

const (
	maxSize = 1 << 28 // 256 MB , do not larger than 4 GB
	padding = 8
)

var (
	blockIDIncrement uint32

	blockMagic = []byte{0x01, 0x01, 0x01, 0x01} // 4 bytes
	blockLen   = uint32(maxSize)                // 4 bytes

)

type Superblock struct {
	id     uint32
	size   uint32 // block size
	reader *os.File
	writer *os.File

	writeOffset uint64
	writeEnable bool
	//index
}

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
