package jsonlog

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/NextronSystems/jsonlog/jsonpointer"
	"golang.org/x/exp/slices"
)

// NewReference creates a new reference to a field of an Object.
// The base must be a pointer to a struct implementing Object.
// The field must be a pointer to a field within the base.
func NewReference(base Object, field any) *Reference {
	if reflect.ValueOf(base).Kind() != reflect.Pointer {
		panic("Base must be a pointer to a struct implementing Object")
	}
	if reflect.ValueOf(field).Kind() != reflect.Pointer {
		panic("field must be a pointer to a field within base")
	}
	return &Reference{
		Base:         base,
		PointedField: field,
	}
}

// Reference is a reference to a field of an Object
type Reference struct {
	Base         any // Must be a pointer to an Object
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
	if pointedValue.Kind() != reflect.Pointer {
		panic("PointedField must be a pointer")
	}
	r.jsonPointer = findRelativeJsonPointer(baseValue, pointedValue)
	if r.jsonPointer == nil {
		panic("pointed field not found in base")
	}
	return r.jsonPointer
}

var referenceType = reflect.TypeOf((*Reference)(nil)).Elem()

// pointerIfPossible returns a pointer to the specified value if it is addressable.
// Otherwise, it returns the value itself.
func pointerIfPossible(v reflect.Value) reflect.Value {
	if v.CanAddr() {
		return v.Addr()
	}
	return v
}

// findRelativeJsonPointer finds a JSON pointer from base to pointedField. The
// base must be pointer. The pointed field should be a pointer to a field
// within the base; if the pointed field is not found within the base, nil is
// returned.
func findRelativeJsonPointer(base reflect.Value, pointedField reflect.Value) jsonpointer.Pointer {
	// If possible, use a pointer to the base; maybe JsonReferenceResolver is implemented as a pointer method
	base = pointerIfPossible(base)
	for {
		if base.Equal(pointedField) {
			return jsonpointer.Pointer{}
		}
		if resolver, isResolver := base.Interface().(JsonReferenceResolver); isResolver {
			return resolver.RelativeJsonPointer(pointedField.Interface())
		}
		if base.Type() == referenceType { // Don't recurse into other references
			return nil
		}
		if base.Kind() == reflect.Pointer || base.Kind() == reflect.Interface {
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
			pointer := findRelativeJsonPointer(field, pointedField)
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
			elem := base.Index(i)
			pointer := findRelativeJsonPointer(elem, pointedField)
			if pointer == nil {
				continue
			}
			return jsonpointer.New(strconv.Itoa(i)).Append(pointer...)
		}
		return nil
	case reflect.Map:
		if base.Type().Key().Kind() != reflect.String {
			return nil
		}
		for _, key := range base.MapKeys() {
			pointer := findRelativeJsonPointer(base.MapIndex(key), pointedField)
			if pointer == nil {
				continue
			}
			return jsonpointer.New(key.String()).Append(pointer...)
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
	if pointedValue.Kind() != reflect.Pointer {
		panic("PointedField must be a pointer")
	}
	r.textLabel, _ = findTextLabel(baseValue, reflect.ValueOf(r.PointedField))
	return r.textLabel
}

// TextReferenceResolver is an interface that can be implemented by a struct to
// specify custom text labels for its fields that are used in references.
type TextReferenceResolver interface {
	// RelativeTextPointer returns a label for the given field of the object that
	// implements this interface.
	// The given field must be a pointer to a field of the object.
	// Fields are marked virtual if the resolver provides a distinct label for
	// this field but a corresponding TextlogMarshaler method does not list this
	// field/label as a separate key.
	// RelativeTextPointer returns the label as a string, a bool indicating
	// whether the label is virtual, and a bool indicating whether the field was
	// found.
	RelativeTextPointer(pointee any) (string, bool, bool)
}

func findTextLabel(base reflect.Value, pointedField reflect.Value) (string, bool) {
	if base.CanAddr() {
		if base.Addr().Equal(pointedField) {
			return "", true
		}
	}
	if base.Kind() == reflect.Pointer || base.Kind() == reflect.Interface {
		base = base.Elem()
	}
	switch base.Kind() {
	case reflect.Struct:
		for i := 0; i < base.NumField(); i++ {
			field := base.Field(i)
			typefield := base.Type().Field(i)
			if !typefield.IsExported() {
				continue
			}
			var fieldPointer = field
			if field.Kind() != reflect.Pointer && field.CanAddr() {
				fieldPointer = field.Addr()
			}
			var label string
			var labelFound bool
			var labelVirtual bool
			resolver, isResolver := fieldPointer.Interface().(TextReferenceResolver)
			if fieldPointer.Equal(pointedField) {
				label, labelFound = "", true
			} else if isResolver {
				label, labelVirtual, labelFound = resolver.RelativeTextPointer(pointedField.Interface())
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
			if typefield.Anonymous || slices.Contains(tagModifiers, TextlogModifierExpand) || isResolver {
				if labelVirtual {
					// Virtual fields are actually not present in the text log (hence the
					// name), so it does not make sense to combine the existent field label
					// with the non-existent subfield label which would confuse the reader.
					fullLabel = label
				} else {
					fullLabel = ConcatTextLabels(fieldlabel, label)
				}
			} else if label == "" {
				fullLabel = fieldlabel
			} else {
				// Unexpanded fields should not use labels of any subfield because it
				// cannot be resolved unambiguously.
				return "", false
			}
			return fullLabel, true
		}
		return "", false
	case reflect.Slice, reflect.Array:
		for i := 0; i < base.Len(); i++ {
			element := base.Index(i)
			label, labelFound := findTextLabel(element, pointedField)
			if !labelFound {
				continue
			}
			return ConcatTextLabels(label, strconv.Itoa(i+1)), true
		}
		return "", false
	case reflect.Map:
		// Map values aren't addressable and thus the pointedField can't point to them;
		// and they aren't expanded, therefore any references to subfields of map values also have no text label.
		return "", false
	default:
		return "", false
	}
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
