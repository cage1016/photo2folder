package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func AddCountToName(folderPath string) {
	re := regexp.MustCompile(`\d+`)
	directoryList := make(map[string]string, 0)
	err := filepath.Walk(folderPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() && re.Match([]byte(info.Name())) && !strings.HasSuffix(folderPath, info.Name()) {
				files, _ := ioutil.ReadDir(path)
				directoryList[path] = fmt.Sprintf("%s_%d", path, len(files))
			}
			return nil
		})

	if err != nil {
		fmt.Print(err)
		return
	}

	for k, v := range directoryList {
		err := os.Rename(k, v)
		if err != nil {
			fmt.Print(err)
		} else {
			fmt.Sprintf("%s → %s\n", k, v)
		}
	}
}
