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
	"strings"
)

// FormatSlowLogToTidb4 format tidb-v2.x slow log to tidb-v4.x
// TODO: use template
// TODO: time format
// TODO: Query_Time unit
func FormatSlowLogToTidb4(slowKv map[string]string) string {
	if sql, ok := slowKv["Sql"]; ok {
		sql = strings.Replace(sql, "\n", "", -1)
		slowKv["Sql"] = sql
	}
	if slowKv["ProcessTime"] == "" {
		slowKv["ProcessTime"] = "0ms"
	}
	if slowKv["WaitTime"] == "" {
		slowKv["WaitTime"] = "0ms"
	}
	if slowKv["RequestCount"] == "" {
		slowKv["RequestCount"] = "0"
	}
	if slowKv["TotalKeys"] == "" {
		slowKv["TotalKeys"] = "0"
	}
	if slowKv["ProcessedKeys"] == "" {
		slowKv["ProcessedKeys"] = "0"
	}
	return fmt.Sprintf(`# Time: 20%s-%s-%sT%s+08:00
# Txn_start_ts: %s
# User@Host: %s
# Conn_ID: %s
# Query_time: %s
# Process_time: %s Wait_time: %s Request_count: %s Total_keys: %s Process_keys: %s
# Succ: %s
%s;
`, slowKv["MONTHDAY"], slowKv["MONTHNUM"], slowKv["YEAR"], slowKv["TIME"],
		slowKv["TxnStartTs"],
		slowKv["User"],
		slowKv["Con"],
		slowKv["CostTime"],
		slowKv["ProcessTime"], slowKv["WaitTime"], slowKv["RequestCount"], slowKv["TotalKeys"], slowKv["ProcessedKeys"],
		slowKv["Succ"],
		slowKv["Sql"])
}
