package thorlog

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/NextronSystems/jsonlog"
	"github.com/NextronSystems/jsonlog/thorlog/common"
)

type AuditEntry struct {
	jsonlog.ObjectHeader

	Id         string               `json:"id"`
	Subject    ObservedObject       `json:"subject"`
	Timestamps map[string]time.Time `json:"timestamps"`
	Reasons    []Reason             `json:"reasons" jsonschema:"nullable"`
	References []AuditReference     `json:"references" jsonschema:"nullable"`
	LogVersion common.Version       `json:"log_version"`
}

type AuditReference struct {
	TargetId     string `json:"target_id"`
	RelationName string `json:"relation_name"`
	RelationType string `json:"relation_type"`
}

const TypeAuditEntry = "THOR audit trail"

func init() { AddLogObjectType(TypeAuditEntry, &AuditEntry{}) }

func NewAuditEntry(id string, subject ObservedObject) *AuditEntry {
	return &AuditEntry{
		ObjectHeader: LogObjectHeader{
			Type: TypeAuditEntry,
		},
		Id:         id,
		Subject:    subject,
		LogVersion: currentVersion,
	}
}

func (a *AuditEntry) UnmarshalJSON(data []byte) error {
	type plainAuditEntry AuditEntry
	var rawAuditEntry struct {
		plainAuditEntry
		Subject EmbeddedObject `json:"subject"`
	}
	if err := json.Unmarshal(data, &rawAuditEntry); err != nil {
		return err
	}
	*a = AuditEntry(rawAuditEntry.plainAuditEntry)
	subject, isObserved := rawAuditEntry.Subject.Object.(ObservedObject)
	if !isObserved {
		return fmt.Errorf("subject of type %T must implement the ObservedObject interface",
			rawAuditEntry.Subject.Object)
	}
	a.Subject = subject
	return nil
}
