package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type PluginStructuredData struct {
	jsonlog.ObjectHeader

	Plugin string `json:"plugin" textlog:"-"`

	Data KeyValueList `json:"data" textlog:",inline"`
}

func (PluginStructuredData) reportable() {}

const typePluginStructuredData = "structured data from plugin"

func init() { AddLogObjectType(typePluginStructuredData, &PluginStructuredData{}) }

func NewPluginStructuredData(plugin string) *PluginStructuredData {
	return &PluginStructuredData{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typePluginStructuredData,
		},
		Plugin: plugin,
	}
}

func (a PluginStructuredData) Truncate(matches []jsonlog.FieldMatch, truncateLimit int, stringContext int) jsonlog.Object {
	a.Data = a.Data.Truncate(matches, truncateLimit, stringContext)
	return &a
}

type PluginString struct {
	jsonlog.ObjectHeader

	Plugin string `json:"plugin" textlog:"-"`

	String string `json:"string" textlog:"string"`
}

func (PluginString) reportable() {}

const typePluginString = "data from plugin"

func init() { AddLogObjectType(typePluginString, &PluginString{}) }

func NewPluginString(plugin string) *PluginString {
	return &PluginString{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typePluginString,
		},
		Plugin: plugin,
	}
}

type PluginFinding struct {
	LogObjectHeader

	Plugin string `json:"plugin" textlog:"-"`

	LogDetails MessageFields `json:"details" textlog:",expand"`
}

func (PluginFinding) reportable() {}

const typePluginFinding = "finding from plugin"

func init() { AddLogObjectType(typePluginFinding, &PluginFinding{}) }

func NewPluginFinding(plugin string) *PluginFinding {
	return &PluginFinding{
		LogObjectHeader: jsonlog.ObjectHeader{
			Type: typePluginFinding,
		},
		Plugin: plugin,
	}
}
