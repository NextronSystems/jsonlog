package thorlog

import (
	"encoding/json"
	"testing"

	"github.com/NextronSystems/jsonlog"
	"github.com/NextronSystems/jsonlog/jsonpointer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyValueList_MarshalJSON(t *testing.T) {
	var kvList = KeyValueList{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
	}
	data, err := json.Marshal(kvList)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `{"key1":"value1","key2":"value2"}` {
		t.Errorf("unexpected JSON: %s", data)
	}
	var unmarshaled KeyValueList
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("unexpected error during unmarshal: %v", err)
	}
	if unmarshaled[0].Key != "key1" || unmarshaled[0].Value != "value1" || unmarshaled[1].Key != "key2" || unmarshaled[1].Value != "value2" {
		t.Errorf("unexpected unmarshaled data: %+v", unmarshaled)
	}
}

func TestKeyValueList_JsonPointers(t *testing.T) {
	var kvList = KeyValueList{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
	}
	reference := jsonlog.Reference{Base: &kvList, PointedField: &kvList[1].Value}
	pointer := reference.ToJsonPointer()
	assert.Equal(t, "/key2", pointer.String())

	reverse, err := jsonpointer.Resolve(&kvList, pointer)
	require.NoError(t, err)
	assert.Equal(t, &kvList[1].Value, reverse)
}
