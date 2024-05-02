package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fi, err := os.Stat(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}

	if !fi.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	size := fi.Size()

	if offset > size {
		return ErrOffsetExceedsFileSize
	}

	if limit > size {
		limit = 0
	}

	source, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer destination.Close()

	var bytesToCopy int64
	bytesToCopy = limit
	if limit == 0 {
		bytesToCopy = size - offset
	}

	if limit > (size - offset) {
		bytesToCopy = size - offset
	}

	bar := pb.StartNew(int(bytesToCopy))

	buf := make([]byte, bytesToCopy)

	var bufOffset int64
	for bufOffset < bytesToCopy {
		read, err := source.ReadAt(buf[bufOffset:], offset)
		bufOffset += int64(read)
		bar.Add64(int64(read))
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	_, err = destination.Write(buf)
	if err != nil {
		return err
	}

	bar.Finish()

	return nil
}
