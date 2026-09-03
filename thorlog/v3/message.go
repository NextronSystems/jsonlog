package thorlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/NextronSystems/jsonlog"
	"github.com/NextronSystems/jsonlog/thorlog/common"
)

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
		LogVersion: CurrentVersion,
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
