package thorlog

import (
	"strings"
)

// StringList is a list of strings with a comma-separated string representation.
//
// Note: StringList does not implement TextReferenceResolver. We do not want to return references to individual strings in the text log because this might lead to confusion; the label would point to the full string list but values like offset are relative to the individual string.
type StringList []string

func (s StringList) String() string {
	return strings.Join(s, ", ")
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
