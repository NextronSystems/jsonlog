package thorlog

import (
	"strings"

	"github.com/NextronSystems/jsonlog"
)

// LogObjectTypes is a map of all log object types. Each log object type must be registered using AddLogObjectType.
var LogObjectTypes = map[string]jsonlog.Object{}

// FindLogObjectType looks up a log object type by name, ignoring case. It
// returns the correctly-folded name (with which the type can be found in
// LogObjectTypes) and true if found, or "" and false if not found.
func FindLogObjectType(name string) (string, bool) {
	for registeredName := range LogObjectTypes {
		if strings.EqualFold(registeredName, name) {
			return registeredName, true
		}
	}
	return "", false
}

// AddLogObjectType registers a new log object type. It panics if a log object type with the same name is already registered.
func AddLogObjectType(name string, obj jsonlog.Object) {
	if _, ok := LogObjectTypes[name]; ok {
		panic("duplicate log object type: " + name)
	}
	LogObjectTypes[name] = obj
}
