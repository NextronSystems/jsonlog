package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
)

type LogObjectHeader = jsonlog.ObjectHeader

type ShellbagEntry struct {
	LogObjectHeader

	Path       string    `json:"path" textlog:"path"`
	Name       string    `json:"name" textlog:"name"`
	DateAccess time.Time `json:"date_access" textlog:"date_access"`
	Hive       string    `json:"hive" textlog:"hive"`
}

func (ShellbagEntry) reportable() {}

const typeShellbagEntry = "shellbag entry"

func init() { AddLogObjectType(typeShellbagEntry, &ShellbagEntry{}) }

func NewShellbagEntry() *ShellbagEntry {
	return &ShellbagEntry{
		LogObjectHeader: LogObjectHeader{
			Type: typeShellbagEntry,
		},
	}
}
