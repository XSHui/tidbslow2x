package collector

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"sync"
)

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
		// write to slow log file
		_, err := slowLogFile.Write([]byte(line + "\n"))
		if err != nil {
			panic(err)
		}
	}
}
