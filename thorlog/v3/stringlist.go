package thorlog

import (
	"strconv"
	"strings"

	"github.com/NextronSystems/jsonlog/jsonpointer"
)

type StringList []string

func (s StringList) String() string {
	return strings.Join(s, ", ")
}

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

func (s StringList) RelativeTextPointer(pointee any) (string, bool) {
	stringPointer, isStringPointer := pointee.(*string)
	if !isStringPointer {
		return "", false
	}
	for i := range s {
		if &s[i] == stringPointer {
			return "", true
		}
	}
	return "", false
}

type ArrowStringList []string

func (a ArrowStringList) String() string {
	return strings.Join(a, ">")
}

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

func (a ArrowStringList) RelativeTextPointer(pointee any) (string, bool) {
	stringPointer, isStringPointer := pointee.(*string)
	if !isStringPointer {
		return "", false
	}
	for i := range a {
		if &a[i] == stringPointer {
			return "", true
		}
	}
	return "", false
}
