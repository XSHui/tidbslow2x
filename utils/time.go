package utils

import "time"

// GetSlowLogFileSuffix get suffix of slow log file
func GetSlowLogFileSuffix() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
