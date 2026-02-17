package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type AuditLogEntry struct {
	jsonlog.ObjectHeader

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
