package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type ProcessConnectionObject struct {
	jsonlog.ObjectHeader
	ProcessConnection
}

func (ProcessConnectionObject) reportable() {}

const typeProcessConnection = "process connection"

func init() { AddLogObjectType(typeProcessConnection, &ProcessConnectionObject{}) }

func NewProcessConnection() *ProcessConnectionObject {
	return &ProcessConnectionObject{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeProcessConnection,
		},
	}
}
