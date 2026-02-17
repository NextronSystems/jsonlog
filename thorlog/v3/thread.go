package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type Thread struct {
	jsonlog.ObjectHeader
	ThreadId uint32     `json:"id"`
	Stack    StringList `json:"stack" jsonschema:"nullable"`
}

func (Thread) observed() {}

const typeThread = "thread"

func init() { AddLogObjectType(typeThread, &Thread{}) }

func NewThread(tid uint32) *Thread {
	return &Thread{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeThread,
		},
		ThreadId: tid,
	}
}
