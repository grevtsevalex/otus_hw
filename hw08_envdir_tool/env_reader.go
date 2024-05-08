package main

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"strings"
)

type Environment map[string]EnvValue

var (
	ErrBadDir         = errors.New("can not open directory")
	ErrBadEnvFileName = errors.New("wrong filename in directory")
)

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return Environment{}, ErrBadDir
	}

	result := make(Environment, len(files))

	for _, f := range files {
		fi, err := f.Info()
		if err != nil {
			return result, err
		}

		if strings.Contains(fi.Name(), "=") {
			return result, ErrBadEnvFileName
		}

		var envValue []byte

		if fi.Size() == 0 {
			result[fi.Name()] = EnvValue{
				Value:      "",
				NeedRemove: true,
			}
			continue
		}

		fd, err := os.Open(dir + "/" + fi.Name())
		if err != nil {
			return result, err
		}
		defer fd.Close()

		sc := bufio.NewScanner(fd)
		for sc.Scan() {
			envValue = sc.Bytes()
			break
		}

		envValue = bytes.ReplaceAll(envValue, []byte("\x00"), []byte("\n"))

		if err := sc.Err(); err != nil {
			return result, err
		}

		result[f.Name()] = EnvValue{
			Value:      strings.TrimRight(string(envValue), " "),
			NeedRemove: false,
		}
	}
	return result, nil
}
