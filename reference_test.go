package jsonlog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReference_ToJsonPointer(t *testing.T) {
	var test testObject
	test.Substruct.SubField1 = "subfield1"
	test.SubField2 = "subfield2"
	test.Nested.Substruct.SubField3 = "subfield3"
	test.Unexpanded.SubField4 = "subfield4"
	test.Subfield5 = "subfield5"
	test.Valuer.Subfield6 = "subfield6"
	test.Valuer.Subfield7 = "subfield7"
	test.SubObject = &SubObject{Subfield8: "subfield8"}

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
		{&test.Valuer, "/valuer"},
		{&test.Valuer.Subfield6, "/valuer/subfield6"},
		{&test.Valuer.Subfield7, "/valuer/subfield7"},
		{&test.SubObject, "/subobject"},
		{&test.SubObject.Subfield8, "/subobject/subfield8"},
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
	test.Valuer.Subfield6 = "subfield6"
	test.Valuer.Subfield7 = "subfield7"
	test.SubObject = &SubObject{Subfield8: "subfield8"}

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
		{&test.Unexpanded.SubField4, "SUBFIELD4"},
		{&test.Subfield5, "SUBFIELD5"},
		{&test.Valuer, "VALUER"},
		{&test.Valuer.Subfield6, "subfield6"},
		{&test.Valuer.Subfield7, "subfield7"},
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
