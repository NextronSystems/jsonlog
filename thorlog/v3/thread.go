package thorlog

import (
	"fmt"

	"github.com/NextronSystems/jsonlog"
)

type Thread struct {
	jsonlog.ObjectHeader
	ThreadId uint32     `json:"id"`
	Stack    StringList `json:"stack" jsonschema:"nullable"`
}

func (Thread) reportable() {}

const typeThread = "thread"

func init() { AddLogObjectType(typeThread, &Thread{}) }

func NewThread(tid uint32) *Thread {
	return &Thread{
		ObjectHeader: jsonlog.ObjectHeader{
			Type:    typeThread,
			Summary: fmt.Sprintf("Thread %d", tid),
		},
		ThreadId: tid,
	}
}
