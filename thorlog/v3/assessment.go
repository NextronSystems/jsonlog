package thorlog

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/NextronSystems/jsonlog"
	"github.com/NextronSystems/jsonlog/jsonpointer"
	"github.com/NextronSystems/jsonlog/thorlog/common"
	"golang.org/x/exp/slices"
)

// Assessment is a summary of a Subject's analysis by THOR.
// The assessed object is not necessarily suspicious; the
// severity can be seen in the Score, and beyond that the
// Reasons contain further information if this Subject is
// considered suspicious.
type Assessment struct {
	jsonlog.ObjectHeader
	Meta LogEventMetadata `json:"meta" textlog:",expand"`
	// Text is the message THOR printed for this assessment.
	// This is usually a summary based on this assessment's subject and level.
	Text string `json:"message" textlog:"message"`
	// Score is a metric that combines severity and certainty. The score is always in a range of 0 to 100;
	// 0 indicates that the assessment found no suspicious indicators, whereas 100 indicates very high
	// severity and certainty.
	Score int64 `json:"score" textlog:"score"`
	// Subject is the object assessed by THOR.
	Subject ObservedObject `json:"subject" textlog:",expand"`
	// Reasons describes the indicators that contributed to the score.
	// This list is not necessarily comprehensive; THOR may cut off all reasons after the first few.
	// If this is the case, an Issue with category IssueCategoryTruncated pointing to this field will be present.
	Reasons []Reason `json:"reasons" textlog:",expand"`
	// ReasonCount contains the total number of reasons (before any truncations).
	ReasonCount int `json:"reason_count,omitempty" textlog:"reasons_count,omitempty"`
	// EventContext contains other objects that may be relevant for an analyst and their relation to the
	// Subject.
	//
	// To give an example: if the Subject is a file in a ZIP archive,
	// the ZIP archive would be listed in the EventContext with a relation type of "derives from"
	// and a relation name of "parent", indicating that the Subject derives from this object,
	// which is its parent.
	EventContext Context `json:"context" textlog:",expand" jsonschema:"nullable"`
	// Issues lists any problems that THOR encountered when trying to create a JSON struct for this assessment.
	// This may include e.g. overly long fields that were truncated, fields that could not be rendered to JSON,
	// or similar problems.
	Issues []Issue `json:"issues,omitempty" textlog:"-"`
	// LogVersion describes the jsonlog version that this event was created with.
	LogVersion common.Version `json:"log_version"`
}

// ObservedObject can be any object type that THOR observes, e.g. File or Process.
type ObservedObject interface {
	observed()
	jsonlog.Object
}

func (a *Assessment) Message() string {
	return a.Text
}

func (a *Assessment) Version() common.Version {
	return a.LogVersion
}

func (a *Assessment) Metadata() *LogEventMetadata {
	return &a.Meta
}

func (a *Assessment) UnmarshalJSON(data []byte) error {
	type plainAssessment Assessment
	var rawAssessment struct {
		plainAssessment                // Embed without unmarshal method to avoid infinite recursion
		Subject         EmbeddedObject `json:"subject"` // EmbeddedObject is used to allow unmarshalling of the subject as a ObservedObject
	}
	if err := json.Unmarshal(data, &rawAssessment); err != nil {
		return err
	}
	subject, ok := rawAssessment.Subject.Object.(ObservedObject)
	if !ok {
		return fmt.Errorf("subject must implement the ObservedObject interface")
	}
	*a = Assessment(rawAssessment.plainAssessment) // Copy the fields from rawAssessment to a
	a.Subject = subject

	// Resolve all references
	// When the event is unmarshalled, the references are not resolved yet and only contain the JSON pointers.
	// Resolve them to the actual values to be able to use them in the text log.
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
	for i := range a.Issues {
		if a.Issues[i].Affected == nil {
			continue
		}
		target, err := jsonpointer.Resolve(a, a.Issues[i].Affected.ToJsonPointer())
		if err != nil {
			return err
		}
		a.Issues[i].Affected = jsonlog.NewReference(a, target)
	}
	return nil
}

