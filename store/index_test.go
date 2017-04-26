package store

import "testing"
import "github.com/stretchr/testify/assert"

func TestIndexBytes(t *testing.T) {
	var err error
	i := NewIndex(1, 1, 1)
	b := i.Bytes()
	i, err = byteToIndex(b)

	t.Log(len(b))
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), i.ID)
}
