package jsonlog

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TextlogTestObject is a test object to test textlog-specific details, mainly
// that is the TextlogFormatter's behavior.
type TextlogTestObject struct {
	ObjectHeader
	Element1   string `json:"element1" textlog:"element1"`
	_          string `json:"-" textlog:"ignored"`
	Element2   string `json:"element2" textlog:"element2"`
	unexported string
	Substruct  struct {
		Subelement1 string `json:"subelement1" textlog:"subelement1"`
	} `json:"substruct" textlog:"substruct,expand"`
	SimpleSubstruct
	Time time.Time `json:"time" textlog:"time"` // Time, despite being a struct, should be handled as a primitive type for the textlog since it does not have the "prefix" modifier set

	Omitempty string `json:"omitempty,omitempty" textlog:"omitempty,omitempty"` // Should be omitted in both JSON and textlog
}

type SimpleSubstruct struct {
	Subelement2 string `json:"subelement2" textlog:"subelement2"`
}

func TestToDetails(t *testing.T) {
	var test = TextlogTestObject{
		ObjectHeader: ObjectHeader{
			Type: "textlogtestobject",
		},
		Element1:   "element1",
		Element2:   "element2",
		unexported: "unexported",
		Substruct: struct {
			Subelement1 string `json:"subelement1" textlog:"subelement1"`
		}{
			Subelement1: "subelement1",
		},
		SimpleSubstruct: SimpleSubstruct{
			Subelement2: "subelement2",
		},
		Omitempty: "",
	}
	formatter := TextlogFormatter{
		FormatValue: func(data any, modifiers []string) string {
			if t, isTime := data.(time.Time); isTime {
				return t.Format(time.RFC3339)
			}
			return fmt.Sprint(data)
		},
	}
	details := formatter.Format(test)
	t.Log(details)
	assert.Equal(t, TextlogEntry{
		{"ELEMENT1", "element1"},
		{"ELEMENT2", "element2"},
		{"SUBSTRUCT_SUBELEMENT1", "subelement1"},
		{"SUBELEMENT2", "subelement2"},
		{"TIME", "0001-01-01T00:00:00Z"},
	}, details)
}

func TestTextlogFormatting(t *testing.T) {
	test := makeTestObject()

	formatter := TextlogFormatter{
		FormatValue: func(data any, modifiers []string) string {
			return fmt.Sprint(data)
		},
	}
	details := formatter.Format(test)
	t.Log(details)
	assert.Equal(t, TextlogEntry{
		{"SUBSTRUCT_SUBFIELD1", "subfield1"},
		{"SUBFIELD2", "subfield2"},
		{"NESTED_SUBFIELD3", "subfield3"},
		{"UNEXPANDED", "subfield4"},
		{"SUBFIELD5", "subfield5"},
		{"RESOLVER", "subfield6, subfield7"},
		{"SUBOBJECT_SUBFIELD8", "subfield8"},
		{"SLICE_SUBFIELD9_1", "slice"},
		{"MAP_KEY", "map"},
	}, details)
}

func TestTextlogOmit(t *testing.T) {
	type omitTestObject struct {
		Kept    string `json:"kept" textlog:"kept"`
		Omitted string `json:"omitted" textlog:"omitted" context:"omit"`
		Nested  struct {
			Omitted string `json:"omitted" textlog:"omitted" context:"omit"`
			Kept    int    `json:"kept" textlog:"kept"`
		} `json:"nested" textlog:"nested,expand"`
	}
	var test omitTestObject
	test.Kept = "kept"
	test.Omitted = "omitted"
	test.Nested.Omitted = "omitted"
	test.Nested.Kept = 42

	formatter := TextlogFormatter{
		FormatValue: func(data any, modifiers []string) string {
			return "formatted:" + fmt.Sprint(data)
		},
		Omit: func(field reflect.StructField, value any) bool {
			return field.Tag.Get("context") == "omit"
		},
	}
	assert.Equal(t, TextlogEntry{
		{"KEPT", "formatted:kept"},
		{"NESTED_KEPT", "formatted:42"},
	}, formatter.Format(test))
}
