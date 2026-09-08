package thorlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
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
	// Ancestors contains information about objects that are the subject's ancestors: E.g. if the subject is a file in a nested ZIP, each of the nested ZIPs is an ancestor.
	// Ancestors does not necessarily include information about all ancestors, but it will always include information about at least the topmost ancestor and the parent.
	Ancestors Ancestors `json:"ancestors" textlog:",expand"`
	// Derivatives contains information about objects that have been referenced in the subject.
	Derivatives []Derivative `json:"derivatives" textlog:"file,expand"`
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
			target, err := jsonpointer.Resolve(a.Subject, a.Reasons[i].StringMatches[j].Field.ToJsonPointer())
			if err != nil {
				return err
			}
			a.Reasons[i].StringMatches[j].Field = jsonlog.NewReference(a.Subject, target)
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

type Ancestors []Ancestor

type Ancestor struct {
	Per      string         `json:"per"`
	Object   ObservedObject `json:"object"`
	TopLevel bool           `json:"top_level,omitempty"`
	Distance int            `json:"distance"`
}

type Derivative struct {
	Object ObservedObject `json:"object" textlog:",expand"`
	// Via is a reference to the field that contains a link to Object.
	Via *jsonlog.Reference `json:"via"`
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
	if !isReportable {
		return fmt.Errorf("object of type %q must implement the ObservedObject interface", rawAncestor.Object.Object.EmbeddedHeader().Type)
	}
	*a = Ancestor(rawAncestor.plainAncestor) // Copy the fields from rawAncestor to a
	a.Object = reportableObject
	return nil
}

const omitInContext = "omitincontext"

func (a Ancestors) MarshalTextLog(t jsonlog.TextlogFormatter) jsonlog.TextlogEntry {
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
	for _, ancestor := range a {
		var prefix string
		if ancestor.Distance == 0 {
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

const typeAssessment = "THOR assessment"

func init() { AddLogObjectType(typeAssessment, &Assessment{}) }

func NewAssessment(subject ObservedObject, message string) *Assessment {
	return &Assessment{
		ObjectHeader: LogObjectHeader{
			Type: typeAssessment,
		},
		Text:       message,
		Subject:    subject,
		LogVersion: currentVersion,
	}
}

// Message describes a THOR message printed during the scan.
// Unlike Assessment, this does not describe an analysis' result,
// but rather something about the scan itself (e.g. how many IOCs were loaded).
type Message struct {
	jsonlog.ObjectHeader
	Meta LogEventMetadata `json:"meta" textlog:",expand"`
	// Text is the message that was logged.
	Text string `json:"message" textlog:"message"`
	// Fields contains additional structured fields that were logged. These
	// contain details about the Text displayed.
	Fields     MessageFields  `json:"fields" textlog:",expand" jsonschema:"nullable"`
	LogVersion common.Version `json:"log_version"`
}

func (m *Message) Message() string {
	return m.Text
}

func (m *Message) Version() common.Version {
	return m.LogVersion
}

func (m *Message) Metadata() *LogEventMetadata {
	return &m.Meta
}

var _ common.Event = (*Message)(nil)

const typeMessage = "THOR message"

func init() { AddLogObjectType(typeMessage, &Message{}) }

func NewMessage(meta LogEventMetadata, message string, kvs ...any) *Message {
	msg := &Message{
		ObjectHeader: LogObjectHeader{
			Type: typeMessage,
		},
		Text:       message,
		Meta:       meta,
		LogVersion: currentVersion,
	}
	if len(kvs)%2 != 0 {
		panic("uneven number of key-value pairs")
	}
	for i := 0; i < len(kvs); i += 2 {
		msg.Fields = append(msg.Fields, MessageField{
			Key:   kvs[i].(string),
			Value: kvs[i+1],
		})
	}
	return msg
}

type MessageField struct {
	Key   string
	Value any
}

type MessageFields []MessageField

func (o MessageFields) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("{")
	for i, kv := range o {
		if i != 0 {
			buf.WriteString(",")
		}
		key, err := json.Marshal(kv.Key)
		if err != nil {
			return nil, err
		}
		buf.Write(key)
		buf.WriteString(":")
		// marshal value
		val, err := json.Marshal(kv.Value)
		if err != nil {
			return nil, err
		}
		buf.Write(val)
	}

	buf.WriteString("}")
	return buf.Bytes(), nil
}

func (o *MessageFields) UnmarshalJSON(data []byte) error {
	value, err := unmarshalJsonValue(data)
	if err != nil {
		return err
	}
	if value == nil {
		return nil
	}
	details, isDetails := value.(MessageFields)
	if !isDetails {
		return &json.UnmarshalTypeError{
			Value:  fmt.Sprint(value),
			Type:   reflect.TypeOf(o).Elem(),
			Offset: 0,
		}
	}
	*o = details
	return nil
}

func (o MessageFields) JSONSchemaAlias() any {
	return map[string]any{}
}

func unmarshalJsonValue(data []byte) (any, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	startToken, err := decoder.Token()
	if err != nil {
		return nil, err
	}
	switch t := startToken.(type) {
	case bool, string, float64, json.Number, nil:
		return t, nil
	}
	if startToken == json.Delim('[') {
		var values []any
		for decoder.More() {
			var value json.RawMessage
			if err := decoder.Decode(&value); err != nil {
				return nil, err
			}
			parsedValue, err := unmarshalJsonValue(value)
			if err != nil {
				return nil, err
			}
			values = append(values, parsedValue)
		}
		return values, nil
	} else if startToken == json.Delim('{') {
		var details MessageFields
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return nil, err
			}
			key, isString := keyToken.(string)
			if !isString {
				return nil, fmt.Errorf("key %v is not a string", keyToken)
			}
			var value json.RawMessage
			if err := decoder.Decode(&value); err != nil {
				return nil, err
			}
			parsedValue, err := unmarshalJsonValue(value)
			if err != nil {
				return nil, err
			}
			details = append(details, MessageField{
				Key:   key,
				Value: parsedValue,
			})
		}
		return details, nil
	} else {
		return nil, fmt.Errorf("invalid JSON token %v", startToken)
	}
}

func (m MessageFields) MarshalTextLog(t jsonlog.TextlogFormatter) jsonlog.TextlogEntry {
	var result jsonlog.TextlogEntry
	for _, kv := range m {
		expandedValues := t.Format(kv.Value)
		if len(expandedValues) == 0 { // FIXME: Better distinguish between types that are expanded and those that aren't
			var formattedValue string
			if t.FormatValue != nil {
				formattedValue = t.FormatValue(kv.Value, nil)
			} else {
				formattedValue = fmt.Sprint(kv.Value)
			}
			result = append(result, jsonlog.TextlogValuePair{
				Key:   kv.Key,
				Value: formattedValue,
			})
		} else {
			result = append(result, expandedValues...)
		}
	}
	return result
}
