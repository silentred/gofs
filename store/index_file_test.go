package store

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileProvider_Append(t *testing.T) {
	var err error
	var file = "/tmp/sblock_0001.index"
	fp := getFileProvider(t, file)
	defer func() {
		fp.Close()
		os.Remove(file)
	}()

	for index := 0; index < 100; index++ {
		err = fp.Append(&Index{uint64(index), 1, 1})
		assert.NoError(t, err)
	}
	assert.Equal(t, uint32(100), fp.itemCnt)
}

func TestFileProvider_LoadIndex(t *testing.T) {
	var err error
	var writeRnd = 4000
	var file = "/tmp/sblock_0001.index"
	fp := getFileProvider(t, file)

	// append
	for index := 0; index < writeRnd; index++ {
		err = fp.Append(&Index{uint64(index), 1, 1})
		assert.NoError(t, err)
	}
	fp.Close()

	fp = getFileProvider(t, file)
	defer os.Remove(file)
	defer fp.Close()

	// reload
	err = fp.LoadIndex()
	assert.NoError(t, err)
	assert.Equal(t, uint32(writeRnd), fp.itemCnt)
	assert.Equal(t, writeRnd, len(fp.indexCache))

	index := fp.FindByID(uint64(writeRnd / 2))
	assert.Equal(t, uint64(writeRnd/2), index.ID)
}

func getFileProvider(t *testing.T, file string) *FileProvider {
	var block = &Superblock{id: 1}
	fp, err := NewFileProvider(file, block)
	assert.NoError(t, err)
	return fp
}