var _ common.Event = (*Assessment)(nil)

type Context []ContextObject

// ContextObject describes a relation of an object to another.
type ContextObject struct {
	Object ObservedObject `json:"object" textlog:",expand"`
	// Relations describes how the object relates to the assessed subject.
	// There may be multiple relations, e.g. if the object is both the parent and the topmost ancestor of the subject.
	//
	// Relations should be ordered by relevance, i.e. the most important relation should be first.
	// Only the first (and most relevant) relation is used for text log formatting.
	Relations []Relation `json:"relations" textlog:",expand" jsonschema:"minItems=1"`
}

type Relation struct {
	Type   string `json:"relation_type"` // RelationType is used to specify the type of relation, e.g. "derives from" or "related to"
	Name   string `json:"relation_name"` // RelationName is used to specify the name of the relation, e.g. "parent". It is optional.
	Unique bool   `json:"unique"`        // Unique indicates whether the relation is unique, i.e. there can only be one object with this relation type / name in the context.
}

func (c *ContextObject) UnmarshalJSON(data []byte) error {
	type plainContextObject ContextObject
	var rawContextObject struct {
		Object EmbeddedObject `json:"object"`
		plainContextObject
	}
	if err := json.Unmarshal(data, &rawContextObject); err != nil {
		return err
	}
	reportableObject, isReportable := rawContextObject.Object.Object.(ObservedObject)
	if !isReportable {
		return fmt.Errorf("object of type %q must implement the ObservedObject interface", rawContextObject.Object.Object.EmbeddedHeader().Type)
	}
	*c = ContextObject(rawContextObject.plainContextObject) // Copy the fields from rawContextObject to c
	c.Object = reportableObject
	return nil
}

const omitInContext = "omitincontext"

func (c Context) MarshalTextLog(t jsonlog.TextlogFormatter) jsonlog.TextlogEntry {
	type objectsByRelation struct {
		Relation Relation
		Objects  []ContextObject
	}
	var elementsByRelation []objectsByRelation
	for _, element := range c {
		var groupExists bool
		if len(element.Relations) == 0 {
			continue
		}
		// only use the first relation for textlog conversion
		relation := element.Relations[0]
		for i := range elementsByRelation {
			if elementsByRelation[i].Relation == relation {
				elementsByRelation[i].Objects = append(elementsByRelation[i].Objects, element)
				groupExists = true
				break
			}
		}
		if !groupExists {
			elementsByRelation = append(elementsByRelation, objectsByRelation{Relation: relation, Objects: []ContextObject{element}})
		}
	}
	oldOmit := t.Omit
	t.Omit = func(modifiers []string, value any) bool {
		if slices.Contains(modifiers, omitInContext) {
			return true // Omit fields that are marked with "omitincontext"
		}
		if oldOmit != nil {
			return oldOmit(modifiers, value) // Call the original omit function if it exists
		}
		return false // Default behavior is to not omit any fields
	}

	var result jsonlog.TextlogEntry
	for _, group := range elementsByRelation {
		for g, element := range group.Objects {
			marshaledElement := t.Format(element)
			for i := range marshaledElement {
				marshaledElement[i].Key = jsonlog.ConcatTextLabels(strings.ToUpper(group.Relation.Name), marshaledElement[i].Key)
				if !group.Relation.Unique {
					marshaledElement[i].Key = jsonlog.ConcatTextLabels(marshaledElement[i].Key, strconv.Itoa(g+1))
				}
			}
			result = append(result, marshaledElement...)
		}
	}
	return result
}

const typeAssessment = "THOR assessment"

func init() { AddLogObjectType(typeAssessment, &Assessment{}) }

func NewAssessment(subject ObservedObject, message string) *Assessment {
	return &Assessment{
		ObjectHeader: LogObjectHeader{
			Type: typeAssessment,
		},
		Text:       message,
		Subject:    subject,
		LogVersion: CurrentVersion,
	}
}
