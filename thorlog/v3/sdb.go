package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type SdbEntry struct {
	jsonlog.ObjectHeader

	File string `json:"file,omitempty" textlog:"file,omitempty"`

	Entry KeyValueList `json:"entry" textlog:"entry"`
}

const typeSdbEntry = "Shim Database entry"

func init() { AddLogObjectType(typeSdbEntry, &SdbEntry{}) }

func NewSdbEntry() *SdbEntry {
	return &SdbEntry{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeSdbEntry,
		},
	}
}

func (a SdbEntry) Truncate(matches []jsonlog.FieldMatch, truncateLimit int, stringContext int) jsonlog.Object {
	a.Entry = a.Entry.Truncate(matches, truncateLimit, stringContext)
	return &a
}
