package thorlog

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	"github.com/NextronSystems/jsonlog/jsonpointer"
)

type KeyValue struct {
	Key   string
	Value string
}

type KeyValueList struct {
	KvList []KeyValue
}

func (d KeyValueList) MarshalJSON() ([]byte, error) {
	var builder strings.Builder
	builder.WriteString("{")
	for i, kv := range d.KvList {
		if err := json.NewEncoder(&builder).Encode(kv.Key); err != nil {
			return nil, err
		}
		builder.WriteString(": ")
		if err := json.NewEncoder(&builder).Encode(kv.Value); err != nil {
			return nil, err
		}
		if i < len(d.KvList)-1 {
			builder.WriteString(", ")
		}
	}
	builder.WriteString("}")
	return []byte(builder.String()), nil
}

func (d *KeyValueList) UnmarshalJSON(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	if delim, isDelim := token.(json.Delim); !isDelim || delim != '{' {
		return errors.New("expected '{'")
	}
	var kvList []KeyValue
	for decoder.More() {
		var key string
		err = decoder.Decode(&key)
		if err != nil {
			return err
		}
		var value string
		err = decoder.Decode(&value)
		if err != nil {
			return err
		}
		kvList = append(kvList, KeyValue{Key: key, Value: value})
	}
	token, err = decoder.Token()
	if err != nil {
		return err
	}
	if delim, isDelim := token.(json.Delim); !isDelim || delim != '}' {
		return errors.New("expected '}'")
	}
	d.KvList = kvList
	return nil
}

func (d KeyValueList) RelativeJsonPointer(pointee any) jsonpointer.Pointer {
	stringPointer, isStringPointer := pointee.(*string)
	if !isStringPointer {
		return nil
	}
	for i := range d.KvList {
		if &d.KvList[i].Value == stringPointer {
			return jsonpointer.New(d.KvList[i].Key)
		}
	}
	return nil
}

func (d KeyValueList) RelativeTextPointer(pointee any) (string, bool) {
	stringPointer, isStringPointer := pointee.(*string)
	if !isStringPointer {
		return "", false
	}
	for i := range d.KvList {
		if &d.KvList[i].Value == stringPointer {
			return d.KvList[i].Key, true
		}
	}
	return "", false
}

func (d KeyValueList) Find(key string) *string {
	for i := range d.KvList {
		if d.KvList[i].Key == key {
			return &d.KvList[i].Value
		}
	}
	return nil
}

func (d KeyValueList) String() string {
	var dataBuilder strings.Builder
	for i, kv := range d.KvList {
		dataBuilder.WriteString(kv.Key)
		dataBuilder.WriteString(": ")
		dataBuilder.WriteString(kv.Value)
		if i < len(d.KvList)-1 {
			dataBuilder.WriteString("  ")
		}
	}
	return dataBuilder.String()
}

func (d KeyValueList) JSONSchemaAlias() any {
	return map[string]string{}
}
