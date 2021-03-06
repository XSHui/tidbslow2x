/*
Copyright © 2021 taylor.xiong

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package collector

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/XSHui/tidbslow2x/utils"
)

// ListFiles list all valid tidb log file in dir
func ListFiles(dir string, start, end string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	fileList := []string{}
	for _, file := range files {
		// TODO: deal with *.log only
		if file.Name() == "tidb_stderr.log" || file.Name() == "tidbslow2x" {
			continue
		}
		if strings.HasPrefix(file.Name(), tidbLogPrefix) {
			if start == "" && end == "" {
				if file.Name() >= fmt.Sprintf("%s-%s", tidbLogPrefix, "2006-01-02T15-04-05") {
					fileList = append(fileList, filepath.Join(dir, file.Name()))
				}
			} else if start == "" && end != "" {
				if file.Name() == "tidb.log" {
					if end > utils.GetCurrentTimeStr() {
						fileList = append(fileList, filepath.Join(dir, file.Name()))
					}
				} else if file.Name() <= fmt.Sprintf("%s-%s", tidbLogPrefix, end) {
					fileList = append(fileList, filepath.Join(dir, file.Name()))
				}
			} else if start != "" && end == "" {
				if file.Name() >= fmt.Sprintf("%s-%s", tidbLogPrefix, start) {
					fileList = append(fileList, filepath.Join(dir, file.Name()))
				}
			} else {
				if file.Name() == "tidb.log" {
					if end > utils.GetCurrentTimeStr() {
						fileList = append(fileList, filepath.Join(dir, file.Name()))
					}
				} else if file.Name() >= fmt.Sprintf("%s-%s", tidbLogPrefix, start) &&
					file.Name() <= fmt.Sprintf("%s-%s", tidbLogPrefix, end) {
					fileList = append(fileList, filepath.Join(dir, file.Name()))
				}
			}
		}
	}
	return fileList
}
