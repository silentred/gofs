package store

import "os"

const (
	O_NOATIME = 0 // no such option

	FALLOC_FL_DEFAULT   = uint32(0)
	FALLOC_FL_KEEP_SIZE = uint32(1) // 1

	POSIX_FADV_NORMAL     = 0 // 0
	POSIX_FADV_RANDOM     = 1 // 1
	POSIX_FADV_SEQUENTIAL = 2 // 2
	POSIX_FADV_WILLNEED   = 3 // 3
	POSIX_FADV_DONTNEED   = 4 // 4
	POSIX_FADV_NOREUSE    = 5 // 5

	SYNC_FILE_RANGE_WAIT_BEFORE = 1 // 1
	SYNC_FILE_RANGE_WRITE       = 2 // 2
	SYNC_FILE_RANGE_WAIT_AFTER  = 4 // 4
)

func Fallocate(fd uintptr, mode uint32, offset int64, length int64) error {
	return nil
}

func Fadvise(fd uintptr, off int64, size int64, advise int) error {
	return nil
}

func Fdatasync(fd uintptr) error {
	file := os.NewFile(fd, "test")
	return file.Sync()
}

func Syncfilerange(fd uintptr, off int64, n int64, flags int) error {
	file := os.NewFile(fd, "test")
	return file.Sync()
}
