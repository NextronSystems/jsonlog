package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
)

type AuditLogEntry struct {
	jsonlog.ObjectHeader

	// Time when the audit log entry was recorded, taken from the record's
	// audit(<seconds>.<milliseconds>:<serial>) header.
	Time  time.Time    `json:"time" textlog:"time"`
	Entry KeyValueList `json:"entry" textlog:"entry"`
}

const TypeAuditLogEntry = "audit log entry"

func init() { AddLogObjectType(TypeAuditLogEntry, &AuditLogEntry{}) }

func NewAuditLogEntry() *AuditLogEntry {
	return &AuditLogEntry{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: TypeAuditLogEntry,
		},
	}
}

func (AuditLogEntry) observed() {}
