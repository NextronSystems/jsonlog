package v2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/NextronSystems/jsonlog"
)

type Field struct {
	Key   string
	Value any
}

type Fields []Field

func (o Fields) MarshalJSON() ([]byte, error) {
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

func (o *Fields) UnmarshalJSON(data []byte) error {
	value, err := unmarshalJsonValue(data)
	if err != nil {
		return err
	}
	if value == nil {
		return nil
	}
	details, isDetails := value.(Fields)
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
		var details Fields
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
			details = append(details, Field{
				Key:   key,
				Value: parsedValue,
			})
		}
		return details, nil
	} else {
		return nil, fmt.Errorf("invalid JSON token %v", startToken)
	}
}

func (m Fields) find(key string) (any, bool) {
	for _, kv := range m {
		if kv.Key == key {
			return kv.Value, true
		}
	}
	return "", false
}

func (m Fields) MarshalTextLog(t jsonlog.TextlogFormatter) jsonlog.TextlogEntry {
	var result jsonlog.TextlogEntry
	for _, kv := range m {
		switch v := kv.Value.(type) {
		case []any:
			var allComplex = true
			for _, value := range v {
				if _, ok := value.(Fields); !ok {
					allComplex = false
					break
				}
			}
			if allComplex {
				if kv.Key == "matched" { // Special case: Match strings
					var matchStrings []string
					for i, value := range v {
						subFields := value.(Fields)
						data, _ := subFields.find("data")
						context, hasContext := subFields.find("context")
						offset, hasOffset := subFields.find("offset")
						field, hasField := subFields.find("field")
						matchString := fmt.Sprintf("%q", data)
						if hasContext {
							matchString += fmt.Sprintf(" in %q", context)
						}
						if hasOffset {
							matchString += fmt.Sprintf(" at %v", offset)
						}
						if hasField {
							matchString += fmt.Sprintf(" in %s", field)
						}
						matchStrings = append(matchStrings, fmt.Sprintf("Str%d: %s", i+1, matchString))
					}
					result = append(result, jsonlog.TextlogValuePair{
						Key:   kv.Key,
						Value: strings.Join(matchStrings, " "),
					})
					continue
				}
				var prefix = kv.Key
				if nonPrefixedFields[kv.Key] { // Special case - do not prefix with key
					prefix = ""
				}
				for i, value := range v {
					subEntry := value.(Fields).MarshalTextLog(t)
					for _, expandedValue := range subEntry {
						key := modifyKeyForTextlog(kv.Key, expandedValue.Key)
						result = append(result, jsonlog.TextlogValuePair{
							Key:   fmt.Sprintf("%s_%d", jsonlog.ConcatTextLabels(prefix, key), i+1),
							Value: expandedValue.Value,
						})
					}
				}

			} else {
				var primitiveValues []string
				for _, value := range v {
					primitiveValues = append(primitiveValues, fmt.Sprint(value))
				}
				result = append(result, jsonlog.TextlogValuePair{
					Key:   kv.Key,
					Value: strings.Join(primitiveValues, ","),
				})
			}
		case Fields:
			subEntry := v.MarshalTextLog(t)
			var prefix = kv.Key
			if nonPrefixedFields[kv.Key] {
				prefix = ""
			}
			for _, expandedValue := range subEntry {
				key := modifyKeyForTextlog(kv.Key, expandedValue.Key)
				result = append(result, jsonlog.TextlogValuePair{
					Key:   jsonlog.ConcatTextLabels(prefix, key),
					Value: expandedValue.Value,
				})
			}
		default:
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
		}
	}
	return result
}

func modifyKeyForTextlog(parent string, key string) string {
	if key == "path" && (parent == "file" || parent == "process" || parent == "files" || parent == "app" || parent == "image" || parent == "archive") {
		key = "file"
	} else if key == "name" && parent == "reasons" {
		key = "reason"
	}
	return key
}

var nonPrefixedFields = map[string]bool{
	"reasons":   true,
	"file":      true,
	"process":   true,
	"files":     true,
	"signature": true,
}
