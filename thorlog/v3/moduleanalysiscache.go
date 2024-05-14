package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type PSMacEntry struct {
	jsonlog.ObjectHeader

	Path    string `json:"path" textlog:"path"`
	Command string `json:"command" textlog:"command"`
}

const typeModuleAnalysisCacheEntry = "PowerShell module analysis cache module entry"

func init() { AddLogObjectType(typeModuleAnalysisCacheEntry, &PSMacEntry{}) }

func NewModuleAnalysisCacheEntry() *PSMacEntry {
	return &PSMacEntry{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeModuleAnalysisCacheEntry,
		},
	}
}
