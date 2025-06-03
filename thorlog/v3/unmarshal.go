package thorlog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/NextronSystems/jsonlog"
)

var ErrNoLogObject = errors.New("JSON does not contain a log object")

// UnknownObject is a log object that is not known to the parser.
type UnknownObject struct {
	jsonlog.ObjectHeader
	Data map[string]any
}

func (UnknownObject) reportable() {}

func (u *UnknownObject) UnmarshalJSON(data []byte) error {
	var details map[string]any
	err := json.Unmarshal(data, &details)
	if err != nil {
		return err
	}
	u.Data = details
	err = json.Unmarshal(data, &u.ObjectHeader)
	if err != nil {
		return err
	}
	return nil
}

func (u UnknownObject) MarshalTextLog(f jsonlog.TextlogFormatter) (jsonlog.TextlogEntry, error) {
	return marshalUnknownJsonObject(f, "", u.Data)
}

func marshalUnknownJsonObject(f jsonlog.TextlogFormatter, k string, v any) (jsonlog.TextlogEntry, error) {
	var fields jsonlog.TextlogEntry
	switch value := v.(type) {
	case map[string]any:
		subfields, err := UnknownObject{Data: value}.MarshalTextLog(f)
		if err != nil {
			return nil, err
		}
		for _, subfield := range subfields {
			fields = append(fields, jsonlog.TextlogValuePair{
				Key:   jsonlog.ConcatTextLabels(k, subfield.Key),
				Value: subfield.Value,
			})
		}
	case []any:
		for i, subvalue := range value {
			subfields, err := marshalUnknownJsonObject(f, k, subvalue)
			if err != nil {
				return nil, err
			}
			for _, subfield := range subfields {
				fields = append(fields, jsonlog.TextlogValuePair{
					Key:   jsonlog.ConcatTextLabels(subfield.Key, strconv.Itoa(i+1)),
					Value: subfield.Value,
				})
			}
		}
	default:
		var formattedValue string
		if f.FormatValue != nil {
			formattedValue = f.FormatValue(value, nil)
		} else {
			formattedValue = fmt.Sprint(value)
		}
		fields = append(fields, jsonlog.TextlogValuePair{
			Key:   k,
			Value: formattedValue,
		})
	}
	return fields, nil
}

// EmbeddedObject is a utility type for unmarshalling THOR log objects from JSON.
type EmbeddedObject struct {
	jsonlog.Object
}

func (e *EmbeddedObject) UnmarshalJSON(data []byte) error {
	var details map[string]any
	err := json.Unmarshal(data, &details)
	if err != nil {
		return err
	}
	objectType, exists := details["type"]
	if !exists {
		return ErrNoLogObject
	}
	objectTypeString, isString := objectType.(string)
	if !isString {
		return ErrNoLogObject
	}

	objectBlank := LogObjectTypes[objectTypeString]
	if objectBlank == nil {
		e.Object = &UnknownObject{
			Data:         details,
			ObjectHeader: jsonlog.ObjectHeader{Type: objectTypeString},
		}
		return nil
	}
	object := reflect.New(reflect.TypeOf(objectBlank).Elem()).Interface().(jsonlog.Object)

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(object)
	if err != nil {
		return err
	}
	e.Object = object
	return nil
}
