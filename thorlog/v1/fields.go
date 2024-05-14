package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/NextronSystems/jsonlog"
)

type Field struct {
	Key   string
	Value string
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
	decoder := json.NewDecoder(bytes.NewReader(data))
	startToken, err := decoder.Token()
	if err != nil {
		return err
	}

	if startToken != json.Delim('{') {
		return &json.UnmarshalTypeError{
			Value:  string(data),
			Type:   reflect.TypeOf(o).Elem(),
			Offset: 0,
		}
	}

	var details Fields
	for decoder.More() {
		keyToken, err := decoder.Token()
		if err != nil {
			return err
		}
		key, isString := keyToken.(string)
		if !isString {
			return fmt.Errorf("key %v is not a string", keyToken)
		}
		var value any
		if err := decoder.Decode(&value); err != nil {
			return err
		}
		details = append(details, Field{
			Key:   key,
			Value: fmt.Sprint(value),
		})
	}
	*o = details
	return nil
}

func (o Fields) MarshalTextLog(_ jsonlog.TextlogFormatter) jsonlog.TextlogEntry {
	var result jsonlog.TextlogEntry
	for _, kv := range o {
		result = append(result, jsonlog.TextlogValuePair{
			Key:   kv.Key,
			Value: kv.Value,
		})
	}
	return result
}
