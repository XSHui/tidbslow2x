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
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"sync"

	"github.com/XSHui/tidbslow2x/utils"
	"github.com/vjeantet/grok"
)

//var pattern string = `%{DATA:Date} %{DATA:Time} %{DATA:File} \[%{LOGLEVEL:Level}\] \[SLOW_QUERY\] cost_time:%{DATA:CostTime}( process_time:%{DATA:ProcessTime} wait_time:%{DATA:WaitTime}| process_time:%{DATA:ProcessTime}| wait_time:%{DATA:WaitTime}) request_count:%{DATA:RequestCount}( total_keys:%{DATA:TotalKeys} processed_keys:%{DATA:ProcessedKeys}| total_keys:%{DATA:TotalKeys}) succ:%{DATA:Succ} con:%{DATA:Con} user:%{GREEDYDATA:User} txn_start_ts:%{DATA:TxnStartTs} database:(%{DATA:Database}) (table_ids:(%{DATA:TableIds}),index_ids:(%{DATA:IndexIds})|table_ids:(%{DATA:TableIds})),sql:%{GREEDYDATA:Sql}`
//var pattern string = `%{DATA:Date} %{DATA:Time} %{DATA:File} \[%{LOGLEVEL:Level}\] \[SLOW_QUERY\] %{GREEDYDATA:Fields},sql:%{GREEDYDATA:Sql}`
var pattern string = `%{DATA:Date} %{DATA:Time} %{DATA:File} \[%{LOGLEVEL:Level}\] \[SLOW_QUERY\] %{GREEDYDATA:Fields} database:%{GREEDYDATA:DatabaseFields}sql:%{GREEDYDATA:Sql}`

// SlowLogCollector collect slow log from tidb log
func SlowLogCollector(fileName string, slowLogFile *os.File) error {
	// open file
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	// collect
	linesPool := sync.Pool{New: func() interface{} {
		lines := make([]byte, chuckSize)
		return lines
	}}

	stringPool := sync.Pool{New: func() interface{} {
		lines := ""
		return lines
	}}

	reader := bufio.NewReader(file)
	var wg sync.WaitGroup

	for {
		buf := linesPool.Get().([]byte)
		n, err := reader.Read(buf)
		buf = buf[:n]
		if n == 0 {
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(err)
				break
			}
			return err
		}
		nextUntillNewline, err := reader.ReadBytes('\n')
		if err != io.EOF {
			buf = append(buf, nextUntillNewline...)
		}
		wg.Add(1)
		go func() {
			processChunk(buf, &linesPool, &stringPool, slowLogFile)
			wg.Done()
		}()
	}

	wg.Wait()
	return nil
}

func processChunk(chunk []byte, linesPool, stringPool *sync.Pool, slowLogFile *os.File) {
	logs := stringPool.Get().(string)
	logs = string(chunk)
	linesPool.Put(chunk)
	logsSlice := strings.Split(logs, "\n")
	stringPool.Put(logs)

	n := len(logsSlice)
	noOfThread := n / rowCount
	if n%rowCount != 0 {
		noOfThread++
	}

	var wg2 sync.WaitGroup
	for i := 0; i < (noOfThread); i++ {
		wg2.Add(1)
		go func(s int, e int, slowLogFile *os.File) {
			defer wg2.Done() //to avaoid deadlocks
			for i := s; i < e; i++ {
				text := logsSlice[i]
				if len(text) == 0 {
					continue
				}
				lineProcess(text, slowLogFile)
			}
		}(i*rowCount, int(math.Min(float64((i+1)*rowCount), float64(len(logsSlice)))), slowLogFile)
	}

	wg2.Wait()
	logsSlice = nil
}

// TODO: match slow log to tidb-v4.0.0
func lineProcess(line string, slowLogFile *os.File) {
	if strings.Contains(line, slowLogLabel) {
		rok, _ := grok.New()
		rokMap, err := rok.Parse(pattern, line)
		if err != nil || len(rokMap) == 0 {
			fmt.Println(line)
			return
		}
		utils.ParseSlowFieds(rokMap["Fields"], rokMap)
		rokMap["database"] = utils.CatchDbField(rokMap["DatabaseFields"])
		// write to slow log file
		_, err = slowLogFile.Write([]byte(FormatSlowLogToTidb4(rokMap)))
		if err != nil {
			panic(err)
		}
	}
}
