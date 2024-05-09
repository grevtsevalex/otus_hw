package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyCustom(t *testing.T) {
	f, err := os.CreateTemp("", "sample")
	if err != nil {
		return
	}
	defer f.Close()

	stringContent := "Hello, Otus!"

	f.WriteString(stringContent)
	newFileName := "/tmp/file1"

	os.Remove(newFileName)

	t.Run("bad file case", func(t *testing.T) {
		err := Copy("/dev/urandom", newFileName, 0, 12)
		os.Remove(newFileName)
		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual error %q", err)
	})

	t.Run("big offset case", func(t *testing.T) {
		err := Copy(f.Name(), newFileName, 13, 12)
		os.Remove(newFileName)
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual error %q", err)
	})

	t.Run("simple case", func(t *testing.T) {
		err := Copy(f.Name(), newFileName, 0, 0)
		require.NoError(t, err)
		require.FileExists(t, newFileName)

		nf, err := os.Open(newFileName)
		if err != nil {
			return
		}
		defer nf.Close()
		require.NoError(t, err)

		buf := make([]byte, 12)
		_, err = nf.Read(buf)

		require.NoError(t, err)
		require.Equal(t, []byte(stringContent), buf)
		os.Remove(newFileName)
	})

	t.Run("simple case with offset", func(t *testing.T) {
		err := Copy(f.Name(), newFileName, 2, 0)
		require.NoError(t, err)
		require.FileExists(t, newFileName)

		nf, err := os.Open(newFileName)
		if err != nil {
			return
		}
		defer nf.Close()
		require.NoError(t, err)

		buf := make([]byte, 10)
		_, err = nf.Read(buf)

		require.NoError(t, err)
		require.Equal(t, []byte(stringContent)[2:], buf)
		os.Remove(newFileName)
	})

	t.Run("simple case with limit", func(t *testing.T) {
		limit := 4
		err := Copy(f.Name(), newFileName, 2, int64(limit))
		require.NoError(t, err)
		require.FileExists(t, newFileName)

		nf, err := os.Open(newFileName)
		if err != nil {
			return
		}
		defer nf.Close()
		require.NoError(t, err)

		buf := make([]byte, limit)
		_, err = nf.Read(buf)

		require.NoError(t, err)
		require.Equal(t, []byte(stringContent)[2:6], buf)
		os.Remove(newFileName)
	})

	t.Run("simple case with big limit", func(t *testing.T) {
		limit := 40
		err := Copy(f.Name(), newFileName, 0, int64(limit))
		require.NoError(t, err)
		require.FileExists(t, newFileName)

		nf, err := os.Open(newFileName)
		if err != nil {
			return
		}
		defer nf.Close()
		require.NoError(t, err)

		buf := make([]byte, 12)
		_, err = nf.Read(buf)

		require.NoError(t, err)
		require.Equal(t, []byte(stringContent), buf)
		os.Remove(newFileName)
	})
}

func TestCopyBigFiles(t *testing.T) {
	f, err := os.Open("./testdata/input.txt")
	if err != nil {
		return
	}

	newFileName := "/tmp/file1"
	os.Remove(newFileName)

	t.Run("simple case", func(t *testing.T) {
		err := Copy(f.Name(), newFileName, 0, 0)
		require.NoError(t, err)
		require.FileExists(t, newFileName)

		expected, err := os.Stat("./testdata/out_offset0_limit0.txt")
		require.NoError(t, err)

		dest, err := os.Stat(newFileName)
		require.NoError(t, err)

		require.NoError(t, err)
		require.Equal(t, expected.Size(), dest.Size())
		os.Remove(newFileName)
	})

	t.Run("simple case with offset and limit", func(t *testing.T) {
		err := Copy(f.Name(), newFileName, 100, 1000)
		require.NoError(t, err)
		require.FileExists(t, newFileName)

		expected, err := os.Stat("./testdata/out_offset100_limit1000.txt")
		require.NoError(t, err)

		dest, err := os.Stat(newFileName)
		require.NoError(t, err)

		require.NoError(t, err)
		require.Equal(t, expected.Size(), dest.Size())
		os.Remove(newFileName)
	})

	t.Run("simple case with big limit", func(t *testing.T) {
		err := Copy(f.Name(), newFileName, 0, 10000)
		require.NoError(t, err)
		require.FileExists(t, newFileName)

		expected, err := os.Stat("./testdata/out_offset0_limit10000.txt")
		require.NoError(t, err)

		dest, err := os.Stat(newFileName)
		require.NoError(t, err)

		require.NoError(t, err)
		require.Equal(t, expected.Size(), dest.Size())
		os.Remove(newFileName)
	})

	t.Run("simple case with big offset", func(t *testing.T) {
		err := Copy(f.Name(), newFileName, 6000, 1000)
		require.NoError(t, err)
		require.FileExists(t, newFileName)

		expected, err := os.Stat("./testdata/out_offset6000_limit1000.txt")
		require.NoError(t, err)

		dest, err := os.Stat(newFileName)
		require.NoError(t, err)

		require.NoError(t, err)
		require.Equal(t, expected.Size(), dest.Size())
		os.Remove(newFileName)
	})
}
