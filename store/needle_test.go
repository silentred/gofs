package store

import (
	"testing"

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
	data := []byte("test test")
	n := NewNeedle(data, 1)

	b := toBytes(t, n)
	needle, err := bytesToNeedle(b)
	t.Logf("needle: %s", needle.String())
	assert.Nil(t, err)
	assert.Equal(t, n.String(), needle.String())
}

func toBytes(t *testing.T, n *Needle) []byte {
	bytes := n.Bytes()
	totalSize := n.GetTotalSize()

	assert.Equal(t, totalSize, len(bytes))

	return bytes
}
