package store

import "hash/crc32"

const (
	hMagicLen   = 4
	fMagicLen   = 4
	idLen       = 8
	flagLen     = 1
	sizeLen     = 4
	checksumLen = 4

	totalHeaderLen = hMagicLen + idLen + flagLen + sizeLen
	totalFooterLen = fMagicLen + checksumLen

	// align by 8 byte
	align = 8
)

var (
	headerMagic = [4]byte{0x01, 0x01, 0x01, 0x01}
	footerMagic = [4]byte{0x02, 0x01, 0x01, 0x04}

	paddings = make([][]byte, 8)

	flagNormal = byte(0)
	flagDel    = byte(1)

	crc32q = crc32.MakeTable(0xD5828281)
)

func init() {
	for index := 0; index < align; index++ {
		paddings[index] = make([]byte, index, index)
	}
}

type Needle struct {
	Header   [4]byte
	ID       uint64
	Flag     byte
	Size     uint32
	Data     []byte
	Footer   [4]byte
	Checksum uint32
	Padding  []byte
}

func getNeedlePadding(dataSize uint32) []byte {
	sizeExceptPadding := totalHeaderLen + dataSize + totalFooterLen
	d := sizeExceptPadding % align
	i := align - d
	return paddings[i]
}

func doChecksum(data []byte) uint32 {
	return crc32.Checksum(data, crc32q)
}

func (n *Needle) GetTotalSize() int {
	return totalHeaderLen + int(n.Size) + totalFooterLen + len(n.Padding)
}

func NewNeedle(data []byte) *Needle {
	dataSize := uint32(len(data))
	padding := getNeedlePadding(dataSize)
	n := &Needle{
		Header:   headerMagic,
		ID:       1,
		Flag:     flagNormal,
		Size:     dataSize,
		Data:     data,
		Footer:   footerMagic,
		Checksum: doChecksum(data),
		Padding:  padding,
	}

	return n
}
