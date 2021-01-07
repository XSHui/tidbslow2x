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
package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// GetSlowLogFileSuffix get suffix of slow log file
func GetSlowLogFileSuffix() string {
	return time.Now().Format("2006-01-02T15-04-05")
}

// SlowLogTimeToSecond convert slow log time to seconds
// 5.5s  --> 5.5
// 600ms --> 0.6
// 4m31.785s --> 271.785
// Todo: opt
func SlowLogTimeToSecond(orgTimeStr string) string {
	if orgTimeStr == "" {
		return "0"
	} else if strings.Contains(orgTimeStr, "ms") {
		return orgTimeStr[:len(orgTimeStr)-2]
	} else if strings.Contains(orgTimeStr, "m") {
		timeArr := strings.Split(orgTimeStr, "m")
		if len(timeArr) != 2 {
			fmt.Println("Error Time: ", orgTimeStr)
			return "0"
		}
		min, err := strconv.ParseFloat(timeArr[0], 32)
		if err != nil {
			fmt.Println(err, "Error Time: ", orgTimeStr)
			return "0"
		}
		secStr := timeArr[1][:len(timeArr[1])-1]
		sec, err := strconv.ParseFloat(secStr, 32)
		if err != nil {
			fmt.Println(err, "Error Time: ", orgTimeStr)
			return "0"
		}
		return fmt.Sprintf("%f", min*60+sec)
	} else {
		return orgTimeStr[:len(orgTimeStr)-1]
	}
}
