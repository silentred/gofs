package store

import "os"

const (
	maxSize = 1 << 28 // 256 MB
	padding = 8
)

type Superblock struct {
	rf      *os.File
	wf      *os.File
	wOffset uint64
}
