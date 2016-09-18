package store

/*
#define _GNU_SOURCE
#include <fcntl.h>
#include <linux/falloc.h>
#include <linux/fs.h>
#include <unistd.h>
*/
import "C"
import "syscall"

const (
	O_NOATIME = syscall.O_NOATIME // do not update access time

	FALLOC_FL_DEFAULT   = uint32(0)
	FALLOC_FL_KEEP_SIZE = uint32(C.FALLOC_FL_KEEP_SIZE) // 1

	POSIX_FADV_NORMAL     = int(C.POSIX_FADV_NORMAL)     // 0
	POSIX_FADV_RANDOM     = int(C.POSIX_FADV_RANDOM)     // 1
	POSIX_FADV_SEQUENTIAL = int(C.POSIX_FADV_SEQUENTIAL) // 2
	POSIX_FADV_WILLNEED   = int(C.POSIX_FADV_WILLNEED)   // 3
	POSIX_FADV_DONTNEED   = int(C.POSIX_FADV_DONTNEED)   // 4
	POSIX_FADV_NOREUSE    = int(C.POSIX_FADV_NOREUSE)    // 5

	SYNC_FILE_RANGE_WAIT_BEFORE = int(C.SYNC_FILE_RANGE_WAIT_BEFORE) // 1
	SYNC_FILE_RANGE_WRITE       = int(C.SYNC_FILE_RANGE_WRITE)       // 2
	SYNC_FILE_RANGE_WAIT_AFTER  = int(C.SYNC_FILE_RANGE_WAIT_AFTER)  // 4
)

// Fallocate pre-allocates disk space for file
func Fallocate(fd uintptr, mode uint32, offset int64, length int64) error {
	return syscall.Fallocate(int(fd), mode, offset, length)
}

// Fadvise set file access mode
func Fadvise(fd uintptr, off int64, size int64, advise int) (err error) {
	var errno int
	if errno = int(C.posix_fadvise(C.int(fd), C.__off_t(off), C.__off_t(size), C.int(advise))); errno != 0 {
		err = syscall.Errno(errno)
	}
	return
}

// Fdatasync flushes data (without metadata, unless it is needed for the following operation) in core to disk
func Fdatasync(fd uintptr) (err error) {
	return syscall.Fdatasync(int(fd))
}

// Syncfilerange flushes data in range to disk
func Syncfilerange(fd uintptr, off int64, n int64, flags int) (err error) {
	return syscall.SyncFileRange(int(fd), off, n, flags)
}
