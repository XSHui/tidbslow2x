package collector

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// ListFiles list all valid tidb log file in dir
func ListFiles(dir string, start, end string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	fileList := []string{}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), tidbLogPrefix) {
			// TODO: start and end optimize
			if file.Name() >= fmt.Sprintf("%s-%s", tidbLogPrefix, start) &&
				file.Name() <= fmt.Sprintf("%s-%s", tidbLogPrefix, end) {
				fileList = append(fileList, filepath.Join(dir, file.Name()))
			}
		}
	}
	return fileList
}
