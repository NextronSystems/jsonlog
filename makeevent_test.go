package jsonlog

import (
	"testing"

	"github.com/NextronSystems/jsonlog/jsonpointer"
	"github.com/stretchr/testify/assert"
)

type testObject struct {
	ObjectHeader

	Substruct struct {
		SubField1 string `json:"subfield1" textlog:"subfield1"`
	} `json:"substruct" textlog:"substruct,expand"`

	unexported string

	AnonymousSubstruct

	Nested NestedSubstruct `json:"nested" textlog:"nested,expand"`

	Unexpanded UnexpandedSubstruct `json:"unexpanded" textlog:"unexpanded"`

	Subfield5 string `json:"subfield5" textlog:"subfield5"`

	Valuer TestEventValuer `json:"valuer" textlog:"valuer"`

	SubObject *SubObject `json:"subobject" textlog:"subobject,expand"`
}

type AnonymousSubstruct struct {
	SubField2 string `json:"subfield2" textlog:"subfield2"`
}

type NestedSubstruct struct {
	Substruct struct {
		SubField3 string `json:"subfield3" textlog:"subfield3"`
	} `json:"substruct" textlog:",expand"`
}

type UnexpandedSubstruct struct {
	SubField4 string `json:"subfield4" textlog:"subfield4"`
}

func (u UnexpandedSubstruct) String() string {
	return u.SubField4
}

type TestEventValuer struct {
	Subfield6 string `json:"subfield6"`
	Subfield7 string `json:"subfield7"`
	Ignore    string
}

func (t *TestEventValuer) Values() []EventValue {
	return []EventValue{
		{&t.Subfield6, t.Subfield6, "subfield6", jsonpointer.New("subfield6")},
		{&t.Subfield7, t.Subfield7, "subfield7", jsonpointer.New("subfield7")},
	}
}

func (t *TestEventValuer) RelativeTextPointer(pointee any) (string, bool) {
	if pointee == &t.Subfield6 {
		return "subfield6", true
	}
	if pointee == &t.Subfield7 {
		return "subfield7", true
	}
	return "", false
}

func (t *TestEventValuer) RelativeJsonPointer(pointee any) jsonpointer.Pointer {
	if pointee == &t.Subfield6 {
		return jsonpointer.New("subfield6")
	}
	if pointee == &t.Subfield7 {
		return jsonpointer.New("subfield7")
	}
	return nil
}

type SubObject struct {
	ObjectHeader
	Subfield8 string `json:"subfield8" textlog:"subfield8"`
}

func TestCreateEvent(t *testing.T) {
	var test testObject
	test.Substruct.SubField1 = "subfield1"
	test.SubField2 = "subfield2"
	test.Nested.Substruct.SubField3 = "subfield3"
	test.Unexpanded.SubField4 = "subfield4"
	test.Subfield5 = "subfield5"
	test.Valuer.Subfield6 = "subfield6"
	test.Valuer.Subfield7 = "subfield7"
	test.SubObject = &SubObject{Subfield8: "subfield8"}

	event := CreateEvent(&test)

	var expected = []ReferencedField{
		{Reference{Base: &test, PointedField: &test.Substruct.SubField1}, "subfield1"},
		{Reference{Base: &test, PointedField: &test.SubField2}, "subfield2"},
		{Reference{Base: &test, PointedField: &test.Nested.Substruct.SubField3}, "subfield3"},
		{Reference{Base: &test, PointedField: &test.Unexpanded}, "subfield4"},
		{Reference{Base: &test, PointedField: &test.Subfield5}, "subfield5"},
		{Reference{Base: &test, PointedField: &test.Valuer.Subfield6}, "subfield6"},
		{Reference{Base: &test, PointedField: &test.Valuer.Subfield7}, "subfield7"},
		{Reference{Base: &test, PointedField: &test.SubObject.Subfield8}, "subfield8"},
	}
	for i := range expected {
		expected[i].Reference.jsonPointer = expected[i].Reference.ToJsonPointer()
		expected[i].Reference.textLabel = expected[i].Reference.ToTextLabel()
	}

	assert.Len(t, event, len(expected))
	for i, field := range event {
		t.Logf("%v: %v", field.Reference.ToJsonPointer(), field.Value)
		assert.Equal(t, expected[i].Reference.Base, event[i].Reference.Base, "Base %d", i)
		assert.Equal(t, expected[i].Reference.PointedField, event[i].Reference.PointedField, "PointedField %d", i)
		assert.Equal(t, expected[i].Reference.textLabel, event[i].Reference.textLabel, "TextLabel %d", i)
		assert.Equal(t, expected[i].Reference.jsonPointer, event[i].Reference.jsonPointer, "JsonPointer %d", i)
		assert.Equal(t, expected[i].Value, event[i].Value, "Value", i)
	}
}
