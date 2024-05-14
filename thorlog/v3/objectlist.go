package thorlog

import "github.com/NextronSystems/jsonlog"

// LogObjectTypes is a map of all log object types. Each log object type must be registered using AddLogObjectType.
var LogObjectTypes = map[string]jsonlog.Object{}

// AddLogObjectType registers a new log object type. It panics if a log object type with the same name is already registered.
func AddLogObjectType(name string, obj jsonlog.Object) {
	if _, ok := LogObjectTypes[name]; ok {
		panic("duplicate log object type: " + name)
	}
	LogObjectTypes[name] = obj
}
