package thorlog

import (
	"fmt"
	"time"

	"github.com/NextronSystems/jsonlog"
)

// PrefetchInfo contains information about a Windows Prefetch file.
//
// Prefetch files are used by Windows to speed up the startup of applications.
// They contain information about an executable such as:
// - The path to the executable
// - The times the executable was run
// - The number of times the executable was run
// - Files accessed by the executable
//
// Prefetch files are located in the C:\Windows\Prefetch directory and have the .pf file extension.
// They rotate, meaning that older prefetch files are deleted when the number of prefetch files exceeds a certain limit.
type PrefetchInfo struct {
	jsonlog.ObjectHeader
	Executable     *File          `json:"executable" textlog:"executable,expand"`
	ExecutionTimes ExecutionTimes `json:"execution_times" textlog:",expand"`
	ExecutionCount int            `json:"execution_count" textlog:"execution_count"`
	AccessedFiles  []string       `json:"accessed_files" textlog:"-"`
}

func (PrefetchInfo) reportable() {}

type ExecutionTimes []time.Time

func (e ExecutionTimes) MarshalTextLog(t jsonlog.TextlogFormatter) jsonlog.TextlogEntry {
	// Only include the most recent execution time in the textlog
	if len(e) == 0 {
		return nil
	}
	var formattedLastTime string
	if t.FormatValue != nil {
		formattedLastTime = t.FormatValue(e[0], nil)
	} else {
		formattedLastTime = fmt.Sprint(e[0])
	}
	return jsonlog.TextlogEntry{
		{Key: "last_execution_time", Value: formattedLastTime},
	}
}

const typePrefetchInfo = "prefetch info"

func init() { AddLogObjectType(typePrefetchInfo, &PrefetchInfo{}) }

func NewPrefetchInfo() *PrefetchInfo {
	return &PrefetchInfo{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typePrefetchInfo,
		},
	}
}
