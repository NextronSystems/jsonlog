package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type ProcessConnectionObject struct {
	jsonlog.ObjectHeader
	Process           *Process          `json:"process" textlog:"process,expand"`
	ConnectionDetails ProcessConnection `json:"connection" textlog:",expand"`
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
