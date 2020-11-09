package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func RemoveEmptyFolder(folderPath string) {
	directoryList := make([]string, 0)
	err := filepath.Walk(folderPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				directoryList = append(directoryList, path)
			}
			return nil
		})

	for _, d := range directoryList {
		isEmpty, err := IsDirEmpty(d)
		if err == nil && isEmpty {
			os.Remove(d)
		}
	}

	if err != nil {
		fmt.Print(err)
	}
}
