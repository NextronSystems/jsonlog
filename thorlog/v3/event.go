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
)

type Finding struct {
	jsonlog.ObjectHeader
	Meta         LogEventMetadata `json:"meta" textlog:",expand"`
	Subject      jsonlog.Object   `json:"subject" textlog:",expand"`
	Score        int64            `json:"score" textlog:"score"`
	Reasons      []Reason         `json:"reasons" textlog:",expand"`
	ReasonCount  int              `json:"-" textlog:"reasons_count"`
	EventContext Context          `json:"context" textlog:",expand" jsonschema:"nullable"`
	LogVersion   common.Version   `json:"log_version"`
}

func (f *Finding) Message() string {
	return f.Summary
}

func (f *Finding) Version() common.Version {
	return f.LogVersion
}

func (f *Finding) Metadata() *LogEventMetadata {
	return &f.Meta
}

func (f *Finding) UnmarshalJSON(data []byte) error {
	var rawFinding struct {
		jsonlog.ObjectHeader
		Meta         LogEventMetadata `json:"meta"`
		Subject      EmbeddedObject   `json:"subject"`
		Score        int64            `json:"score"`
		Reasons      []Reason         `json:"reasons"`
		EventContext Context          `json:"context"`
		LogVersion   common.Version   `json:"log_version"`
	}
	if err := json.Unmarshal(data, &rawFinding); err != nil {
		return err
	}
	f.ObjectHeader = rawFinding.ObjectHeader
	f.Meta = rawFinding.Meta
	f.Subject = rawFinding.Subject.Object
	f.Score = rawFinding.Score
	f.Reasons = rawFinding.Reasons
	f.EventContext = rawFinding.EventContext
	f.LogVersion = rawFinding.LogVersion

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
	return nil
}

var _ common.Event = (*Finding)(nil)

type Context []ContextObject

type ContextObject struct {
	Object   jsonlog.Object `json:"object" textlog:",expand"`
	Relation string         `json:"relation"`
	Unique   bool           `json:"unique"`
}

func (c *ContextObject) UnmarshalJSON(data []byte) error {
	var rawContextObject struct {
		Object   EmbeddedObject `json:"object"`
		Relation string         `json:"relation"`
		Unique   bool           `json:"unique"`
	}
	if err := json.Unmarshal(data, &rawContextObject); err != nil {
		return err
	}
	c.Object = rawContextObject.Object.Object
	c.Relation = rawContextObject.Relation
	c.Unique = rawContextObject.Unique
	return nil
}

func (c Context) MarshalTextLog(t jsonlog.TextlogFormatter) jsonlog.TextlogEntry {
	var elementsByRelation [][]ContextObject
	for _, element := range c {
		var groupExists bool
		for i := range elementsByRelation {
			if elementsByRelation[i][0].Relation == element.Relation {
				elementsByRelation[i] = append(elementsByRelation[i], element)
				groupExists = true
				break
			}
		}
		if !groupExists {
			elementsByRelation = append(elementsByRelation, []ContextObject{element})
		}
	}

	var result jsonlog.TextlogEntry
	for _, group := range elementsByRelation {
		for g, element := range group {
			marshaledElement := t.Format(element)
			for i := range marshaledElement {
				marshaledElement[i].Key = jsonlog.ConcatTextLabels(strings.ToUpper(element.Relation), marshaledElement[i].Key)
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

func NewFinding(subject jsonlog.Object, message string) *Finding {
	return &Finding{
		ObjectHeader: LogObjectHeader{
			Type:    typeFinding,
			Summary: message,
		},
		Subject:    subject,
		LogVersion: currentVersion,
	}
}

type Message struct {
	jsonlog.ObjectHeader
	Meta       LogEventMetadata `json:"meta" textlog:",expand"`
	Fields     MessageFields    `json:"fields" textlog:",expand" jsonschema:"nullable"`
	LogVersion common.Version   `json:"log_version"`
}

func (m *Message) Message() string {
	return m.Summary
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
			Type:    typeMessage,
			Summary: message,
		},
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
