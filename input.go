package main

import (
	"errors"
	"io/fs"
	"os"
	"path"
)

var (
	ErrArgumentCount  = errors.New("incorrect number of arguments")
	ErrFileNotExist   = errors.New("file does not exist")
	ErrFileNotRegular = errors.New("file is not a regular file")
	ErrFileExtension  = errors.New("incorrect file extension")
	ErrFileEmpty      = errors.New("file is empty")
)

func GetInputPath() (string, error) {
	if len(os.Args[1:]) != 1 {
		return "", ErrArgumentCount
	}

	return os.Args[1], nil
}

func GetInputFile(inputPath string) (*os.File, error) {
	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, ErrFileNotExist
		}

		return nil, err
	}

	if !fileInfo.Mode().IsRegular() {
		return nil, ErrFileNotRegular
	}

	if ext := path.Ext(fileInfo.Name()); ext != ".txt" {
		return nil, ErrFileExtension
	}

	if fileInfo.Size() == 0 {
		return nil, ErrFileEmpty
	}

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}

	return inputFile, nil
}
