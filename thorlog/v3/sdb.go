package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type SdbEntry struct {
	jsonlog.ObjectHeader

	Entry KeyValueList `json:"entry" textlog:"entry"`
}

func (SdbEntry) observed() {}

const typeSdbEntry = "shim database entry"

func init() { AddLogObjectType(typeSdbEntry, &SdbEntry{}) }

func NewSdbEntry() *SdbEntry {
	return &SdbEntry{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeSdbEntry,
		},
	}
}
