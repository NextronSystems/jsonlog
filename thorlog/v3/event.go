package thorlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/NextronSystems/jsonlog"
	"github.com/NextronSystems/jsonlog/jsonpointer"
	"github.com/NextronSystems/jsonlog/thorlog/common"
	"golang.org/x/exp/slices"
)

type Finding struct {
	jsonlog.ObjectHeader
	Meta         LogEventMetadata `json:"meta" textlog:",expand"`
	Text         string           `json:"message" textlog:"message"`
	Subject      ReportableObject `json:"subject" textlog:",expand"`
	Score        int64            `json:"score" textlog:"score"`
	Reasons      []Reason         `json:"reasons" textlog:",expand"`
	ReasonCount  int              `json:"reason_count,omitempty" textlog:"reasons_count,omitempty"`
	EventContext Context          `json:"context" textlog:",expand" jsonschema:"nullable"`
	Issues       []Issue          `json:"issues,omitempty" textlog:"-"`
	LogVersion   common.Version   `json:"log_version"`
}

type ReportableObject interface {
	reportable()
	jsonlog.Object
}

func (f *Finding) Message() string {
	return f.Text
}

func (f *Finding) Version() common.Version {
	return f.LogVersion
}

func (f *Finding) Metadata() *LogEventMetadata {
	return &f.Meta
}

func (f *Finding) UnmarshalJSON(data []byte) error {
	type plainFinding Finding
	var rawFinding struct {
		plainFinding                // Embed without unmarshal method to avoid infinite recursion
		Subject      EmbeddedObject `json:"subject"` // EmbeddedObject is used to allow unmarshalling of the subject as a ReportableObject
	}
	if err := json.Unmarshal(data, &rawFinding); err != nil {
		return err
	}
	subject, ok := rawFinding.Subject.Object.(ReportableObject)
	if !ok {
		return fmt.Errorf("subject must implement the reportable interface")
	}
	*f = Finding(rawFinding.plainFinding) // Copy the fields from rawFinding to f
	f.Subject = subject

	// Resolve all references
	// When the event is unmarshalled, the references are not resolved yet and only contain the JSON pointers.
	// Resolve them to the actual values to be able to use them in the text log.
	for i := range f.Reasons {
		for j := range f.Reasons[i].StringMatches {
			if f.Reasons[i].StringMatches[j].Field == nil {
				continue
			}
			target, err := jsonpointer.Resolve(f.Subject, f.Reasons[i].StringMatches[j].Field.ToJsonPointer())
			if err != nil {
				return err
			}
			f.Reasons[i].StringMatches[j].Field = jsonlog.NewReference(f.Subject, target)
		}
	}
	for i := range f.Issues {
		if f.Issues[i].Affected == nil {
			continue
		}
		target, err := jsonpointer.Resolve(f, f.Issues[i].Affected.ToJsonPointer())
		if err != nil {
			return err
		}
		f.Issues[i].Affected = jsonlog.NewReference(f, target)
	}
	return nil
}

var _ common.Event = (*Finding)(nil)

type Context []ContextObject

type ContextObject struct {
	Object       ReportableObject `json:"object" textlog:",expand"`
	RelationType string           `json:"relation_type"` // RelationType is used to specify the type of relation, e.g. "derives from" or "related to"
	RelationName string           `json:"relation_name"` // RelationName is used to specify the name of the relation, e.g. "parent". It is optional.
	Unique       bool             `json:"unique"`        // Unique indicates whether the relation is unique, i.e. there can only be one object with this relation type / name in the context.
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
	reportableObject, isReportable := rawContextObject.Object.Object.(ReportableObject)
	if !isReportable {
		return fmt.Errorf("object of type %q must implement the reportable interface", rawContextObject.Object.Object.EmbeddedHeader().Type)
	}
	*c = ContextObject(rawContextObject.plainContextObject) // Copy the fields from rawContextObject to c
	c.Object = reportableObject
	return nil
}

const omitInContext = "omitincontext"

func (c Context) MarshalTextLog(t jsonlog.TextlogFormatter) jsonlog.TextlogEntry {
	var elementsByRelation [][]ContextObject
	for _, element := range c {
		var groupExists bool
		for i := range elementsByRelation {
			if elementsByRelation[i][0].RelationName == element.RelationName {
				elementsByRelation[i] = append(elementsByRelation[i], element)
				groupExists = true
				break
			}
		}
		if !groupExists {
			elementsByRelation = append(elementsByRelation, []ContextObject{element})
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
		for g, element := range group {
			marshaledElement := t.Format(element)
			for i := range marshaledElement {
				marshaledElement[i].Key = jsonlog.ConcatTextLabels(strings.ToUpper(element.RelationName), marshaledElement[i].Key)
				if !element.Unique {
					marshaledElement[i].Key = jsonlog.ConcatTextLabels(marshaledElement[i].Key, strconv.Itoa(g+1))
				}
			}
			result = append(result, marshaledElement...)
		}
	}
	return result
}

const typeFinding = "THOR finding"

func init() { AddLogObjectType(typeFinding, &Finding{}) }

func NewFinding(subject ReportableObject, message string) *Finding {
	return &Finding{
		ObjectHeader: LogObjectHeader{
			Type: typeFinding,
		},
		Text:       message,
		Subject:    subject,
		LogVersion: currentVersion,
	}
}

type Message struct {
	jsonlog.ObjectHeader
	Meta       LogEventMetadata `json:"meta" textlog:",expand"`
	Text       string           `json:"message" textlog:"message"`
	Fields     MessageFields    `json:"fields" textlog:",expand" jsonschema:"nullable"`
	LogVersion common.Version   `json:"log_version"`
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
