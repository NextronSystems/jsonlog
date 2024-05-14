package jsonpointer

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Resolve resolves the given JSON pointer against the given base object.
// The base object must be a pointer to a struct or a pointer to a slice/array.
// The returned value is a pointer to the pointed field.
func Resolve(base any, pointer Pointer) (any, error) {
	baseValue := reflect.ValueOf(base)

	var currentValue = baseValue
	for _, token := range pointer {
		var ok bool
		currentValue, ok = findByLabel(currentValue, token)
		if !ok {
			return nil, fmt.Errorf("could not resolve %s: could not find %s", pointer, token)
		}
	}
	return currentValue.Addr().Interface(), nil
}

func findByLabel(base reflect.Value, jsonLabel string) (reflect.Value, bool) {
	for base.Kind() == reflect.Ptr || base.Kind() == reflect.Interface {
		if base.IsNil() {
			return reflect.Value{}, false
		}
		base = base.Elem()
	}
	switch base.Kind() {
	case reflect.Struct:
		return findStructFieldByLabel(base, jsonLabel)
	case reflect.Array, reflect.Slice:
		return findArrayElementByLabel(base, jsonLabel)
	default:
		return reflect.Value{}, false
	}
}

func findArrayElementByLabel(base reflect.Value, label string) (reflect.Value, bool) {
	index, err := strconv.Atoi(label)
	if err != nil {
		return reflect.Value{}, false
	}
	if index < 0 || index >= base.Len() {
		return reflect.Value{}, false
	}
	return base.Index(index), true
}

func findStructFieldByLabel(base reflect.Value, label string) (reflect.Value, bool) {
	structType := base.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if !field.IsExported() {
			continue
		}
		if strings.SplitN(field.Tag.Get("json"), ",", 2)[0] == label {
			return base.Field(i), true
		}
	}
	return reflect.Value{}, false
}
