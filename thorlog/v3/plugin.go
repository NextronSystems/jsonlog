package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

// PluginStructuredData contains data that was passed to THOR by a plugin in order to be scanned.
//
// The keys and values in the given Data are defined by the plugin and thus cannot be known beforehand.
type PluginStructuredData struct {
	jsonlog.ObjectHeader

	// The plugin that passed the data to THOR
	Plugin string `json:"plugin" textlog:"-"`

	// The data that was passed to THOR by the plugin
	Data KeyValueList `json:"data" textlog:"data"`
}

func (PluginStructuredData) observed() {}

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
