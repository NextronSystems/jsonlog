package jsonlog

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestObject struct {
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
	var test = TestObject{
		ObjectHeader: ObjectHeader{
			Type: "testobject",
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
