package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type SdbEntry struct {
	jsonlog.ObjectHeader

	File string `json:"file,omitempty" textlog:"file,omitempty"`

	Entry KeyValueList `json:"entry" textlog:"entry"`
}

func (SdbEntry) reportable() {}

const typeSdbEntry = "shim database entry"

func init() { AddLogObjectType(typeSdbEntry, &SdbEntry{}) }

func NewSdbEntry() *SdbEntry {
	return &SdbEntry{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeSdbEntry,
		},
	}
}
