package utils

import "strings"

// ParseSlowFieds parse slow log files from string
func ParseSlowFieds(fields string, slowKv map[string]string) {
	if slowKv == nil {
		slowKv = make(map[string]string)
	}
	fieldArr := strings.Split(fields, " ")
	for _, field := range fieldArr {
		kv := strings.Split(field, ":")
		if len(kv) > 0 {
			if len(kv) > 1 {
				slowKv[kv[0]] = kv[1]
			} else {
				slowKv[kv[0]] = ""
			}
		}
	}
	return
}

// CatchDbField catch db name from database fields
func CatchDbField(dbstr string) string {
	strArr := strings.Split(dbstr, " ")
	if len(strArr) != 0 {
		return strArr[0]
	}
	return ""
}
