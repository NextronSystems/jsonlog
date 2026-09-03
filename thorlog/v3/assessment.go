package thorlog

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/NextronSystems/jsonlog"
	"github.com/NextronSystems/jsonlog/jsonpointer"
	"github.com/NextronSystems/jsonlog/thorlog/common"
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
	Reasons []Reason `json:"reasons" textlog:",expand" jsonschema:"nullable"`
	// ReasonCount contains the total number of reasons (before any truncations).
	ReasonCount int `json:"reason_count,omitempty" textlog:"reasons_count,omitempty"`
	// Ancestors contains information about objects that are the subject's ancestors.
	//
	// E.g. if the subject is a file in a nested ZIP, all the ZIPs (the topmost one and each nested one) are ancestors.
	// Ancestors does not necessarily include information about all ancestors, but it will always include information about at least the topmost ancestor and the parent.
	Ancestors Ancestors `json:"ancestors" textlog:",expand" jsonschema:"nullable"`
	// Derivatives contains information about objects that have been referenced in the subject.
	Derivatives []Derivative `json:"derivatives" textlog:"file,expand" jsonschema:"nullable"`
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
	for i := range a.Derivatives {
		if a.Derivatives[i].Via == nil {
			continue
		}
		target, err := jsonpointer.Resolve(a, a.Derivatives[i].Via.ToJsonPointer())
		if err != nil {
			return err
		}
		a.Derivatives[i].Via = jsonlog.NewReference(a, target)
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

type Ancestors []Ancestor

type Ancestor struct {
	// Per describes the action that caused the ancestor to create its child.
	Per string `json:"per"`
	// TopLevel is true if the ancestor does not have a parent.
	TopLevel bool `json:"top_level,omitempty"`
	// Distance is the number of links between the subject and the ancestor (parent = 1, grandparent = 2, ...)
	Distance int `json:"distance"`
	// Object describes the ancestor.
	Object ObservedObject `json:"object"`
}

func (a *Ancestor) UnmarshalJSON(data []byte) error {
	type plainAncestor Ancestor
	var rawAncestor struct {
		Object EmbeddedObject `json:"object"`
		plainAncestor
	}
	if err := json.Unmarshal(data, &rawAncestor); err != nil {
		return err
	}
	reportableObject, isReportable := rawAncestor.Object.Object.(ObservedObject)
	if !isReportable && rawAncestor.Object.Object != nil {
		return fmt.Errorf("object of type %q must implement the ObservedObject interface", rawAncestor.Object.Object.EmbeddedHeader().Type)
	}
	*a = Ancestor(rawAncestor.plainAncestor) // Copy the fields from rawAncestor to a
	a.Object = reportableObject
	return nil
}

func (a Ancestors) MarshalTextLog(t jsonlog.TextlogFormatter) jsonlog.TextlogEntry {
	t = withOmitInContext(t)
	var result jsonlog.TextlogEntry
	for _, ancestor := range a {
		var prefix string
		if ancestor.Distance == 1 {
			prefix = "parent"
		} else if ancestor.TopLevel {
			prefix = "origin"
		} else {
			continue
		}
		marshaledElement := t.Format(ancestor.Object)
		for i := range marshaledElement {
			marshaledElement[i].Key = jsonlog.ConcatTextLabels(strings.ToUpper(prefix), marshaledElement[i].Key)
		}
		result = append(result, marshaledElement...)
	}
	return result
}

type Derivative struct {
	// Via is a reference to the field that contains a link to Object.
	Via *jsonlog.Reference `json:"via"`
	// Object is the object that is derived from the assessment's Subject.
	Object ObservedObject `json:"object" textlog:",expand"`
}

// withOmitInContext returns a formatter that additionally omits all fields tagged with `context:"omit"`.
// This tag marks fields that are too verbose to be repeated for an object that is only
// logged as context (ancestor or derivative) of another object.
func withOmitInContext(t jsonlog.TextlogFormatter) jsonlog.TextlogFormatter {
	oldOmit := t.Omit
	t.Omit = func(field reflect.StructField, value any) bool {
		if field.Tag.Get("context") == "omit" {
			return true
		}
		if oldOmit != nil {
			return oldOmit(field, value)
		}
		return false
	}
	return t
}

func (d *Derivative) UnmarshalJSON(data []byte) error {
	type plainDerivative Derivative
	var unmarshalableDerivative struct {
		Object EmbeddedObject `json:"object"`
		plainDerivative
	}
	if err := json.Unmarshal(data, &unmarshalableDerivative); err != nil {
		return err
	}
	observedObject, isObservedObject := unmarshalableDerivative.Object.Object.(ObservedObject)
	if !isObservedObject && unmarshalableDerivative.Object.Object != nil {
		return fmt.Errorf("object of type %q must implement the ObservedObject interface", unmarshalableDerivative.Object.Object.EmbeddedHeader().Type)
	}
	*d = Derivative(unmarshalableDerivative.plainDerivative)
	d.Object = observedObject
	return nil
}

func (d Derivative) MarshalTextLog(t jsonlog.TextlogFormatter) jsonlog.TextlogEntry {
	t = withOmitInContext(t)
	type plainDerivative Derivative // Wrap this struct to not implement TextlogMarshaler
	return t.Format(plainDerivative(d))
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
