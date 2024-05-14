package thorlog

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	"github.com/NextronSystems/jsonlog"
	"github.com/NextronSystems/jsonlog/jsonpointer"
	"github.com/NextronSystems/jsonlog/thorlog/truncate"
	"golang.org/x/exp/slices"
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
		json.NewEncoder(&builder).Encode(kv.Key)
		builder.WriteString(": ")
		json.NewEncoder(&builder).Encode(kv.Value)
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

func (d KeyValueList) Values() []jsonlog.EventValue {
	var values []jsonlog.EventValue
	for i := range d.KvList {
		values = append(values, jsonlog.EventValue{
			FieldPointer: &d.KvList[i].Value,
			Value:        d.KvList[i].Value,
			TextLabel:    d.KvList[i].Key,
			JsonPointer:  jsonpointer.New(d.KvList[i].Key),
		})
	}
	return values
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

func (d KeyValueList) Truncate(matches []jsonlog.FieldMatch, truncateLimit int, stringContext int) KeyValueList {
	var fieldIndicesBySize = make([]int, len(d.KvList))
	for i := range d.KvList {
		fieldIndicesBySize[i] = i
	}
	slices.SortFunc(fieldIndicesBySize, func(a, b int) int {
		sizeA := len(d.KvList[a].Key) + len(d.KvList[a].Value)
		sizeB := len(d.KvList[b].Key) + len(d.KvList[b].Value)
		if sizeA < sizeB {
			return -1
		} else if sizeB > sizeA {
			return 1
		} else {
			return 0
		}
	})
	var availableSize = truncateLimit
	// Try to include as many full fields as possible, starting with the smallest ones
	for len(fieldIndicesBySize) > 0 {
		requiredSize := len(d.KvList[0].Key) + len(d.KvList[0].Value) + 4
		availableSizePerField := availableSize / len(fieldIndicesBySize)
		if requiredSize > availableSizePerField {
			break
		}
		availableSize -= requiredSize
		fieldIndicesBySize = fieldIndicesBySize[1:]
	}
	if len(fieldIndicesBySize) == 0 {
		// No truncation necessary
		return d
	}
	availableSizePerField := availableSize / len(fieldIndicesBySize)
	var relevantMatches = make([][]truncate.Match, len(d.KvList))
	for i := range d.KvList {
		for _, match := range matches {
			if match.FieldPointer == &d.KvList[i].Value {
				relevantMatches[i] = append(relevantMatches[i], match.Match)
			}
		}
	}
	var truncatedFields []KeyValue
	for i := range d.KvList {
		availableSizeForValue := availableSizePerField - len(d.KvList[i].Key) - 4
		truncatedFields = append(truncatedFields, KeyValue{
			Key:   d.KvList[i].Key,
			Value: truncate.SmartTruncate(d.KvList[i].Value, relevantMatches[i], availableSizeForValue, stringContext),
		})
	}
	return KeyValueList{KvList: truncatedFields}
}
