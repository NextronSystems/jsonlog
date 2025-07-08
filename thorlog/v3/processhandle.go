package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type ProcessHandle struct {
	jsonlog.ObjectHeader

	Name   string `json:"name" textlog:"name"`
	Handle uint64 `json:"handle" textlog:"handle,omitempty"`
	Type   string `json:"type,omitempty" textlog:"type,omitempty"`
}

func (ProcessHandle) reportable() {}

const typeProcessHandle = "process handle"

func init() { AddLogObjectType(typeProcessHandle, &ProcessHandle{}) }

func NewProcessHandle() *ProcessHandle {
	return &ProcessHandle{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeProcessHandle,
		},
	}
}
