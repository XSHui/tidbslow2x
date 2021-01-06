/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/XSHui/tidbslow2x/collector"
	"github.com/XSHui/tidbslow2x/utils"
)

const (
	slowLogFileNamePrefix = "tidb_slow_query-"
)

// collectCmd represents the collect command
var collectCmd = &cobra.Command{
	Use:   "collect [flags]",
	Short: "A slow query log collector for tidb-v2.x.",
	Long:  `A slow query log collector for tidb-v2.x.`,
	Run: func(cmd *cobra.Command, args []string) {
		// arges
		logDir, startDate, endDate := loadFlags(cmd)
		// tidb log file
		fileList := collector.ListFiles(logDir, startDate, endDate)
		fmt.Println(fileList)
		// slow log file
		slowFileDir := filepath.Join(logDir, slowLogFileNamePrefix+utils.GetSlowLogFileSuffix())

		slowFile, err := os.OpenFile(slowFileDir, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer slowFile.Close()

		// collect slow log to slow log file
		for _, file := range fileList {
			collector.SlowLogCollector(file, slowFile)
		}
	},
}

func loadFlags(cmd *cobra.Command) (logDir string, startDate string, endDate string) {
	var err error
	logDir, err = cmd.Flags().GetString("log_dir")
	if err != nil {
		fmt.Printf("log_dir error: %v", err)
		os.Exit(1)
	}
	startDate, err = cmd.Flags().GetString("start_date")
	if err != nil {
		fmt.Printf("start_date error: %v", err)
		os.Exit(1)
	}
	// TODO: start date format check
	endDate, err = cmd.Flags().GetString("end_date")
	if err != nil {
		fmt.Printf("end_date error: %v", err)
		os.Exit(1)
	}
	// TODO: end date format check
	// TODO: start and end date validity check
	return
}

func init() {
	rootCmd.AddCommand(collectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// collectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// collectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	collectCmd.Flags().StringP("log_dir", "d", "./", "tidb slow log location")
	collectCmd.Flags().StringP("start_date", "s", "", "start date")
	collectCmd.Flags().StringP("end_date", "e", "", "end date")
}
