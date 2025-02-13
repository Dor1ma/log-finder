package mmap

import (
	"os"

	"golang.org/x/sys/unix"
)

func MapFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	size := int(info.Size())
	if size == 0 {
		return nil, nil
	}

	data, err := unix.Mmap(int(file.Fd()), 0, size, unix.PROT_READ, unix.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func Unmap(data []byte) error {
	return unix.Munmap(data)
}
