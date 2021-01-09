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

	"github.com/XSHui/tidbslow2x/utils"
)

// FormatSlowLogToTidb4 format tidb-v2.x slow log to tidb-v4.x
func FormatSlowLogToTidb4(slowKv map[string]string) string {
	if date, ok := slowKv["Date"]; ok {
		slowKv["Date"] = strings.Replace(date, "/", "-", -1)
	}
	if sql, ok := slowKv["Sql"]; ok {
		slowKv["Sql"] = strings.Replace(sql, "\n", "", -1)
	}
	return fmt.Sprintf(`# Time: %sT%s+08:00
# Txn_start_ts: %s
# User@Host: %s
# Conn_ID: %s
# Query_time: %s
# Process_time: %s Wait_time: %s Request_count: %s Total_keys: %s Process_keys: %s
# DB: %s
# Succ: %s
%s;
`, slowKv["Date"], slowKv["Time"],
		slowKv["txn_start_ts"],
		slowKv["user"],
		slowKv["con"],
		utils.SlowLogTimeToSecond(slowKv["cost_time"]),
		utils.SlowLogTimeToSecond(slowKv["process_time"]),
		utils.SlowLogTimeToSecond(slowKv["wait_time"]),
		slowKv["request_count"],
		slowKv["total_keys"],
		slowKv["processed_keys"],
		slowKv["database"],
		slowKv["succ:true"],
		slowKv["Sql"])
}
