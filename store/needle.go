package store

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/crc32"
)

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
	headerMagic = []byte{0x01, 0x09, 0x08, 0x09}
	footerMagic = []byte{0x01, 0x01, 0x01, 0x04}

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
	Header   []byte
	ID       uint64
	Flag     byte
	Size     uint32
	Data     []byte
	Footer   []byte
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

func NewNeedle(data []byte, id uint64) *Needle {
	dataSize := uint32(len(data))
	padding := getNeedlePadding(dataSize)
	n := &Needle{
		Header:   headerMagic,
		ID:       id,
		Flag:     flagNormal,
		Size:     dataSize,
		Data:     data,
		Footer:   footerMagic,
		Checksum: doChecksum(data),
		Padding:  padding,
	}

	return n
}

func (n *Needle) GetTotalSize() int {
	return totalHeaderLen + int(n.Size) + totalFooterLen + len(n.Padding)
}

// ToBytes converts needle to []byte
func (n *Needle) Bytes() []byte {
	b := make([]byte, 0, n.GetTotalSize())
	buf := bytes.NewBuffer(b)
	tmpBytes := make([]byte, 8)

	buf.Write(n.Header)

	binary.PutUvarint(tmpBytes, n.ID)
	buf.Write(tmpBytes)
	buf.Write([]byte{n.Flag})

	byte32 := tmpBytes[:4]
	binary.LittleEndian.PutUint32(byte32, n.Size)
	buf.Write(byte32)

	buf.Write(n.Data)
	buf.Write(n.Footer)

	binary.LittleEndian.PutUint32(byte32, n.Checksum)
	buf.Write(byte32)
	buf.Write(n.Padding)

	return buf.Bytes()
}

func (n *Needle) String() string {
	bytes, err := json.Marshal(n)
	if err != nil {
		fmt.Println(err)
	}

	return string(bytes)
}

func bytesToNeedle(b []byte) (n *Needle, err error) {
	n = new(Needle)
	i := 0

	n.Header = b[i : i+4]
	i += 4

	idByte := b[i : i+8]
	i += 8
	id, _ := binary.Uvarint(idByte)
	n.ID = id

	n.Flag = b[i]
	i++

	n.Size = binary.LittleEndian.Uint32(b[i : i+4])
	i += 4

	n.Data = b[i : i+int(n.Size)]
	i += int(n.Size)

	n.Footer = b[i : i+4]
	i += 4

	n.Checksum = binary.LittleEndian.Uint32(b[i : i+4])
	i += 4

	n.Padding = b[i:]

	return
}
