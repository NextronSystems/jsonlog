package jsonlog

import (
	"strings"
	"testing"

	"github.com/NextronSystems/jsonlog/jsonpointer"
	"github.com/stretchr/testify/assert"
)

type testObject struct {
	ObjectHeader

	Recursive *Reference

	Substruct struct {
		SubField1 string `json:"subfield1" textlog:"subfield1"`
	} `json:"substruct" textlog:"substruct,expand"`

	// nolint:unused // not used, just used to check that unexported fields are not included in the event
	unexported string

	AnonymousSubstruct

	Nested NestedSubstruct `json:"nested" textlog:"nested,expand"`

	Unexpanded UnexpandedSubstruct `json:"unexpanded" textlog:"unexpanded"`

	Subfield5 string `json:"subfield5" textlog:"subfield5"`

	Resolver TestResolver `json:"resolver" textlog:"resolver"`

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

type TestResolver struct {
	Subfield6 string `json:"subfield6"`
	Subfield7 string `json:"subfield7"`
	Ignore    string
}

func (t TestResolver) String() string {
	return strings.Join([]string{t.Subfield6, t.Subfield7}, ", ")
}

var _ TextReferenceResolver = (*TestResolver)(nil)

func (t *TestResolver) RelativeTextPointer(pointee any) (string, bool, bool) {
	virtual := true // We do not implement TextlogMarshaler
	if pointee == &t.Subfield6 {
		return "subfield6", virtual, true
	}
	if pointee == &t.Subfield7 {
		return "subfield7", virtual, true
	}
	return "", false, false
}

var _ JsonReferenceResolver = (*TestResolver)(nil)

func (t *TestResolver) RelativeJsonPointer(pointee any) jsonpointer.Pointer {
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

func TestReference_ToJsonPointer(t *testing.T) {
	var test testObject
	test.Substruct.SubField1 = "subfield1"
	test.SubField2 = "subfield2"
	test.Nested.Substruct.SubField3 = "subfield3"
	test.Unexpanded.SubField4 = "subfield4"
	test.Subfield5 = "subfield5"
	test.Resolver.Subfield6 = "subfield6"
	test.Resolver.Subfield7 = "subfield7"
	test.SubObject = &SubObject{Subfield8: "subfield8"}
	test.Recursive = NewReference(&test, &test.Substruct)

	var tests = []struct {
		PointedField any
		want         string
	}{
		{&test.Substruct, "/substruct"},
		{&test.Substruct.SubField1, "/substruct/subfield1"},
		{&test.SubField2, "/subfield2"},
		{&test.Nested, "/nested"},
		{&test.Nested.Substruct, "/nested/substruct"},
		{&test.Nested.Substruct.SubField3, "/nested/substruct/subfield3"},
		{&test.Unexpanded, "/unexpanded"},
		{&test.Unexpanded.SubField4, "/unexpanded/subfield4"},
		{&test.Subfield5, "/subfield5"},
		{&test.Resolver, "/resolver"},
		{&test.Resolver.Subfield6, "/resolver/subfield6"},
		{&test.Resolver.Subfield7, "/resolver/subfield7"},
		{&test.SubObject, "/subobject"},
		{&test.SubObject.Subfield8, "/subobject/subfield8"},
		// test.Recursive not required here: if there is a flaw in the pointer
		// search logic, it will panic when encountering this cycle.
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			ref := NewReference(&test, tt.PointedField)
			assert.Equal(t, tt.want, ref.ToJsonPointer().String())
		})
	}
}

func TestReference_ToTextPointer(t *testing.T) {
	var test testObject
	test.Substruct.SubField1 = "subfield1"
	test.SubField2 = "subfield2"
	test.Nested.Substruct.SubField3 = "subfield3"
	test.Unexpanded.SubField4 = "subfield4"
	test.Subfield5 = "subfield5"
	test.Resolver.Subfield6 = "subfield6"
	test.Resolver.Subfield7 = "subfield7"
	test.SubObject = &SubObject{Subfield8: "subfield8"}
	test.Recursive = NewReference(&test, &test.Substruct)

	var tests = []struct {
		PointedField any
		want         string
	}{
		{&test.Substruct, "SUBSTRUCT"},
		{&test.Substruct.SubField1, "SUBSTRUCT_SUBFIELD1"},
		{&test.SubField2, "SUBFIELD2"},
		{&test.Nested, "NESTED"},
		{&test.Nested.Substruct, "NESTED"},
		{&test.Nested.Substruct.SubField3, "NESTED_SUBFIELD3"},
		{&test.Unexpanded, "UNEXPANDED"},
		{&test.Unexpanded.SubField4, ""},
		{&test.Subfield5, "SUBFIELD5"},
		{&test.Resolver, "RESOLVER"},
		{&test.Resolver.Subfield6, "subfield6"},
		{&test.Resolver.Subfield7, "subfield7"},
		{&test.SubObject, "SUBOBJECT"},
		{&test.SubObject.Subfield8, "SUBOBJECT_SUBFIELD8"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			ref := NewReference(&test, tt.PointedField)
			assert.Equal(t, tt.want, ref.ToTextLabel())
		})
	}
}
