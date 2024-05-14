package jsonlog

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/NextronSystems/jsonlog/jsonpointer"
	"github.com/NextronSystems/jsonlog/thorlog/truncate"
	"golang.org/x/exp/slices"
)

type Event []ReferencedField

type ReferencedField struct {
	Reference Reference
	Value     string
}

type FieldMatch struct {
	FieldPointer any
	truncate.Match
}

type Truncator interface {
	Truncate(matches []FieldMatch, truncateLimit int, stringContext int) Object
}

func CreateEvent(object Object) Event {
	// Walk the object and create a list of all fields
	value := reflect.ValueOf(object)
	if value.Kind() != reflect.Ptr {
		panic("object must be a pointer")
	}
	value = value.Elem()
	return walkObject(object, "", nil, value)
}

type EventValue struct {
	FieldPointer any
	Value        string

	// TextLabel is the label that should be used for this value in the text log. If it's not specified, it is looked up using the text ref resolver.
	TextLabel string
	// JsonPointer is the pointer to the field in the JSON log. If it's not specified, it is looked up using the json ref resolver.
	JsonPointer jsonpointer.Pointer
}

type EventValuer interface {
	Values() []EventValue
}

func walkObject(base any, textlabelPrefix string, pointerPrefix jsonpointer.Pointer, value reflect.Value) []ReferencedField {
	var fields []ReferencedField
	for i := 0; i < value.NumField(); i++ {
		typeField := value.Type().Field(i)
		if !typeField.IsExported() {
			continue
		}
		field := value.Field(i)

		textlogTag := typeField.Tag.Get("textlog")
		tagModifiers := strings.Split(textlogTag, ",")
		logfield := strings.ToUpper(tagModifiers[0])
		tagModifiers = tagModifiers[1:]
		if logfield == "-" {
			continue
		}

		jsonTag := strings.SplitN(typeField.Tag.Get("json"), ",", 2)[0]
		fieldJsonPointer := jsonpointer.New(pointerPrefix...)
		if jsonTag != "" {
			fieldJsonPointer = fieldJsonPointer.Append(jsonTag)
		}
		var fieldLabel = ConcatTextLabels(textlabelPrefix, logfield)
		var subfieldPrefix string
		if slices.Contains(tagModifiers, "expand") {
			subfieldPrefix = fieldLabel
		}

		var fieldPointer = field
		if field.Kind() != reflect.Ptr {
			fieldPointer = field.Addr()
		}
		if valuer, isValuer := fieldPointer.Interface().(EventValuer); isValuer {
			for _, eventValue := range valuer.Values() {
				jsonPointer := jsonpointer.New(fieldJsonPointer...)
				jsonPointer = append(jsonPointer, eventValue.JsonPointer...)
				fields = append(fields, ReferencedField{
					Value: eventValue.Value,
					Reference: Reference{
						Base:         base,
						PointedField: eventValue.FieldPointer,
						textLabel:    ConcatTextLabels(subfieldPrefix, eventValue.TextLabel),
						jsonPointer:  jsonPointer,
					},
				})
			}
			continue
		}
		if field.Kind() == reflect.Interface {
			if field.IsNil() {
				continue
			}
			field = field.Elem()
			if field.Kind() != reflect.Ptr {
				// Interface values that are not pointers can't be addressed, thus we can't get a reference to them
				continue
			}
		}
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				continue
			}
			field = field.Elem()
		}
		shouldExpandStruct := typeField.Anonymous || slices.Contains(tagModifiers, "expand")
		if field.Kind() == reflect.Struct && shouldExpandStruct {
			subfields := walkObject(base, subfieldPrefix, fieldJsonPointer, field)
			fields = append(fields, subfields...)
		} else if logfield != "" {
			if slices.Contains(tagModifiers, "omitempty") && isZero(field) {
				continue
			}
			fields = append(fields, ReferencedField{
				Value: fmt.Sprint(field.Interface()),
				Reference: Reference{
					Base:         base,
					PointedField: field.Addr().Interface(),
					textLabel:    fieldLabel,
					jsonPointer:  fieldJsonPointer,
				},
			})
		}
	}
	return fields
}

type isZeroer interface {
	IsZero() bool
}

func isZero(v reflect.Value) bool {
	if v.IsZero() {
		return true
	}
	if zeroer, ok := v.Interface().(isZeroer); ok {
		return zeroer.IsZero()
	}
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
