package thorlog

import (
	"strconv"
	"strings"

	"github.com/NextronSystems/jsonlog"
	"github.com/NextronSystems/jsonlog/jsonpointer"
)

// StringList is a list of strings with a comma-separated string representation.
//
// Note: StringList does not implement TextReferenceResolver. We do not want to return references to individual strings in the text log because this might lead to confusion; the label would point to the full string list but values like offset are relative to the individual string.
type StringList []string

func (s StringList) String() string {
	return strings.Join(s, ", ")
}

var _ jsonlog.JsonReferenceResolver = StringList{}

func (s StringList) RelativeJsonPointer(pointee any) jsonpointer.Pointer {
	stringPointer, isStringPointer := pointee.(*string)
	if !isStringPointer {
		return nil
	}
	for i := range s {
		if &s[i] == stringPointer {
			return jsonpointer.New(strconv.Itoa(i))
		}
	}
	return nil
}

// ArrowStringList is a list of strings with an arrow-separated string representation.
//
// Note: ArrowStringList does not implement TextReferenceResolver. We do not
// want to return references to individual strings in the text log because this
// might lead to confusion; the label would point to the full string list but
// values like offset are relative to the individual string.
type ArrowStringList []string

func (a ArrowStringList) String() string {
	return strings.Join(a, ">")
}

var _ jsonlog.JsonReferenceResolver = ArrowStringList{}

func (a ArrowStringList) RelativeJsonPointer(pointee any) jsonpointer.Pointer {
	stringPointer, isStringPointer := pointee.(*string)
	if !isStringPointer {
		return nil
	}
	for i := range a {
		if &a[i] == stringPointer {
			return jsonpointer.New(strconv.Itoa(i))
		}
	}
	return nil
}
