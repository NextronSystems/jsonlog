package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
)

type JournaldEntry struct {
	jsonlog.ObjectHeader

	Time    time.Time    `json:"time" textlog:"time"`
	Details KeyValueList `json:"details" textlog:"entry"`
}

const TypeJournaldEntry = "journal log entry"

func init() { AddLogObjectType(TypeJournaldEntry, &JournaldEntry{}) }

func NewJournaldEntry() *JournaldEntry {
	return &JournaldEntry{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: TypeJournaldEntry,
		},
	}
}

func (a JournaldEntry) Truncate(matches []jsonlog.FieldMatch, truncateLimit int, stringContext int) jsonlog.Object {
	a.Details = a.Details.Truncate(matches, truncateLimit, stringContext)
	return &a
}

func (JournaldEntry) reportable() {}
