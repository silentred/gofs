package store

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetPadding(t *testing.T) {
	fixSize := totalHeaderLen + totalFooterLen
	dataSize := 12132
	paddingSize := len(getNeedlePadding(uint32(dataSize)))
	total := dataSize + fixSize + paddingSize
	assert.Equal(t, 0, total%align)
}

func TestNeedle(t *testing.T) {
	data := []byte{0x96, 0x95}
	n := NewNeedle(data, 1)

	b := toBytes(t, n)
	needle, err := bytesToNeedle(b)
	assert.Nil(t, err)
	assert.Equal(t, n.String(), needle.String())

}

func toBytes(t *testing.T, n *Needle) []byte {
	bytes := n.Bytes()
	totalSize := n.GetTotalSize()

	assert.Equal(t, totalSize, len(bytes))

	return bytes
}

func TestCh(t *testing.T) {
	ch := make(chan bool, 1)
	ch2 := make(chan bool, 1)

	go func() {
		time.Sleep(1 * time.Second)
		ch <- false
	}()

	time.Sleep(1 * time.Second)
	ch2 <- false

	select {
	case <-ch:
		fmt.Println(1)
	case <-ch2:
		fmt.Println(2)
	}
}
