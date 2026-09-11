package audittrail

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/NextronSystems/jsonlog"
	"github.com/NextronSystems/jsonlog/jsonpointer"
	"github.com/NextronSystems/jsonlog/thorlog/common"
	"github.com/NextronSystems/jsonlog/thorlog/v3"
)

// AuditRecord describes a single object that THOR observed during a scan.
// Audit records are written regardless of whether THOR reported anything about the
// object. Audit records are written to a separate log rather than
// to THOR's regular output.
type AuditRecord struct {
	jsonlog.ObjectHeader

	// ID identifies this record within its scan.
	ID string `json:"id"`
	// Subject is the object that THOR observed.
	Subject thorlog.ObservedObject `json:"subject"`
	// Times contains the timestamps that are associated with the Subject.
	// This is a summary of all timestamps listed in the Subject.
	// If the Subject does not contain any timestamps, the time at which the Subject
	// was scanned is listed.
	Times map[string]time.Time `json:"timestamps"`
	// Reasons describes the indicators that THOR found for the Subject.
	Reasons []thorlog.Reason `json:"reasons" jsonschema:"nullable"`
	// References points to the other audit records that the Subject is related to.
	References []AuditReference `json:"references" jsonschema:"nullable"`
	// LogVersion describes the jsonlog version that this record was created with.
	LogVersion common.Version `json:"log_version"`
}

// AuditReference describes a relation of an audit record to another one.
type AuditReference struct {
	// TargetID is the Id of the referenced record.
	TargetID string `json:"target_id"`
	// RelationName is the name of the relation, e.g. "parent". It is optional.
	RelationName string `json:"relation_name"`
	// RelationType is the type of the relation, e.g. "derives from" or "related to".
	RelationType string `json:"relation_type"`
}

const TypeAuditRecord = "THOR audit record"

func init() { thorlog.AddLogObjectType(TypeAuditRecord, &AuditRecord{}) }

func NewAuditRecord(id string, subject thorlog.ObservedObject) *AuditRecord {
	return &AuditRecord{
		ObjectHeader: thorlog.LogObjectHeader{
			Type: TypeAuditRecord,
		},
		ID:         id,
		Subject:    subject,
		LogVersion: thorlog.CurrentVersion,
	}
}

func (a *AuditRecord) UnmarshalJSON(data []byte) error {
	type plainAuditRecord AuditRecord
	var rawAuditRecord struct {
		plainAuditRecord
		Subject thorlog.EmbeddedObject `json:"subject"`
	}
	if err := json.Unmarshal(data, &rawAuditRecord); err != nil {
		return err
	}
	*a = AuditRecord(rawAuditRecord.plainAuditRecord)
	subject, isObserved := rawAuditRecord.Subject.Object.(thorlog.ObservedObject)
	if !isObserved {
		return fmt.Errorf("subject of type %T must implement the ObservedObject interface",
			rawAuditRecord.Subject.Object)
	}
	a.Subject = subject

	// Resolve all references
	// When the audit record is unmarshalled, the references are not resolved yet and only contain
	// the JSON pointers. Resolve them to the actual values to be able to use them afterwards.
	for i := range a.Reasons {
		for j := range a.Reasons[i].StringMatches {
			if a.Reasons[i].StringMatches[j].Field == nil {
				continue
			}
			target, err := jsonpointer.Resolve(subject, a.Reasons[i].StringMatches[j].Field.ToJsonPointer())
			if err != nil {
				return err
			}
			a.Reasons[i].StringMatches[j].Field = jsonlog.NewReference(subject, target)
		}
	}
	return nil
}

func (a *AuditRecord) Version() common.Version {
	return a.LogVersion
}

func (a *AuditRecord) Timestamps() map[string]time.Time {
	return a.Times
}
