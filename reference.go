package jsonlog

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/NextronSystems/jsonlog/jsonpointer"
	"golang.org/x/exp/slices"
)

// NewReference creates a new reference to a field of a Object.
// The base must be a pointer to a struct implementing Object.
// The field must be a pointer to a field within the base.
func NewReference(base Object, field any) *Reference {
	if reflect.ValueOf(base).Kind() != reflect.Ptr {
		panic("Base must be a pointer to a struct implementing Object")
	}
	if reflect.ValueOf(field).Kind() != reflect.Ptr {
		panic("field must be a pointer to a field within base")
	}
	return &Reference{
		Base:         base,
		PointedField: field,
	}
}

// Reference is a reference to a field of a Object
type Reference struct {
	Base         any // Must be a pointer to a Object
	PointedField any // Must be a pointer to a (possibly nested) field of Base

	textLabel   string
	jsonPointer jsonpointer.Pointer
}

// ToJsonPointer returns a JSON pointer to the pointed field.
func (r *Reference) ToJsonPointer() jsonpointer.Pointer {
	if r.jsonPointer != nil {
		return r.jsonPointer
	}
	baseValue := reflect.ValueOf(r.Base)
	pointedValue := reflect.ValueOf(r.PointedField)
	if pointedValue.Kind() != reflect.Ptr {
		panic("PointedField must be a pointer")
	}
	r.jsonPointer = findRelativeJsonPointer(baseValue, pointedValue)
	if r.jsonPointer == nil {
		panic("pointed field not found in base")
	}
	return r.jsonPointer
}

func findRelativeJsonPointer(base reflect.Value, pointedField reflect.Value) jsonpointer.Pointer {
	for {
		if base.Equal(pointedField) {
			return jsonpointer.Pointer{}
		}
		if resolver, isResolver := base.Interface().(JsonReferenceResolver); isResolver {
			return resolver.RelativeJsonPointer(pointedField.Interface())
		}
		if base.Kind() == reflect.Ptr || base.Kind() == reflect.Interface {
			if base.IsNil() {
				return nil
			}
			base = base.Elem()
		} else {
			break
		}
	}
	switch base.Kind() {
	case reflect.Struct:
		for i := 0; i < base.NumField(); i++ {
			field := base.Field(i)
			typefield := base.Type().Field(i)
			if !typefield.IsExported() {
				continue
			}
			pointer := findRelativeJsonPointer(field.Addr(), pointedField)
			if pointer == nil {
				continue
			}
			if typefield.Anonymous {
				return pointer
			}
			jsonTag := strings.SplitN(typefield.Tag.Get("json"), ",", 2)[0]
			return jsonpointer.New(jsonTag).Append(pointer...)
		}
		return nil
	case reflect.Slice, reflect.Array:
		for i := 0; i < base.Len(); i++ {
			pointer := findRelativeJsonPointer(base.Index(i).Addr(), pointedField)
			if pointer == nil {
				continue
			}
			return jsonpointer.New(strconv.Itoa(i)).Append(pointer...)
		}
		return nil
	default:
		return nil
	}
}

// JsonReferenceResolver is an interface that can be implemented by a struct to create custom JSON pointers to its fields.
// This interface should be implemented by log objects when the object has a custom MarshalJSON() and thus the default JSON pointer would not work.
type JsonReferenceResolver interface {
	// RelativeJsonPointer returns a pointer to the given field of the object that implements this interface.
	// If the field is not found, nil is returned.
	// The given field must be a pointer to a field of the object.
	RelativeJsonPointer(pointee any) jsonpointer.Pointer
}

// ToTextLabel returns a text label for the pointed field.
func (r *Reference) ToTextLabel() string {
	if r.textLabel != "" {
		return r.textLabel
	}
	baseValue := reflect.ValueOf(r.Base).Elem()
	pointedValue := reflect.ValueOf(r.PointedField)
	if pointedValue.Kind() != reflect.Ptr {
		panic("PointedField must be a pointer")
	}
	r.textLabel, _ = findTextLabel(baseValue, reflect.ValueOf(r.PointedField))
	return r.textLabel
}

// TextReferenceResolver is an interface that can be implemented by a struct to specify custom text labels for its fields
// that are used in references.
type TextReferenceResolver interface {
	// RelativeTextPointer returns a label for the given field of the object that implements this interface.
	// If the field is not found, nil is returned.
	// The given field must be a pointer to a field of the object.
	RelativeTextPointer(pointee any) (string, bool)
}

func findTextLabel(base reflect.Value, pointedField reflect.Value) (string, bool) {
	if base.Addr().Equal(pointedField) {
		return "", true
	}
	if base.Kind() == reflect.Ptr || base.Kind() == reflect.Interface {
		base = base.Elem()
	}
	if base.Kind() != reflect.Struct {
		return "", false
	}
	for i := 0; i < base.NumField(); i++ {
		field := base.Field(i)
		typefield := base.Type().Field(i)
		if !typefield.IsExported() {
			continue
		}
		var fieldPointer = field
		if field.Kind() != reflect.Ptr {
			fieldPointer = field.Addr()
		}
		var label string
		var labelFound bool
		if fieldPointer.Equal(pointedField) {
			label, labelFound = "", true
		} else if resolver, isResolver := fieldPointer.Interface().(TextReferenceResolver); isResolver {
			label, labelFound = resolver.RelativeTextPointer(pointedField.Interface())
		} else {
			label, labelFound = findTextLabel(field, pointedField)
		}
		if !labelFound {
			continue
		}
		textlogTag := typefield.Tag.Get("textlog")
		tagModifiers := strings.Split(textlogTag, ",")
		fieldlabel := strings.ToUpper(tagModifiers[0])
		tagModifiers = tagModifiers[1:]
		var fullLabel string
		if slices.Contains(tagModifiers, "expand") {
			fullLabel = ConcatTextLabels(fieldlabel, label)
		} else if label == "" {
			fullLabel = fieldlabel
		} else {
			fullLabel = label
		}
		return fullLabel, true
	}
	return "", false
}

func ConcatTextLabels(prefix string, label string) string {
	if prefix == "" {
		return label
	}
	if label == "" {
		return prefix
	}
	if prefix == label {
		return label
	}
	return prefix + "_" + label
}

// String returns the text label of the pointed field.
func (r Reference) String() string {
	return r.ToTextLabel()
}

// MarshalJSON returns the JSON pointer to the pointed field as a JSON string.
func (r Reference) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.ToJsonPointer().String())
}

func (r *Reference) UnmarshalJSON(data []byte) error {
	var pointerString string
	err := json.Unmarshal(data, &pointerString)
	if err != nil {
		return err
	}
	pointer, err := jsonpointer.Parse(pointerString)
	if err != nil {
		return err
	}
	r.jsonPointer = pointer
	return nil
}

// Value returns the pointed field.
func (r Reference) Value() any {
	return reflect.ValueOf(r.PointedField).Elem().Interface()
}

func (r Reference) JSONSchemaAlias() any {
	return ""
}

// SetLabels sets the JSON pointer and text label for the reference.
// This is an unsafe operation since it does not update the base and pointed field,
// nor does it check if the passed values match the pointed field.
// However, it can be used on a fresh reference where the labels are already known
// to avoid the overhead of looking them up.
func (r *Reference) SetLabels(jsonPointer jsonpointer.Pointer, textLabel string) {
	r.jsonPointer = jsonPointer
	r.textLabel = textLabel
}
