package jsonlog

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

// TextlogEntry represents a single entry in a text log. It is a list of key-value pairs.
type TextlogEntry []TextlogValuePair

// TextlogValuePair is a key-value pair in a text log entry.
type TextlogValuePair struct {
	Key   string
	Value string
}

// TextlogFormatter is a formatter for text logs. It provides a way to format log objects into text log entries.
// Text log keys are derived from the textlog tag in the struct field tags. If the tag is empty, the field is not marshalled.
// This tag can contain modifiers that control how the field is marshalled. Supported modifiers are:
//
//   - expand: causes the field to be expanded into its subfields. Each subfield is prefixed with this field's tag name.
//   - omitempty: causes the field to be omitted if it is the zero value or implements IsZero() and it returns true.
//   - explicit: causes the field to be marshalled even if the tag name is empty (requires that the containing struct has a non-empty tag name to avoid an empty text log key).
//
// Slice and map values are marshalled as a list of key-value pairs. The key is the index (one-based) for slices and the map key for maps.
type TextlogFormatter struct {
	// FormatValue is a function that formats a single value into a string. If it is nil, fmt.Sprint is used.
	FormatValue func(data any, modifiers []string) string
	// Omit is a function that determines whether a field should be omitted from the log entry.
	// If it is nil, no fields are omitted.
	Omit func(modifiers []string, value any) bool
}

func (t TextlogFormatter) format(data any, modifiers []string) string {
	if t.FormatValue != nil {
		return t.FormatValue(data, modifiers)
	}
	return fmt.Sprint(data)
}

// Format formats an object into a text log entry.
// The object must be a struct, pointer to a struct, slice, or map.
func (t TextlogFormatter) Format(object any) TextlogEntry {
	entry := t.toEntry(reflect.ValueOf(object))
	// Deduplicate keys
	// Keys should already be unique, but this is not guaranteed and we need to guarantee this property for downstream consumers
	keys := make(map[string]struct{})
	for j := range entry {
		pair := &entry[j]
		if _, exists := keys[pair.Key]; exists {
			for i := 2; ; i++ {
				newKey := pair.Key + "_" + strconv.Itoa(i)
				if _, exists := keys[newKey]; !exists {
					pair.Key = newKey
					break
				}
			}
		}
		keys[pair.Key] = struct{}{}
	}
	return entry
}

const (
	// modifierExpand, when applied to a struct field, causes the field to be expanded into its subfields
	modifierExpand = "expand"
	// modifierOmitempty causes a field to be omitted if it is the zero value or implements IsZero() and it returns true
	modifierOmitempty = "omitempty"
	// modifierExplicit causes a field to be marshalled even if the tag name is empty
	modifierExplicit = "explicit"
)

type TextlogMarshaler interface {
	MarshalTextLog(formatter TextlogFormatter) TextlogEntry
}

func withUppercaseKeys(entry TextlogEntry) TextlogEntry {
	for i := range entry {
		entry[i].Key = strings.ToUpper(entry[i].Key)
	}
	return entry
}

func (t TextlogFormatter) toEntry(object reflect.Value) TextlogEntry {
	for object.Kind() == reflect.Ptr || object.Kind() == reflect.Interface {
		if object.IsNil() {
			return nil
		}
		if marshaler, ok := object.Interface().(TextlogMarshaler); ok {
			return withUppercaseKeys(marshaler.MarshalTextLog(t))
		}
		object = object.Elem()
	}
	if object.Kind() == reflect.Invalid {
		return nil
	}
	if marshaler, ok := object.Interface().(TextlogMarshaler); ok {
		return withUppercaseKeys(marshaler.MarshalTextLog(t))
	}
	switch object.Kind() {
	case reflect.Struct:
		var details TextlogEntry
		for i := 0; i < object.NumField(); i++ {
			typeField := object.Type().Field(i)
			if !typeField.IsExported() {
				continue
			}
			field := object.Field(i)
			textlogTag := typeField.Tag.Get("textlog")
			tagModifiers := strings.Split(textlogTag, ",")
			logfield := strings.ToUpper(tagModifiers[0])
			tagModifiers = tagModifiers[1:]
			if logfield == "-" {
				continue
			}
			if !typeField.Anonymous && textlogTag == "" && !slices.Contains(tagModifiers, modifierExplicit) {
				continue
			}
			if slices.Contains(tagModifiers, modifierOmitempty) && isZero(field) {
				continue
			}
			if t.Omit != nil && t.Omit(tagModifiers, field.Interface()) {
				continue
			}
			if typeField.Anonymous || slices.Contains(tagModifiers, modifierExpand) {
				// Use the tag as a prefix for the subfields
				subentry := t.toEntry(field)
				for _, subentryValue := range subentry {
					details = append(details, TextlogValuePair{
						Key:   ConcatTextLabels(logfield, subentryValue.Key),
						Value: subentryValue.Value,
					})
				}
			} else {
				// Add the field as a single value
				key := logfield
				if field.Kind() == reflect.Ptr {
					field = field.Elem()
				}
				details = append(details, TextlogValuePair{
					Key:   key,
					Value: t.format(field.Interface(), tagModifiers),
				})
			}
		}
		return details
	case reflect.Slice:
		var details TextlogEntry
		for i := 0; i < object.Len(); i++ {
			subentry := t.toEntry(object.Index(i))
			for _, subentryValue := range subentry {
				details = append(details, TextlogValuePair{
					Key:   ConcatTextLabels(subentryValue.Key, strconv.Itoa(i+1)),
					Value: subentryValue.Value,
				})
			}
		}
		return details
	case reflect.Map:
		if object.Type().Key().Kind() != reflect.String {
			return nil
		}
		var details TextlogEntry
		for _, key := range object.MapKeys() {
			details = append(details, TextlogValuePair{
				Key:   key.String(),
				Value: t.format(object.MapIndex(key).Interface(), nil),
			})
		}
		return details
	default:
		return nil
	}
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
	if v.Comparable() {
		return v.Equal(reflect.Zero(v.Type()))
	}
	return false
}
