package store

import (
	"bytes"
	"encoding/binary"
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

// Needle represents the object item in Superblock
type Needle struct {
	Header   []byte // magic header
	ID       uint64
	Flag     byte   // delete flag
	Size     uint32 // data size
	Data     []byte
	Footer   []byte // magic footer
	Checksum uint32 // crc32(data)
	Padding  []byte
}

func getNeedlePadding(dataSize uint32) []byte {
	sizeExceptPadding := totalHeaderLen + dataSize + totalFooterLen
	d := sizeExceptPadding % align
	i := align - d
	return paddings[i]
}

// calculate total size of needle from dataSize
func totalSize(dataSize int) int {
	return dataSize + len(getNeedlePadding(uint32(dataSize))) + totalHeaderLen + totalFooterLen
}

func checksum(data []byte) uint32 {
	return crc32.Checksum(data, crc32q)
}

// NewNeedle returns a new needle
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
		Checksum: checksum(data),
		Padding:  padding,
	}

	return n
}

// GetTotalSize get the size of the needle
func (n *Needle) GetTotalSize() int {
	return totalHeaderLen + int(n.Size) + totalFooterLen + len(n.Padding)
}

// Bytes converts needle to []byte
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
	f := `[needle] id=%d size=%d flag=%d checksum=%d t_size=%d`
	return fmt.Sprintf(f, n.ID, n.Size, n.Flag, n.Checksum, n.GetTotalSize())
}

func (n *Needle) Parse(b []byte) error {
	var err error
	var i int

	if len(b)%align != 0 {
		return errNotAlign
	}

	if len(b) < (totalFooterLen + totalHeaderLen) {
		return errInvalidNeddleByte
	}

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

	if len(b)-i < (int(n.Size) + fMagicLen + checksumLen) {
		return errInvalidNeddleByte
	}

	n.Data = b[i : i+int(n.Size)]
	i += int(n.Size)

	n.Footer = b[i : i+4]
	i += 4

	n.Checksum = binary.LittleEndian.Uint32(b[i : i+4])
	i += 4

	n.Padding = b[i:]

	return err
}
