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
// object. Audit records are written to the audit trail and not
// to THOR's regular output.
type AuditRecord struct {
	jsonlog.ObjectHeader

	// ID identifies this record within its scan.
	ID string `json:"id"`
	// Object is the object that THOR observed.
	Object thorlog.ObservedObject `json:"object"`
	// Times contains the timestamps that are associated with the Object.
	// This is a list of all timestamps found anywhere in the Object.
	// It is guaranteed that Times is not empty:
	// if the object does not contain any timestamps, the time at which the Object
	// was scanned is listed.
	Times map[string]time.Time `json:"timestamps"`
	// Reasons describes the indicators that THOR found for the Object.
	Reasons []thorlog.Reason `json:"reasons" jsonschema:"nullable"`
	// References points to the other audit records that the Object is related to.
	References []AuditReference `json:"references" jsonschema:"nullable"`
	// LogVersion describes the jsonlog version that this record was created with.
	LogVersion common.Version `json:"log_version"`
}

// AuditReference describes a relation of an audit record to another one.
type AuditReference struct {
	// TargetID is the Id of the referenced record.
	TargetID string `json:"target_id"`
	// Relation is the type of the relation, either "child of" or "points to".
	Relation string `json:"relation"`
}

const TypeAuditRecord = "THOR audit record"

func init() { thorlog.AddLogObjectType(TypeAuditRecord, &AuditRecord{}) }

func NewAuditRecord(id string, object thorlog.ObservedObject) *AuditRecord {
	return &AuditRecord{
		ObjectHeader: thorlog.LogObjectHeader{
			Type: TypeAuditRecord,
		},
		ID:         id,
		Object:     object,
		LogVersion: thorlog.CurrentVersion,
	}
}

func (a *AuditRecord) UnmarshalJSON(data []byte) error {
	type plainAuditRecord AuditRecord
	var rawAuditRecord struct {
		plainAuditRecord
		Object thorlog.EmbeddedObject `json:"object"`
	}
	if err := json.Unmarshal(data, &rawAuditRecord); err != nil {
		return err
	}
	*a = AuditRecord(rawAuditRecord.plainAuditRecord)
	object, isObserved := rawAuditRecord.Object.Object.(thorlog.ObservedObject)
	if !isObserved {
		return fmt.Errorf("object of type %T must implement the ObservedObject interface",
			rawAuditRecord.Object.Object)
	}
	a.Object = object

	// Resolve all references
	// When the audit record is unmarshalled, the references are not resolved yet and only contain
	// the JSON pointers. Resolve them to the actual values to be able to use them afterwards.
	for i := range a.Reasons {
		for j := range a.Reasons[i].StringMatches {
			if a.Reasons[i].StringMatches[j].Field == nil {
				continue
			}
			target, err := jsonpointer.Resolve(object, a.Reasons[i].StringMatches[j].Field.ToJsonPointer())
			if err != nil {
				return err
			}
			a.Reasons[i].StringMatches[j].Field = jsonlog.NewReference(object, target)
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
