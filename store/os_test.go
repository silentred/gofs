package store

import (
	"fmt"
	"os"
	"testing"
)

func TestFallocate(t *testing.T) {
	f, err := os.Create("test")
	if err != nil {
		panic("cannot create test")
	}
	fi, err := f.Stat()
	if err != nil {
		panic("cannot get stat of test")
	}

	fmt.Printf("fname=%s size=%d \n", fi.Name(), fi.Size())

	err = Fallocate(f.Fd(), FALLOC_FL_DEFAULT, 0, 100)
	if err != nil {
		fmt.Println(err)
	}

	fi, _ = f.Stat()
	fmt.Printf("The file is %d bytes long \n", fi.Size())

	f.WriteAt([]byte("tttt"), 1000)
	f.WriteAt([]byte("ff"), 1000)

	//func Fadvise(fd uintptr, off int64, size int64, advise int) (err error) {
	err = Fadvise(f.Fd(), 0, 1024, POSIX_FADV_RANDOM)
	if err != nil {
		fmt.Println(err)
	}

	err = Fdatasync(f.Fd())
	if err != nil {
		fmt.Println(err)
	}

	//func Syncfilerange(fd uintptr, off int64, n int64, flags int) (err error) {
	err = Syncfilerange(f.Fd(), 1000, 4, SYNC_FILE_RANGE_WRITE)
	if err != nil {
		fmt.Println(err)
	}

	f.Close()
	os.Remove(f.Name())
}
